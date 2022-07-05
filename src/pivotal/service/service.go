package service

import (
	"errors"
	"github.com/VlasovArtem/hob/src/common/dependency"
	interrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/transactional"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	"github.com/VlasovArtem/hob/src/pivotal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PivotalServiceStr struct {
	housePivotalRepository repository.PivotalRepository[model.HousePivotal]
	groupPivotalRepository repository.PivotalRepository[model.GroupPivotal]
	houseService           houseService.HouseService
	groupService           groupService.GroupService
}

func NewPivotalService(
	housePivotalRepository repository.PivotalRepository[model.HousePivotal],
	groupPivotalRepository repository.PivotalRepository[model.GroupPivotal],
	houseService houseService.HouseService,
	groupService groupService.GroupService,
) PivotalService {
	return &PivotalServiceStr{
		housePivotalRepository: housePivotalRepository,
		groupPivotalRepository: groupPivotalRepository,
		houseService:           houseService,
		groupService:           groupService,
	}
}

func (p *PivotalServiceStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalService(
		dependency.FindRequiredDependency[repository.HousePivotalRepository, repository.PivotalRepository[model.HousePivotal]](factory),
		dependency.FindRequiredDependency[repository.GroupPivotalRepository, repository.PivotalRepository[model.GroupPivotal]](factory),
		dependency.FindRequiredDependency[houseService.HouseServiceStr, houseService.HouseService](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceStr, groupService.GroupService](factory),
	)
}

type PivotalService interface {
	transactional.Transactional[PivotalService]
	Find(houseId uuid.UUID) (model.PivotalResponseDto, error)
	AddPayment(paymentSum float64, updateTime time.Time, houseId uuid.UUID) error
	AddIncome(incomeSum float64, updateTime time.Time, houseId *uuid.UUID, groupIds []uuid.UUID) error
	ExistsByHouseId(houseId uuid.UUID) bool
}

func (p *PivotalServiceStr) AddPayment(paymentSum float64, updateTime time.Time, houseId uuid.UUID) error {
	return p.housePivotalRepository.DB().Transaction(func(db *gorm.DB) error {
		trx := p.Transactional(db).(*PivotalServiceStr)

		if houseDto, err := trx.houseService.FindById(houseId); err != nil {
			return err
		} else {
			for _, group := range houseDto.Groups {
				var groupPivotal model.GroupPivotal
				if err = trx.groupPivotalRepository.FindBySourceId(group.Id, &groupPivotal); err != nil && !errors.Is(err, interrors.ErrNotFound{}) {
					return err
				} else if groupPivotal.Pivotal.Id != uuid.Nil {
					groupPivotal.Pivotal.Payments += paymentSum
					groupPivotal.Pivotal.Total -= paymentSum
					groupPivotal.Pivotal.LatestPaymentUpdateDate = updateTime
					if err := trx.groupPivotalRepository.Update(groupPivotal.Pivotal.Id, groupPivotal.Pivotal); err != nil {
						return err
					}
				}
			}
			var housePivotal model.HousePivotal
			if err = trx.housePivotalRepository.FindBySourceId(houseId, &housePivotal); err != nil {
				return err
			}
			housePivotal.Pivotal.Payments += paymentSum
			housePivotal.Pivotal.Total -= paymentSum
			housePivotal.Pivotal.LatestPaymentUpdateDate = updateTime
			if err = trx.housePivotalRepository.Update(housePivotal.Pivotal.Id, housePivotal.Pivotal); err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PivotalServiceStr) AddIncome(incomeSum float64, updateTime time.Time, houseId *uuid.UUID, groupIds []uuid.UUID) error {
	return p.housePivotalRepository.DB().Transaction(func(db *gorm.DB) error {
		trx := p.Transactional(db).(*PivotalServiceStr)

		if houseId == nil && len(groupIds) == 0 {
			return nil
		}

		calcHousePivotal := func(houseId uuid.UUID) error {
			var housePivotal model.HousePivotal
			if err := trx.housePivotalRepository.FindBySourceId(houseId, &housePivotal); err != nil {
				return err
			}
			housePivotal.Pivotal.Income += incomeSum
			housePivotal.Pivotal.Total += incomeSum
			housePivotal.Pivotal.LatestIncomeUpdateDate = updateTime
			if err := p.housePivotalRepository.Update(housePivotal.Pivotal.Id, housePivotal.Pivotal); err != nil {
				return err
			}
			return nil
		}

		if len(groupIds) == 0 && houseId != nil {
			return calcHousePivotal(*houseId)
		} else {
			for _, dto := range trx.houseService.FindHousesByGroupIds(groupIds) {
				if err := calcHousePivotal(dto.Id); err != nil {
					return err
				}
			}
			for _, groupId := range groupIds {
				var groupPivotal model.GroupPivotal
				if err := trx.groupPivotalRepository.FindBySourceId(groupId, &groupPivotal); err != nil && !errors.Is(err, interrors.ErrNotFound{}) {
					return err
				} else if groupPivotal.Pivotal.Id != uuid.Nil {
					groupPivotal.Pivotal.Income += incomeSum
					groupPivotal.Pivotal.Total += incomeSum
					groupPivotal.Pivotal.LatestPaymentUpdateDate = updateTime
					if err := trx.groupPivotalRepository.Update(groupPivotal.Pivotal.Id, groupPivotal.Pivotal); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (p *PivotalServiceStr) Transactional(tx *gorm.DB) PivotalService {
	return &PivotalServiceStr{
		housePivotalRepository: p.housePivotalRepository.Transactional(tx),
		groupPivotalRepository: p.groupPivotalRepository.Transactional(tx),
		houseService:           p.houseService.Transactional(tx),
		groupService:           p.groupService.Transactional(tx),
	}
}

func (p *PivotalServiceStr) Find(houseId uuid.UUID) (response model.PivotalResponseDto, err error) {
	houseDto, err := p.houseService.FindById(houseId)
	if err != nil {
		return
	}

	var groupPivotals []model.GroupPivotal

	for _, group := range houseDto.Groups {
		var groupPivotal model.GroupPivotal

		err = p.groupPivotalRepository.FindBySourceId(group.Id, &groupPivotal)
		if err != nil {
			return response, err
		}
		if groupPivotal.Pivotal.Id != uuid.Nil {
			response.Groups = append(response.Groups, groupPivotal.ToDto())
			groupPivotals = append(groupPivotals, groupPivotal)
		} else {
			return response, errors.New("pivotal not correct")
		}
	}

	var housePivotal model.HousePivotal

	if err = p.housePivotalRepository.FindBySourceId(houseId, &housePivotal); err != nil {
		return
	}

	if housePivotal.Pivotal.Id != uuid.Nil {
		response.House = housePivotal.ToDto()
	} else {
		return response, errors.New("pivotal not correct")
	}

	response.Total = calculateTotal(housePivotal, groupPivotals)

	return
}

func (p *PivotalServiceStr) FindByGroupId(groupId uuid.UUID) (response model.GroupPivotalDto, err error) {
	if !p.groupService.ExistsById(groupId) {
		return model.GroupPivotalDto{}, interrors.NewErrNotFound("Group with id %s not found", groupId)
	}

	var pivotal model.GroupPivotal

	err = p.groupPivotalRepository.FindBySourceId(groupId, &pivotal)

	return
}

func (p *PivotalServiceStr) ExistsByHouseId(houseId uuid.UUID) bool {
	return p.housePivotalRepository.ExistsBy("house_id = ?", houseId)
}

func calculateTotal(house model.HousePivotal, groups []model.GroupPivotal) (pivotal model.TotalPivotalDto) {
	if house.Pivotal.Id != uuid.Nil {
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
