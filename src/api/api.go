package api

import (
	"github.com/VlasovArtem/hob/src/app"
	"github.com/VlasovArtem/hob/src/common/dependency"
	countryHandler "github.com/VlasovArtem/hob/src/country/handler"
	healthHandler "github.com/VlasovArtem/hob/src/health/handler"
	houseHandler "github.com/VlasovArtem/hob/src/house/handler"
	incomeHandler "github.com/VlasovArtem/hob/src/income/handler"
	incomeSchedulerHandler "github.com/VlasovArtem/hob/src/income/scheduler/handler"
	meterHandler "github.com/VlasovArtem/hob/src/meter/handler"
	paymentHandler "github.com/VlasovArtem/hob/src/payment/handler"
	paymentSchedulerHandler "github.com/VlasovArtem/hob/src/payment/scheduler/handler"
	providerHandler "github.com/VlasovArtem/hob/src/provider/handler"
	userHandler "github.com/VlasovArtem/hob/src/user/handler"
	"github.com/gorilla/mux"
)

type ApplicationHandler interface {
	dependency.ObjectDependencyInitializer
	Init(*mux.Router)
}

func InitApi(router *mux.Router, application *app.RootApplication) {
	addHandler(router, application, new(countryHandler.CountryHandlerObject))
	addHandler(router, application, new(userHandler.UserHandlerObject))
	addHandler(router, application, new(houseHandler.HouseHandlerObject))
	addHandler(router, application, new(providerHandler.ProviderHandlerObject))
	addHandler(router, application, new(paymentHandler.PaymentHandlerObject))
	addHandler(router, application, new(paymentSchedulerHandler.PaymentSchedulerHandlerObject))
	addHandler(router, application, new(meterHandler.MeterHandlerObject))
	addHandler(router, application, new(incomeHandler.IncomeHandlerObject))
	addHandler(router, application, new(incomeSchedulerHandler.IncomeSchedulerHandlerObject))
	addHandler(router, application, new(healthHandler.HealthHandlerObject))
}

func addHandler(router *mux.Router, application *app.RootApplication, handler ApplicationHandler) {
	application.DependenciesFactory.AddAutoDependency(handler).(ApplicationHandler).Init(router)
}
