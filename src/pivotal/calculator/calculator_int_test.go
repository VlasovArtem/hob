package calculator

import (
	"github.com/VlasovArtem/hob/src/db"
	groupMocks "github.com/VlasovArtem/hob/src/group/mocks"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	groupRepository "github.com/VlasovArtem/hob/src/group/repository"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	houseRepository "github.com/VlasovArtem/hob/src/house/repository"
	houseService "github.com/VlasovArtem/hob/src/house/service"
	incomeMocks "github.com/VlasovArtem/hob/src/income/mocks"
	incomeModel "github.com/VlasovArtem/hob/src/income/model"
	incomeRepository "github.com/VlasovArtem/hob/src/income/repository"
	incomeService "github.com/VlasovArtem/hob/src/income/service"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	paymentRepository "github.com/VlasovArtem/hob/src/payment/repository"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/VlasovArtem/hob/src/pivotal/model"
	pivotalRepository "github.com/VlasovArtem/hob/src/pivotal/repository"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
	providerRepository "github.com/VlasovArtem/hob/src/provider/repository"
	providerService "github.com/VlasovArtem/hob/src/provider/service"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/test/testhelper/integration"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	userRepository "github.com/VlasovArtem/hob/src/user/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CalculatorIntegrationTestSuite struct {
	integration.Suite[*PivotalCalculatorServiceStr]
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (c *CalculatorIntegrationTestSuite) SetupSuite() {
	c.InitSuite(func(service db.DatabaseService) *PivotalCalculatorServiceStr {
		users := userService.NewUserService(userRepository.NewUserRepository(service))
		groups := groupService.NewGroupService(users, groupRepository.NewGroupRepository(service))
		houses := houseService.NewHouseService(
			testhelper.InitCountryService(),
			users,
			houseRepository.NewHouseRepository(service),
			groups,
		)
		pivotals := pivotalService.NewPivotalService(
			pivotalRepository.NewHousePivotalRepository(service),
			pivotalRepository.NewGroupPivotalRepository(service),
			houses,
			groups,
		)
		return &PivotalCalculatorServiceStr{
			housePivotalRepository: pivotalRepository.NewHousePivotalRepository(service),
			groupPivotalRepository: pivotalRepository.NewGroupPivotalRepository(service),
			houseService:           houses,
			incomeService:          incomeService.NewIncomeService(houses, groups, incomeRepository.NewIncomeRepository(service), pivotals),
			paymentService:         paymentService.NewPaymentService(users, houses, providerService.NewProviderService(providerRepository.NewProviderRepository(service)), paymentRepository.NewPaymentRepository(service), pivotals),
		}
	})

	c.
		AddAfterTest(func(object *PivotalCalculatorServiceStr) {
			testhelper.TruncateTableCascade(c.DatabaseService, "house_groups")
			testhelper.TruncateTableCascade(c.DatabaseService, "income_groups")
			testhelper.TruncateTable(c.DatabaseService, groupModel.Group{})
			testhelper.TruncateTable(c.DatabaseService, paymentModel.Payment{})
			testhelper.TruncateTable(c.DatabaseService, incomeModel.Income{})
			testhelper.TruncateTable(c.DatabaseService, model.HousePivotal{})
			testhelper.TruncateTable(c.DatabaseService, model.GroupPivotal{})
		}).
		ExecuteMigration(userModel.User{}, groupModel.Group{}, houseModel.House{}, providerModel.Provider{}, paymentModel.Payment{}, incomeModel.Income{}, model.HousePivotal{}, model.GroupPivotal{})

	c.createdUser = userMocks.GenerateUser()
	c.CreateEntity(&c.createdUser)

	c.createdHouse = houseMocks.GenerateHouse(c.createdUser.Id)
	c.CreateEntity(&c.createdHouse)
}

func TestCalculatorIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CalculatorIntegrationTestSuite))
}

