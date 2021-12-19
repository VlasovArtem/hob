package tui

import (
	"github.com/VlasovArtem/hob/src/app"
	countries "github.com/VlasovArtem/hob/src/country/service"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	houses "github.com/VlasovArtem/hob/src/house/service"
	incomes "github.com/VlasovArtem/hob/src/income/service"
	payments "github.com/VlasovArtem/hob/src/payment/service"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	users "github.com/VlasovArtem/hob/src/user/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

type TerminalApp struct {
	*tview.Application
	Main           *tview.Pages
	root           *app.RootApplication
	AuthorizedUser *userModel.UserDto
	House          *houseModel.HouseDto
	CountriesCodes []string
	actions        KeyActions
}

func NewTApp(rootApplication *app.RootApplication) *TerminalApp {
	tapp := &TerminalApp{
		Application: tview.NewApplication(),
		root:        rootApplication,
		Main:        tview.NewPages(),
	}

	tapp.SetUpBorders()

	for _, country := range tapp.getCountryService().FindAllCountries() {
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
	return t.root.DependenciesFactory.FindRequiredByObject(houses.HouseServiceObject{}).(houses.HouseService)
}

func (t *TerminalApp) GetUserService() users.UserService {
	return t.root.DependenciesFactory.FindRequiredByObject(users.UserServiceObject{}).(users.UserService)
}

func (t *TerminalApp) getCountryService() countries.CountryService {
	return t.root.DependenciesFactory.FindRequiredByObject(countries.CountryServiceObject{}).(countries.CountryService)
}

func (t *TerminalApp) GetIncomeService() incomes.IncomeService {
	return t.root.DependenciesFactory.FindRequiredByObject(incomes.IncomeServiceObject{}).(incomes.IncomeService)
}

func (t *TerminalApp) GetPaymentService() payments.PaymentService {
	return t.root.DependenciesFactory.FindRequiredByObject(payments.PaymentServiceObject{}).(payments.PaymentService)
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
