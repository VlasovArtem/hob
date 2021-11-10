package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"house/model"
	"test"
	"test/testhelper"
	"testing"
)

var countriesService = testhelper.InitCountryService()

func serviceGenerator() HouseService { return NewHouseService(countriesService) }

func TestAddHouse(t *testing.T) {
	houseService := serviceGenerator()

	type args struct {
		house model.CreateHouseRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "with new house",
			args: args{house: test.GenerateCreateHouseRequest()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err2, response := houseService.AddHouse(tt.args.house)

			assert.Nil(t, err2)

			err, house := houseService.FindById(response.Id)
			assert.Nil(t, err)
			assert.Equal(t, test.GenerateHouseResponse(house.Id, house.Name), house, "Houses should be the same")
		})
	}
}

func TestFindAllHouses(t *testing.T) {
	houseService := serviceGenerator()

	request := test.GenerateCreateHouseRequest()
	err, response := houseService.AddHouse(request)

	assert.Nil(t, err)

	tests := []struct {
		name            string
		wantResult      []model.HouseResponse
		serviceProvider func() HouseService
	}{
		{
			name:            "houses",
			wantResult:      []model.HouseResponse{response},
			serviceProvider: func() HouseService { return houseService },
		},
		{
			name:            "empty",
			wantResult:      []model.HouseResponse{},
			serviceProvider: serviceGenerator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.serviceProvider().FindAllHouses()

			assert.Equal(t, tt.wantResult, actual)
		})
	}
}

func TestFindById(t *testing.T) {
	houseService := serviceGenerator()

	request := test.GenerateCreateHouseRequest()
	err, response := houseService.AddHouse(request)

	assert.Nil(t, err)

	notExistingId := uuid.New()

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		wantError error
		wantModel model.HouseResponse
	}{
		{
			name: "with existing",
			args: args{
				id: response.Id,
			},
			wantModel: response,
			wantError: nil,
		}, {
			name: "with not existing",
			args: args{
				id: notExistingId,
			},
			wantModel: model.HouseResponse{},
			wantError: errors.New(fmt.Sprintf("House with id - %s not exists", notExistingId)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, house := houseService.FindById(tt.args.id)
			assert.Equalf(t, tt.wantError, err, "FindById(%v)", tt.args.id)
			assert.Equalf(t, tt.wantModel, house, "FindById(%v)", tt.args.id)
		})
	}
}
