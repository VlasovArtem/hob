package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"house/model"
)

var houses = make(map[uuid.UUID]model.House)

func AddHouse(house model.House) {
	houses[house.Id] = house
}

func FindAllHouses() []model.House {
	result := []model.House{}

	for _, house := range houses {
		result = append(result, house)
	}

	return result
}

func FindById(id uuid.UUID) (error, model.House) {
	if house, ok := houses[id]; ok {
		return nil, house
	}

	return errors.New(fmt.Sprintf("House with id - %s not exists", id)), model.DEFAULT
}

func deleteAll() {
	houses = make(map[uuid.UUID]model.House)
}
