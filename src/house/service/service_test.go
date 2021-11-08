package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"house/model"
	"reflect"
	"testing"
)

func generateHouse(id uuid.UUID) model.House {
	return model.House{
		Id:          id,
		Name:        "Test House",
		Country:     "Country",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		Deleted:     false,
	}
}

func TestAddHouse(t *testing.T) {
	t.Cleanup(func() { deleteAll() })

	type args struct {
		house model.House
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "with new house",
			args: args{house: generateHouse(uuid.New())},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddHouse(tt.args.house)

			err, house := FindById(tt.args.house.Id)
			assert.Nil(t, err)
			assert.Equal(t, tt.args.house, house, "Houses should be the same")
		})
	}
}

func TestFindAllHouses(t *testing.T) {
	t.Cleanup(func() { deleteAll() })

	house := generateHouse(uuid.New())

	tests := []struct {
		name       string
		wantResult []model.House
		prepare    func()
	}{
		{
			name:       "houses",
			wantResult: []model.House{house},
			prepare: func() {
				AddHouse(house)
			},
		},
		{
			name:       "empty",
			wantResult: []model.House{},
			prepare: func() {
				deleteAll()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			allHouses := FindAllHouses()

			if len(tt.wantResult) == 0 && len(allHouses) != 0 {
				t.Errorf("FindAllHouses() should be empty")
			} else if len(tt.wantResult) != 0 {
				if !reflect.DeepEqual(allHouses, tt.wantResult) {
					t.Errorf("FindAllHouses() = %v, want %v", allHouses, tt.wantResult)
				}
			}
		})
	}
}

func TestFindById(t *testing.T) {
	t.Cleanup(func() { deleteAll() })
	house := generateHouse(uuid.New())
	AddHouse(house)

	notExistingId := uuid.New()

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		wantError error
		wantModel model.House
	}{
		{
			name: "with existing",
			args: args{
				id: house.Id,
			},
			wantModel: house,
			wantError: nil,
		}, {
			name: "with not existing",
			args: args{
				id: notExistingId,
			},
			wantModel: model.DEFAULT,
			wantError: errors.New(fmt.Sprintf("House with id - %s not exists", notExistingId)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, house := FindById(tt.args.id)
			assert.Equalf(t, tt.wantError, err, "FindById(%v)", tt.args.id)
			assert.Equalf(t, tt.wantModel, house, "FindById(%v)", tt.args.id)
		})
	}
}
