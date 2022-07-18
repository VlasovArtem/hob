package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	houses "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/repository"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	providers "github.com/VlasovArtem/hob/src/provider/service"
	users "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PaymentServiceStr struct {
	userService       users.UserService
	houseService      houses.HouseService
	providerService   providers.ProviderService
	paymentRepository repository.PaymentRepository
	pivotalService    pivotalService.PivotalService
}

func NewPaymentService(
	userService users.UserService,
	houseService houses.HouseService,
	providerService providers.ProviderService,
	paymentRepository repository.PaymentRepository,
	pivotalService pivotalService.PivotalService,
) PaymentService {
	return &PaymentServiceStr{
		userService:       userService,
		houseService:      houseService,
		providerService:   providerService,
		paymentRepository: paymentRepository,
		pivotalService:    pivotalService,
	}
}

func (p *PaymentServiceStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(users.UserServiceStr{}),
		dependency.FindNameAndType(houses.HouseServiceStr{}),
		dependency.FindNameAndType(providers.ProviderServiceStr{}),
		dependency.FindNameAndType(repository.PaymentRepositoryStr{}),
		dependency.FindNameAndType(pivotalService.PivotalServiceStr{}),
	}
}

func (p *PaymentServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentService(
		dependency.FindRequiredDependency[users.UserServiceStr, users.UserService](factory),
		dependency.FindRequiredDependency[houses.HouseServiceStr, houses.HouseService](factory),
		dependency.FindRequiredDependency[providers.ProviderServiceStr, providers.ProviderService](factory),
		dependency.FindRequiredDependency[repository.PaymentRepositoryStr, repository.PaymentRepository](factory),
		dependency.FindRequiredDependency[pivotalService.PivotalServiceStr, pivotalService.PivotalService](factory),
	)
}

type PaymentService interface {
	Add(request model.CreatePaymentRequest) (model.PaymentDto, error)
	AddBatch(request model.CreatePaymentBatchRequest) ([]model.PaymentDto, error)
	FindById(id uuid.UUID) (model.PaymentDto, error)
	FindByHouseId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByGroupId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByUserId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	FindByProviderId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdatePaymentRequest) error
	Transactional(db *gorm.DB) PaymentService
}

func (p *PaymentServiceStr) Add(request model.CreatePaymentRequest) (response model.PaymentDto, err error) {
	if !p.userService.ExistsById(request.UserId) {
		return response, fmt.Errorf("user with id %s not found", request.UserId)
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, fmt.Errorf("house with id %s not found", request.HouseId)
	}

	if request.ProviderId != nil {
		if !p.providerService.ExistsById(*request.ProviderId) {
			return response, fmt.Errorf("provider with id %s not found", request.ProviderId)
		}
	}

	var payment model.Payment

	err = p.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
		trx := p.Transactional(tx).(*PaymentServiceStr)

		payment = request.ToEntity()

		if err = trx.paymentRepository.Create(&payment); err != nil {
			return err
		}

		if trx.pivotalService.ExistsByHouseId(request.HouseId) {
			return trx.pivotalService.AddPayment(float64(request.Sum), request.Date.Add(1*time.Microsecond), request.HouseId)
		}
		return nil
	})

	return payment.ToDto(), err
}

func (p *PaymentServiceStr) AddBatch(request model.CreatePaymentBatchRequest) (response []model.PaymentDto, err error) {
	if len(request.Payments) == 0 {
		return make([]model.PaymentDto, 0), nil
	}

	userIds := make(map[uuid.UUID]bool)
	houseIds := make(map[uuid.UUID]bool)
	providerIds := make(map[uuid.UUID]bool)

	entities := common.MapSlice(request.Payments, func(paymentRequest model.CreatePaymentRequest) model.Payment {
		userIds[paymentRequest.UserId] = true
		houseIds[paymentRequest.HouseId] = true
		if paymentRequest.ProviderId != nil {
			providerIds[*paymentRequest.ProviderId] = true
		}

		return paymentRequest.ToEntity()
	})

	builder := interrors.NewBuilder()

	for userId := range userIds {
		if !p.userService.ExistsById(userId) {
			builder.WithDetail(fmt.Sprintf("user with id %s not found", userId))
		}
	}

	for houseId := range houseIds {
		if !p.houseService.ExistsById(houseId) {
			builder.WithDetail(fmt.Sprintf("house with id %s not found", houseId))
		}
	}

	for providerId := range providerIds {
		if !p.providerService.ExistsById(providerId) {
			builder.WithDetail(fmt.Sprintf("provider with id %s not found", providerId))
		}
	}

	if builder.HasErrors() {
		return nil, interrors.NewErrResponse(builder.WithMessage("Create payment batch failed"))
	}

	err = p.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
		trx := p.Transactional(tx).(*PaymentServiceStr)

		if err = p.paymentRepository.Create(&entities); err != nil {
			return err
		}

		for _, payment := range entities {
			if trx.pivotalService.ExistsByHouseId(payment.HouseId) {
				return trx.pivotalService.AddPayment(float64(payment.Sum), payment.Date.Add(1*time.Microsecond), payment.HouseId)
			}
		}

		return nil
	})

	return common.MapSlice(entities, model.EntityToDto), nil
}