func (c *CalculatorIntegrationTestSuite) TestCalculate_WithoutGroup() {
	payment := paymentMocks.GeneratePaymentWithoutProvider(c.createdHouse.Id, c.createdUser.Id)
	c.CreateEntity(&payment)
	income := incomeMocks.GenerateIncome(&c.createdHouse.Id)
	c.CreateEntity(&income)

	pivotal, err := c.O.Calculate(c.createdHouse.Id)

	assert.Nil(c.T(), err)

	assert.True(c.T(), len(pivotal.Groups) == 0)

	paymentsSum := c.sumPayments(payment)
	incomesSum := c.sumIncomes(income)

	assert.Equal(c.T(), model.HousePivotalDto{
		Pivotal: model.Pivotal{
			Id:                      pivotal.House.Id,
			Income:                  incomesSum,
			Payments:                paymentsSum,
			Total:                   incomesSum - paymentsSum,
			LatestIncomeUpdateDate:  income.Date.Add(1 * time.Microsecond),
			LatestPaymentUpdateDate: payment.Date.Add(1 * time.Microsecond),
		},
		HouseId: c.createdHouse.Id,
	}, pivotal.House)

	assert.Equal(c.T(), model.TotalPivotalDto{
		Income:   incomesSum,
		Payments: paymentsSum,
		Total:    incomesSum - paymentsSum,
	}, pivotal.Total)

	assert.True(c.T(), c.O.housePivotalRepository.Exists(pivotal.House.Id))
}

func (c *CalculatorIntegrationTestSuite) TestCalculate_WithoutGroupAndMultipleData() {
	firstPayment := paymentMocks.GeneratePaymentWithoutProvider(c.createdHouse.Id, c.createdUser.Id)
	c.CreateEntity(&firstPayment)
	firstIncome := incomeMocks.GenerateIncome(&c.createdHouse.Id)
	c.CreateEntity(&firstIncome)
	secondPayment := paymentMocks.GeneratePaymentWithoutProvider(c.createdHouse.Id, c.createdUser.Id)
	c.CreateEntity(&secondPayment)
	secondIncome := incomeMocks.GenerateIncome(&c.createdHouse.Id)
	c.CreateEntity(&secondIncome)

	pivotal, err := c.O.Calculate(c.createdHouse.Id)

	assert.Nil(c.T(), err)

	assert.True(c.T(), len(pivotal.Groups) == 0)

	paymentsSum := c.sumPayments(firstPayment, secondPayment)
	incomesSum := c.sumIncomes(firstIncome, secondIncome)

	assert.Equal(c.T(), model.HousePivotalDto{
		Pivotal: model.Pivotal{
			Id:                      pivotal.House.Id,
			Income:                  incomesSum,
			Payments:                paymentsSum,
			Total:                   incomesSum - paymentsSum,
			LatestIncomeUpdateDate:  secondIncome.Date.Add(1 * time.Microsecond),
			LatestPaymentUpdateDate: secondPayment.Date.Add(1 * time.Microsecond),
		},
		HouseId: c.createdHouse.Id,
	}, pivotal.House)

	assert.Equal(c.T(), model.TotalPivotalDto{
		Income:   incomesSum,
		Payments: paymentsSum,
		Total:    incomesSum - paymentsSum,
	}, pivotal.Total)

	assert.True(c.T(), c.O.housePivotalRepository.Exists(pivotal.House.Id))
}

