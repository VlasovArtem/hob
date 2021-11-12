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
			response, err2 := houseService.Add(tt.args.house)

			assert.Nil(t, err2)

			house, err := houseService.FindById(response.Id)
			assert.Nil(t, err)
			assert.Equal(t, test.GenerateHouseResponse(house.Id, house.Name), house, "Houses should be the same")
		})
	}
}

func TestFindAllHouses(t *testing.T) {
	houseService := serviceGenerator()

	request := test.GenerateCreateHouseRequest()
	response, err := houseService.Add(request)

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
			actual := tt.serviceProvider().FindAll()

			assert.Equal(t, tt.wantResult, actual)
		})
	}
}

func TestFindById(t *testing.T) {
	houseService := serviceGenerator()

	request := test.GenerateCreateHouseRequest()
	response, err := houseService.Add(request)

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
			house, err := houseService.FindById(tt.args.id)
			assert.Equalf(t, tt.wantError, err, "FindById(%v)", tt.args.id)
			assert.Equalf(t, tt.wantModel, house, "FindById(%v)", tt.args.id)
		})
	}
}

func Test_ExistsById(t *testing.T) {
	houseService := serviceGenerator()

	request := test.GenerateCreateHouseRequest()
	response, err := houseService.Add(request)

	assert.Nil(t, err)

	notExistingId := uuid.New()

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with existing",
			args: args{
				id: response.Id,
			},
			want: true,
		}, {
			name: "with not existing",
			args: args{
				id: notExistingId,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := houseService.ExistsById(tt.args.id)
			assert.Equalf(t, tt.want, exists, "FindById(%v)", tt.args.id)
		})
	}
}