func (p *PaymentServiceStr) FindById(id uuid.UUID) (model.PaymentDto, error) {
	if payment, err := p.paymentRepository.Find(id); err != nil {
		return model.PaymentDto{}, database.HandlerFindError(err, fmt.Sprintf("payment with id %s not found", id))
	} else {
		return payment.ToDto(), nil
	}
}

func (p *PaymentServiceStr) FindByHouseId(houseId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto {
	return p.paymentRepository.FindByHouseId(houseId, limit, offset, from, to)
}

func (p *PaymentServiceStr) FindByUserId(userId uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto {
	return p.paymentRepository.FindByUserId(userId, limit, offset, from, to)
}

func (p *PaymentServiceStr) FindByGroupId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto {
	return p.paymentRepository.FindByGroupId(id, limit, offset, from, to)
}

func (p *PaymentServiceStr) FindByProviderId(id uuid.UUID, limit int, offset int, from, to *time.Time) []model.PaymentDto {
	return p.paymentRepository.FindByProviderId(id, limit, offset, from, to)
}

func (p *PaymentServiceStr) ExistsById(id uuid.UUID) bool {
	return p.paymentRepository.Exists(id)
}

func (p *PaymentServiceStr) DeleteById(id uuid.UUID) error {
	if !p.ExistsById(id) {
		return fmt.Errorf("payment with id %s not found", id)
	}
	return p.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
		trx := p.Transactional(tx).(*PaymentServiceStr)

		if payment, err := trx.paymentRepository.Find(id); err != nil {
			return err
		} else {
			if trx.pivotalService.ExistsByHouseId(payment.HouseId) {
				if err = trx.pivotalService.DeletePayment(float64(payment.Sum), payment.HouseId); err != nil {
					return err
				}
			}
			return trx.DeleteById(id)
		}
	})
}

func (p *PaymentServiceStr) Update(id uuid.UUID, request model.UpdatePaymentRequest) error {
	if !p.ExistsById(id) {
		return fmt.Errorf("payment with id %s not found", id)
	}
	if request.ProviderId != nil && !p.providerService.ExistsById(*request.ProviderId) {
		return fmt.Errorf("provider with id %s not found", request.ProviderId)
	}
	if request.Date.After(time.Now()) {
		return errors.New("date should not be after current date")
	}
	return p.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
		trx := p.Transactional(tx).(*PaymentServiceStr)

		if payment, err := trx.paymentRepository.Find(id); err != nil {
			return err
		} else {
			if trx.pivotalService.ExistsByHouseId(payment.HouseId) {
				if err = trx.pivotalService.UpdatePayment(float64(payment.Sum), float64(request.Sum), request.Date, payment.HouseId); err != nil {
					return err
				}
			}
			return trx.paymentRepository.Update(id, request.UpdateToEntity(id), "HouseId", "House", "UserId", "User")
		}
	})
}

func (p *PaymentServiceStr) Transactional(db *gorm.DB) PaymentService {
	return &PaymentServiceStr{
		userService:       p.userService.Transactional(db),
		houseService:      p.houseService.Transactional(db),
		providerService:   p.providerService.Transactional(db),
		paymentRepository: p.paymentRepository.Transactional(db),
		pivotalService:    p.pivotalService.Transactional(db),
	}
}
