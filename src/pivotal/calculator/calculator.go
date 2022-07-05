package calculator

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/transactional"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/VlasovArtem/hob/src/pivotal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PivotalCalculatorServiceStr struct {
	housePivotalRepository repository.PivotalRepository[model.HousePivotal]
	groupPivotalRepository repository.PivotalRepository[model.GroupPivotal]
	houseService           houseService.HouseService
	incomeService          incomeService.IncomeService
	paymentService         paymentService.PaymentService
}

func NewPivotalCalculatorService(
	housePivotalRepository repository.PivotalRepository[model.HousePivotal],
	groupPivotalRepository repository.PivotalRepository[model.GroupPivotal],
	houseService houseService.HouseService,
	incomeService incomeService.IncomeService,
	paymentService paymentService.PaymentService,
) PivotalCalculatorService {
	return &PivotalCalculatorServiceStr{
		housePivotalRepository: housePivotalRepository,
		groupPivotalRepository: groupPivotalRepository,
		houseService:           houseService,
		incomeService:          incomeService,
		paymentService:         paymentService,
	}
}

func (p *PivotalCalculatorServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalCalculatorService(
		dependency.FindRequiredDependency[repository.HousePivotalRepository, repository.PivotalRepository[model.HousePivotal]](factory),
		dependency.FindRequiredDependency[repository.GroupPivotalRepository, repository.PivotalRepository[model.GroupPivotal]](factory),
		dependency.FindRequiredDependency[houseService.HouseServiceStr, houseService.HouseService](factory),
		dependency.FindRequiredDependency[incomeService.IncomeServiceStr, incomeService.IncomeService](factory),
		dependency.FindRequiredDependency[paymentService.PaymentServiceStr, paymentService.PaymentService](factory),
	)
}

type PivotalCalculatorService interface {
	transactional.Transactional[PivotalCalculatorService]
	Calculate(houseId uuid.UUID) (model.PivotalResponseDto, error)
}

func (p *PivotalCalculatorServiceStr) Transactional(tx *gorm.DB) PivotalCalculatorService {
	return &PivotalCalculatorServiceStr{
		housePivotalRepository: p.housePivotalRepository.Transactional(tx),
		groupPivotalRepository: p.groupPivotalRepository.Transactional(tx),
		houseService:           p.houseService.Transactional(tx),
		incomeService:          p.incomeService.Transactional(tx),
		paymentService:         p.paymentService.Transactional(tx),
	}
}

func (p *PivotalCalculatorServiceStr) Calculate(houseId uuid.UUID) (response model.PivotalResponseDto, err error) {
	housePivotal := model.HousePivotal{
		HouseId: houseId,
		Pivotal: model.Pivotal{
			Id: uuid.New(),
		},
	}
	var groupPivotals []model.GroupPivotal

	err = p.housePivotalRepository.DB().Transaction(func(tx *gorm.DB) error {
		trx := p.Transactional(tx).(*PivotalCalculatorServiceStr)

		houseDto, err := trx.houseService.FindById(houseId)
		if err != nil {
			return err
		}

		if err = trx.housePivotalRepository.DeleteBy("house_id", houseId); err != nil {
			return err
		}

		if len(houseDto.Groups) != 0 {
			for _, groupId := range common.MapSlice[groupModel.GroupDto, uuid.UUID](houseDto.Groups, func(group groupModel.GroupDto) uuid.UUID { return group.Id }) {
				pivotal := model.GroupPivotal{
					GroupId: groupId,
					Pivotal: model.Pivotal{
						Id: uuid.New(),
					},
				}
				if err := trx.groupPivotalRepository.DeleteBy("group_id", groupId); err != nil {
					return err
				}
				for i, paymentDto := range trx.paymentService.FindByGroupId(groupId, -1, -1, nil, nil) {
					if i == 0 {
						pivotal.Pivotal.LatestPaymentUpdateDate = paymentDto.Date.Add(1 * time.Microsecond)
					}
					pivotal.Pivotal.Payments += float64(paymentDto.Sum)
					pivotal.Pivotal.Total -= float64(paymentDto.Sum)
				}
				for i, incomeDto := range trx.incomeService.FindByGroupIds([]uuid.UUID{groupId}, -1, -1, nil, nil) {
					if i == 0 {
						pivotal.Pivotal.LatestIncomeUpdateDate = incomeDto.Date.Add(1 * time.Microsecond)
					}
					pivotal.Pivotal.Income += float64(incomeDto.Sum)
					pivotal.Pivotal.Total += float64(incomeDto.Sum)
				}
				groupPivotals = append(groupPivotals, pivotal)
			}
			if err = trx.groupPivotalRepository.Create(&groupPivotals, "Groups.*"); err != nil {
				return err
			}
		}

		for i, paymentDto := range trx.paymentService.FindByHouseId(houseId, -1, -1, nil, nil) {
			if i == 0 {
				housePivotal.Pivotal.LatestPaymentUpdateDate = paymentDto.Date.Add(1 * time.Microsecond)
			}
			housePivotal.Pivotal.Payments += float64(paymentDto.Sum)
			housePivotal.Pivotal.Total -= float64(paymentDto.Sum)
		}

		for i, incomeDto := range trx.incomeService.FindByHouseId(houseId, -1, -1, nil, nil) {
			if i == 0 {
				housePivotal.Pivotal.LatestIncomeUpdateDate = incomeDto.Date.Add(1 * time.Microsecond)
			}
			housePivotal.Pivotal.Income += float64(incomeDto.Sum)
			housePivotal.Pivotal.Total += float64(incomeDto.Sum)
		}

		if err = trx.housePivotalRepository.Create(&housePivotal); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	response.Groups = common.MapSlice(groupPivotals, func(groupPivotal model.GroupPivotal) model.GroupPivotalDto { return groupPivotal.ToDto() })
	response.House = housePivotal.ToDto()
	response.Total = calculateTotal(housePivotal, groupPivotals)

	return
}

func calculateTotal(house model.HousePivotal, groups []model.GroupPivotal) (pivotal model.TotalPivotalDto) {
	if house.Pivotal.Id != uuid.Nil && len(groups) == 0 {
		pivotal.Income = house.Pivotal.Income
		pivotal.Payments = house.Pivotal.Payments
		pivotal.Total = house.Pivotal.Income - house.Pivotal.Payments
	}
	for _, group := range groups {
		if group.Pivotal.Id != uuid.Nil {
			pivotal.Income += group.Pivotal.Income
			pivotal.Payments += group.Pivotal.Payments
			pivotal.Total += group.Pivotal.Income - group.Pivotal.Payments
		}
	}
	return
}
