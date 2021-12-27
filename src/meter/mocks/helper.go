package mocks

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GenerateMeter(paymentId uuid.UUID, houseId uuid.UUID) model.Meter {
	details := map[string]float64{
		"first":  1.1,
		"second": 2.2,
	}
	if detailsBytes, err := json.Marshal(details); err != nil {
		log.Fatal().Err(err)
		return model.Meter{}
	} else {
		return model.Meter{
			Id:          uuid.New(),
			Name:        "Name",
			Details:     detailsBytes,
			Description: "Description",
			PaymentId:   paymentId,
			HouseId:     houseId,
		}
	}
}

func GenerateCreateMeterRequest() model.CreateMeterRequest {
	return model.CreateMeterRequest{
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   uuid.New(),
		HouseId:     uuid.New(),
	}
}

func GenerateUpdateMeterRequest() (uuid.UUID, model.UpdateMeterRequest) {
	return uuid.New(), model.UpdateMeterRequest{
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
			"third":  3.0,
		},
		Description: "Description",
	}
}

func GenerateMeterResponse(id uuid.UUID) model.MeterDto {
	return model.MeterDto{
		Id:   id,
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   uuid.New(),
		HouseId:     uuid.New(),
	}
}
