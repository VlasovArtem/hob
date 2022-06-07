package service

import (
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
	"time"
)

type PivotalServiceObject struct {
	housePivotalRepository repository.PivotalRepository
	groupPivotalRepository repository.PivotalRepository
	houseService           houseService.HouseService
	groupService           groupService.GroupService
	pivotalCache           cache.PivotalCache
	incomeService          incomeService.IncomeService
	paymentService         paymentService.PaymentService
}

func NewPivotalService(
	housePivotalRepository repository.PivotalRepository,
	groupPivotalRepository repository.PivotalRepository,
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

func (p PivotalServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalService(
		dependency.FindRequiredDependency[repository.HousePivotalRepository, repository.PivotalRepository](factory),
		dependency.FindRequiredDependency[repository.GroupPivotalRepository, repository.PivotalRepository](factory),
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
}

func (p *PivotalServiceObject) FindByHouseId(houseId uuid.UUID) (response model.HousePivotalDto, err error) {
	if !p.houseService.ExistsById(houseId) {
		return model.HousePivotalDto{}, interrors.NewErrNotFound("House with id %s not found", houseId)
	}

	var pivotal model.HousePivotal

	if err = p.housePivotalRepository.FindBySourceId(houseId, &pivotal); err != nil {
		return
	}

	if err = p.findById(houseId, &pivotal.Pivotal,
		func(from *time.Time) float64 {
			return p.incomeService.CalculateSumByHouseId(houseId, from)
		},
		func(from *time.Time) float64 {
			return p.paymentService.CalculateSum([]uuid.UUID{houseId}, from)
		},
	); err != nil {
		return
	}

	pivotal.HouseId = houseId
	return pivotal.ToDto(), nil
}

func (p *PivotalServiceObject) FindByGroupId(groupId uuid.UUID) (response model.GroupPivotalDto, err error) {
	if !p.groupService.ExistsById(groupId) {
		return model.GroupPivotalDto{}, interrors.NewErrNotFound("Group with id %s not found", groupId)
	}

	var pivotal model.GroupPivotal

	if err = p.groupPivotalRepository.FindBySourceId(groupId, &pivotal); err != nil {
		return
	}

	if err = p.findById(groupId, &pivotal.Pivotal,
		func(from *time.Time) float64 {
			return p.incomeService.CalculateSumByGroupId(groupId, from)
		},
		func(from *time.Time) float64 {
			houseIds := common.MapSlice(p.houseService.FindHousesByGroupId(groupId), func(house houseModel.HouseDto) uuid.UUID { return house.Id })

			return p.paymentService.CalculateSum(houseIds, from)
		}); err != nil {
		return model.GroupPivotalDto{}, err
	}

	pivotal.GroupId = groupId
	return pivotal.ToDto(), nil
}

func (p *PivotalServiceObject) findById(
	id uuid.UUID,
	pivotal *model.Pivotal,
	incomeSumProvider func(from *time.Time) float64,
	paymentSumProvider func(from *time.Time) float64,
) (err error) {
	if pivotal = p.pivotalCache.Find(id); pivotal != nil {
		return
	}

	var latestIncomesDate, latestPaymentDate *time.Time

	if pivotal.Id != uuid.Nil {
		latestIncomesDate = &pivotal.LatestIncomeUpdateDate
		latestPaymentDate = &pivotal.LatestPaymentUpdateDate

		pivotal.LatestIncomeUpdateDate = pivotal.LatestIncomeUpdateDate.Add(time.Millisecond * 1)
		pivotal.LatestPaymentUpdateDate = pivotal.LatestPaymentUpdateDate.Add(time.Millisecond * 1)
	}

	incomeSum := incomeSumProvider(latestIncomesDate)
	paymentSum := paymentSumProvider(latestPaymentDate)

	pivotal.Total = incomeSum - paymentSum
	pivotal.Income = incomeSum
	pivotal.Payments = paymentSum

	p.pivotalCache.Add(id, pivotal)

	return
}