func (c *CalculatorIntegrationTestSuite) TestCalculate() {
	group := groupMocks.GenerateGroup(c.createdUser.Id)
	c.CreateEntity(&group)

	firstHouse := houseMocks.GenerateHouse(c.createdUser.Id)
	firstHouse.Groups = append(firstHouse.Groups, group)
	c.CreateEntity(&firstHouse)

	secondHouse := houseMocks.GenerateHouse(c.createdUser.Id)
	secondHouse.Groups = append(secondHouse.Groups, group)
	c.CreateEntity(&secondHouse)

	firstPayment := paymentMocks.GeneratePaymentWithoutProvider(firstHouse.Id, c.createdUser.Id)
	c.CreateEntity(&firstPayment)

	firstIncome := incomeMocks.GenerateIncome(nil)
	firstIncome.Groups = append(firstIncome.Groups, group)
	c.CreateEntity(&firstIncome)

	secondPayment := paymentMocks.GeneratePaymentWithoutProvider(secondHouse.Id, c.createdUser.Id)
	c.CreateEntity(&secondPayment)

	secondIncome := incomeMocks.GenerateIncome(nil)
	secondIncome.Groups = append(secondIncome.Groups, group)
	c.CreateEntity(&secondIncome)

	firstPivotal, err := c.O.Calculate(firstHouse.Id)

	assert.Nil(c.T(), err)

	secondPivotal, err := c.O.Calculate(secondHouse.Id)

	assert.Nil(c.T(), err)

	assert.True(c.T(), len(firstPivotal.Groups) == 1)
	assert.True(c.T(), len(secondPivotal.Groups) == 1)

	paymentsSum := c.sumPayments(firstPayment, secondPayment)
	incomesSum := c.sumIncomes(firstIncome, secondIncome)

	assert.Equal(c.T(), model.HousePivotalDto{
		Pivotal: model.Pivotal{
			Id:                      firstPivotal.House.Id,
			Income:                  incomesSum,
			Payments:                float64(firstPayment.Sum),
			Total:                   incomesSum - float64(firstPayment.Sum),
			LatestIncomeUpdateDate:  secondIncome.Date.Add(1 * time.Microsecond),
			LatestPaymentUpdateDate: secondPayment.Date.Add(1 * time.Microsecond),
		},
		HouseId: firstHouse.Id,
	}, firstPivotal.House)

	assert.Equal(c.T(), model.HousePivotalDto{
		Pivotal: model.Pivotal{
			Id:                      secondPivotal.House.Id,
			Income:                  incomesSum,
			Payments:                float64(secondPayment.Sum),
			Total:                   incomesSum - float64(secondPayment.Sum),
			LatestIncomeUpdateDate:  secondIncome.Date.Add(1 * time.Microsecond),
			LatestPaymentUpdateDate: secondPayment.Date.Add(1 * time.Microsecond),
		},
		HouseId: secondHouse.Id,
	}, secondPivotal.House)

	assert.Equal(c.T(), []model.GroupPivotalDto{
		{
			Pivotal: model.Pivotal{
				Id:                      firstPivotal.Groups[0].Id,
				Income:                  incomesSum,
				Payments:                paymentsSum,
				Total:                   incomesSum - paymentsSum,
				LatestIncomeUpdateDate:  secondIncome.Date.Add(1 * time.Microsecond),
				LatestPaymentUpdateDate: secondPayment.Date.Add(1 * time.Microsecond),
			},
			GroupId: group.Id,
		},
	}, firstPivotal.Groups)

	assert.Equal(c.T(), []model.GroupPivotalDto{
		{
			Pivotal: model.Pivotal{
				Id:                      secondPivotal.Groups[0].Id,
				Income:                  incomesSum,
				Payments:                paymentsSum,
				Total:                   incomesSum - paymentsSum,
				LatestIncomeUpdateDate:  secondIncome.Date.Add(1 * time.Microsecond),
				LatestPaymentUpdateDate: secondPayment.Date.Add(1 * time.Microsecond),
			},
			GroupId: group.Id,
		},
	}, secondPivotal.Groups)

	assert.Equal(c.T(), model.TotalPivotalDto{
		Income:   incomesSum,
		Payments: paymentsSum,
		Total:    incomesSum - paymentsSum,
	}, firstPivotal.Total)

	assert.True(c.T(), c.O.housePivotalRepository.Exists(firstPivotal.House.Id))
	assert.True(c.T(), c.O.housePivotalRepository.Exists(secondPivotal.House.Id))

	assert.Equal(c.T(), firstPivotal.Total, secondPivotal.Total)
	assert.NotEqual(c.T(), firstPivotal.House, secondPivotal.House)
	assert.NotEqual(c.T(), firstPivotal.Groups, secondPivotal.Groups)
}

func (c *CalculatorIntegrationTestSuite) sumPayments(payments ...paymentModel.Payment) (sum float64) {
	for _, payment := range payments {
		sum += float64(payment.Sum)
	}
	return
}

func (c *CalculatorIntegrationTestSuite) sumIncomes(incomes ...incomeModel.Income) (sum float64) {
	for _, income := range incomes {
		sum += float64(income.Sum)
	}
	return
}
