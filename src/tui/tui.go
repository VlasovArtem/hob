package tui

import (
	"github.com/VlasovArtem/hob/src/app"
	"github.com/VlasovArtem/hob/src/common/dependency"
	countries "github.com/VlasovArtem/hob/src/country/service"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	houses "github.com/VlasovArtem/hob/src/house/service"
	incomeSchedulers "github.com/VlasovArtem/hob/src/income/scheduler/service"
	incomes "github.com/VlasovArtem/hob/src/income/service"
	meters "github.com/VlasovArtem/hob/src/meter/service"
	paymentSchedulers "github.com/VlasovArtem/hob/src/payment/scheduler/service"
	payments "github.com/VlasovArtem/hob/src/payment/service"
	providers "github.com/VlasovArtem/hob/src/provider/service"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	users "github.com/VlasovArtem/hob/src/user/service"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"os"
)

var DefaultInputFieldWidth = 60
var DefaultUUID = uuid.UUID{}

type TerminalApp struct {
	*tview.Application
	Main           *tview.Pages
	root           *app.RootApplication
	AuthorizedUser *userModel.UserDto
	House          *houseModel.HouseDto
	Countries      map[string]string
	CountriesCodes []string
	CountriesNames []string
	actions        KeyActions
}

func NewTApp(rootApplication *app.RootApplication) *TerminalApp {
	tapp := &TerminalApp{
		Application: tview.NewApplication(),
		root:        rootApplication,
		Main:        tview.NewPages(),
	}

	tapp.SetUpBorders()

	allCountries := tapp.getCountryService().FindAllCountries()
	tapp.Countries = make(map[string]string, len(allCountries))
	for _, country := range allCountries {
		tapp.Countries[country.Code] = country.Name
		tapp.CountriesNames = append(tapp.CountriesNames, country.Name)
		tapp.CountriesCodes = append(tapp.CountriesCodes, country.Code)
	}

	return tapp
}

func (t *TerminalApp) Init() {
	log.Info().Msg("Initializing Terminal UI")

	config := t.root.Config
	if config != nil && config.User.Email != "" {
		user := config.User
		userResponse, err := t.GetUserService().VerifyUser(user.Email, user.Password)

		if err != nil {
			log.Error().Err(err).Msg("XDG user configuration is not valid")
		} else {
			t.AuthorizedUser = &userResponse
		}
	}

	if t.AuthorizedUser == nil {
		t.Main.AddAndSwitchToPage(SignInPageName, NewSignIn(t), true)
	} else {
		t.Main.AddAndSwitchToPage(HomePageName, NewHome(t), true)
	}

	t.SetRoot(t.Main, true).EnableMouse(true)
}

func (t *TerminalApp) GetHouseService() houses.HouseService {
	return dependency.FindRequiredDependency[houses.HouseServiceObject, houses.HouseService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetUserService() users.UserService {
	return dependency.FindRequiredDependency[users.UserServiceObject, users.UserService](t.root.DependenciesFactory)
}

func (t *TerminalApp) getCountryService() countries.CountryService {
	return dependency.FindRequiredDependency[countries.CountryServiceObject, countries.CountryService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetIncomeService() incomes.IncomeService {
	return dependency.FindRequiredDependency[incomes.IncomeServiceObject, incomes.IncomeService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetPaymentService() payments.PaymentService {
	return dependency.FindRequiredDependency[payments.PaymentServiceObject, payments.PaymentService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetProviderService() providers.ProviderService {
	return dependency.FindRequiredDependency[providers.ProviderServiceObject, providers.ProviderService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetMeterService() meters.MeterService {
	return dependency.FindRequiredDependency[meters.MeterServiceObject, meters.MeterService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetPaymentSchedulerService() paymentSchedulers.PaymentSchedulerService {
	return dependency.FindRequiredDependency[paymentSchedulers.PaymentSchedulerServiceObject, paymentSchedulers.PaymentSchedulerService](t.root.DependenciesFactory)
}

func (t *TerminalApp) GetIncomeSchedulerService() incomeSchedulers.IncomeSchedulerService {
	return dependency.FindRequiredDependency[incomeSchedulers.IncomeSchedulerServiceObject, incomeSchedulers.IncomeSchedulerService](t.root.DependenciesFactory)
}

func AsKey(evt *tcell.EventKey) tcell.Key {
	if evt.Key() != tcell.KeyRune {
		return evt.Key()
	}
	key := tcell.Key(evt.Rune())
	if evt.Modifiers() == tcell.ModAlt {
		key = tcell.Key(int16(evt.Rune()) * int16(evt.Modifiers()))
	}
	return key
}

func (t *TerminalApp) CreateErrAndToParent(pageName string, err error, doneFunc func(key tcell.Key)) {
	t.Main.AddAndSwitchToPage(pageName, NewInfoWithError(err, doneFunc), true)
}

func (t *TerminalApp) AddToPage(pageName string, primitive tview.Primitive) {
	t.Main.AddAndSwitchToPage(pageName, primitive, true)
}

func (t *TerminalApp) SetUpBorders() {
	tview.Borders = struct {
		Horizontal       rune
		Vertical         rune
		TopLeft          rune
		TopRight         rune
		BottomLeft       rune
		BottomRight      rune
		LeftT            rune
		RightT           rune
		TopT             rune
		BottomT          rune
		Cross            rune
		HorizontalFocus  rune
		VerticalFocus    rune
		TopLeftFocus     rune
		TopRightFocus    rune
		BottomLeftFocus  rune
		BottomRightFocus rune
	}{
		Horizontal:  tview.BoxDrawingsLightHorizontal,
		Vertical:    tview.BoxDrawingsLightVertical,
		TopLeft:     tview.BoxDrawingsLightDownAndRight,
		TopRight:    tview.BoxDrawingsLightDownAndLeft,
		BottomLeft:  tview.BoxDrawingsLightUpAndRight,
		BottomRight: tview.BoxDrawingsLightUpAndLeft,

		LeftT:   tview.BoxDrawingsLightVerticalAndRight,
		RightT:  tview.BoxDrawingsLightVerticalAndLeft,
		TopT:    tview.BoxDrawingsLightDownAndHorizontal,
		BottomT: tview.BoxDrawingsLightUpAndHorizontal,
		Cross:   tview.BoxDrawingsLightVerticalAndHorizontal,

		HorizontalFocus:  tview.BoxDrawingsLightHorizontal,
		VerticalFocus:    tview.BoxDrawingsLightVertical,
		TopLeftFocus:     tview.BoxDrawingsLightDownAndRight,
		TopRightFocus:    tview.BoxDrawingsLightDownAndLeft,
		BottomLeftFocus:  tview.BoxDrawingsLightUpAndRight,
		BottomRightFocus: tview.BoxDrawingsLightUpAndLeft,
	}
}

func (t *TerminalApp) Quit(key *tcell.EventKey) *tcell.EventKey {
	t.Stop()
	os.Exit(0)
	return key
}
