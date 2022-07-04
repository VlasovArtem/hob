package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/VlasovArtem/hob/src/pivotal/cache"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/VlasovArtem/hob/src/pivotal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PivotalServiceObject struct {
	housePivotalRepository repository.PivotalRepository[model.HousePivotal]
	groupPivotalRepository repository.PivotalRepository[model.GroupPivotal]
	houseService           houseService.HouseService
	groupService           groupService.GroupService
	pivotalCache           cache.PivotalCache
	incomeService          incomeService.IncomeService
	paymentService         paymentService.PaymentService
}

func NewPivotalService(
	housePivotalRepository repository.PivotalRepository[model.HousePivotal],
	groupPivotalRepository repository.PivotalRepository[model.GroupPivotal],
	houseService houseService.HouseService,
	groupService groupService.GroupService,
	pivotalCache cache.PivotalCache,
	incomeService incomeService.IncomeService,
	paymentService paymentService.PaymentService,
) PivotalService {
	return &PivotalServiceObject{
		housePivotalRepository: housePivotalRepository,
		groupPivotalRepository: groupPivotalRepository,
		houseService:           houseService,
		groupService:           groupService,
		pivotalCache:           pivotalCache,
		incomeService:          incomeService,
		paymentService:         paymentService,
	}
}

func (p *PivotalServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalService(
		dependency.FindRequiredDependency[repository.HousePivotalRepository[model.HousePivotal], repository.PivotalRepository[model.HousePivotal]](factory),
		dependency.FindRequiredDependency[repository.GroupPivotalRepository[model.GroupPivotal], repository.PivotalRepository[model.GroupPivotal]](factory),
		dependency.FindRequiredDependency[houseService.HouseServiceObject, houseService.HouseService](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceObject, groupService.GroupService](factory),
		dependency.FindRequiredDependency[cache.PivotalCacheObject, cache.PivotalCache](factory),
		dependency.FindRequiredDependency[incomeService.IncomeServiceObject, incomeService.IncomeService](factory),
		dependency.FindRequiredDependency[paymentService.PaymentServiceObject, paymentService.PaymentService](factory),
	)
}

type PivotalService interface {
	FindByHouseId(houseId uuid.UUID) (model.HousePivotalDto, error)
	FindByGroupId(groupId uuid.UUID) (model.GroupPivotalDto, error)
	AddPayment(db *gorm.DB, paymentSum float64, houseId uuid.UUID) error
	AddIncome(db *gorm.DB, incomeSum float64, houseId uuid.UUID) error
}

func (p *PivotalServiceObject) AddPayment(db *gorm.DB, paymentSum float64, houseId uuid.UUID) error {
	if houseDto, err := p.houseService.FindById(houseId); err != nil {
		return err
	} else {
		updateTime := time.Now()

		for _, group := range houseDto.Groups {
			var groupPivotal model.GroupPivotal
			if err = p.groupPivotalRepository.FindBySourceId(group.Id, &groupPivotal); err != nil && !errors.Is(err, interrors.ErrNotFound{}) {
				return err
			} else if groupPivotal.Pivotal.Id != uuid.Nil {
				groupPivotal.Pivotal.Payments += paymentSum
				groupPivotal.Pivotal.Total -= paymentSum
				groupPivotal.Pivotal.LatestPaymentUpdateDate = updateTime
				p.groupPivotalRepository.UpdateTransactional(db)
			}
		}
	}
}

func (p *PivotalServiceObject) AddIncome(db *gorm.DB, incomeSum float64, houseId uuid.UUID) error {

}

