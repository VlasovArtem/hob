package handler

import (
	"common/rest"
	"encoding/json"
	"net/http"
	"payment/model"
	ps "payment/service"
)

type paymentHandlerObject struct {
	paymentService ps.PaymentService
}

func NewPaymentHandler(paymentService ps.PaymentService) PaymentHandler {
	return &paymentHandlerObject{
		paymentService: paymentService,
	}
}

type PaymentHandler interface {
	AddPayment() http.HandlerFunc
	FindPaymentById() http.HandlerFunc
	FindPaymentByHouseId() http.HandlerFunc
	FindPaymentByUserId() http.HandlerFunc
}

func (p *paymentHandlerObject) AddPayment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		paymentRequest := model.CreatePaymentRequest{}

		if err := rest.PerformRequest(&paymentRequest, writer, request); err != nil {
			return
		}

		if payment, err := p.paymentService.AddPayment(paymentRequest); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(payment)
		}
	}
}

func (p *paymentHandlerObject) FindPaymentById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := p.paymentService.FindPaymentById(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}

func (p *paymentHandlerObject) FindPaymentByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentService.FindPaymentByHouseId(id), nil)
		}
	}
}

func (p *paymentHandlerObject) FindPaymentByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentService.FindPaymentByUserId(id), nil)
		}
	}
}