func (p *PivotalServiceObject) FindByHouseId(houseId uuid.UUID) (response model.HousePivotalDto, err error) {
	if !p.houseService.ExistsById(houseId) {
		return model.HousePivotalDto{}, interrors.NewErrNotFound("House with id %s not found", houseId)
	}

	var pivotal model.HousePivotal

	if err = p.housePivotalRepository.FindBySourceId(houseId, &pivotal); err != nil {
		return
	}

	return

	//pivotalExists := pivotal.HouseId != uuid.Nil
	//
	//if err = p.findById(houseId, &pivotal.Pivotal,
	//	func(from *time.Time) (float64, error) {
	//		return p.incomeService.CalculateSumByHouseId(houseId, from)
	//	},
	//	func(from *time.Time) (float64, error) {
	//		return p.paymentService.CalculateSum([]uuid.UUID{houseId}, from)
	//	},
	//); err != nil {
	//	return
	//}
	//
	//pivotal.HouseId = houseId
	//
	//if !pivotalExists {
	//	pivotal.Pivotal.Id = uuid.New()
	//	createdPivotal, err := p.housePivotalRepository.Create(pivotal)
	//	if err != nil {
	//		return model.HousePivotalDto{}, err
	//	} else {
	//		return createdPivotal.ToDto(), nil
	//	}
	//}
	//
	//if err := p.housePivotalRepository.Update(pivotal.Pivotal.Id, pivotal.Pivotal.Total, pivotal.Pivotal.LatestIncomeUpdateDate, pivotal.Pivotal.LatestPaymentUpdateDate); err != nil {
	//	return model.HousePivotalDto{}, err
	//}
	//
	//return pivotal.ToDto(), nil
}

func (p *PivotalServiceObject) FindByGroupId(groupId uuid.UUID) (response model.GroupPivotalDto, err error) {
	if !p.groupService.ExistsById(groupId) {
		return model.GroupPivotalDto{}, interrors.NewErrNotFound("Group with id %s not found", groupId)
	}

	var pivotal model.GroupPivotal

	if err = p.groupPivotalRepository.FindBySourceId(groupId, &pivotal); err != nil {
		return
	}

	pivotalExists := pivotal.GroupId != uuid.Nil

	if err = p.findById(groupId, &pivotal.Pivotal,
		func(from *time.Time) (float64, error) {
			return p.incomeService.CalculateSumByGroupId(groupId, from)
		},
		func(from *time.Time) (float64, error) {
			houseIds := common.MapSlice(p.houseService.FindHousesByGroupId(groupId), func(house houseModel.HouseDto) uuid.UUID { return house.Id })

			return p.paymentService.CalculateSum(houseIds, from)
		}); err != nil {
		return model.GroupPivotalDto{}, err
	}

	pivotal.GroupId = groupId

	if !pivotalExists {
		pivotal.Pivotal.Id = uuid.New()
		createdPivotal, err := p.groupPivotalRepository.Create(pivotal)
		if err != nil {
			return model.GroupPivotalDto{}, err
		} else {
			return createdPivotal.ToDto(), nil
		}
	}

	if err := p.groupPivotalRepository.Update(pivotal.Pivotal.Id, pivotal.Pivotal.Total, pivotal.Pivotal.LatestIncomeUpdateDate, pivotal.Pivotal.LatestPaymentUpdateDate); err != nil {
		return model.GroupPivotalDto{}, err
	}

	return pivotal.ToDto(), nil
}

func (p *PivotalServiceObject) findById(
	id uuid.UUID,
	pivotal *model.Pivotal,
	incomeSumProvider func(from *time.Time) (float64, error),
	paymentSumProvider func(from *time.Time) (float64, error),
) (err error) {
	if cachedPivotal := p.pivotalCache.Find(id); cachedPivotal != nil {
		return
	}

	var latestIncomesDate, latestPaymentDate *time.Time

	if pivotal.Id != uuid.Nil {
		latestIncomesDate = &pivotal.LatestIncomeUpdateDate
		latestPaymentDate = &pivotal.LatestPaymentUpdateDate

		pivotal.LatestIncomeUpdateDate = pivotal.LatestIncomeUpdateDate.Add(time.Millisecond * 1)
		pivotal.LatestPaymentUpdateDate = pivotal.LatestPaymentUpdateDate.Add(time.Millisecond * 1)
	}

	incomeSum, err := incomeSumProvider(latestIncomesDate)
	if err != nil {
		return
	}
	paymentSum, err := paymentSumProvider(latestPaymentDate)
	if err != nil {
		return
	}

	pivotal.Total = incomeSum - paymentSum
	pivotal.Income = incomeSum
	pivotal.Payments = paymentSum

	p.pivotalCache.Add(id, pivotal)

	return
}
