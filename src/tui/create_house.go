package tui

import (
	"fmt"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/rivo/tview"
)

const CreateHousePageName = "create-page"

type CreateHouse struct {
	*tview.Form
	*NavigationBack
	app     *TerminalApp
	request houseModel.CreateHouseRequest
}

func (c *CreateHouse) my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo {
	return NewNavigationInfo(CreateHousePageName, func() tview.Primitive { return NewCreateHouse(app, parent) })
}

func (c *CreateHouse) enrichNavigation(app *TerminalApp, parent *NavigationInfo) {
	c.MyNavigation = interface{}(c).(MyNavigation)
	c.enrich(app, parent)
}

func NewCreateHouse(app *TerminalApp, parent *NavigationInfo) *CreateHouse {
	f := &CreateHouse{
		Form:           tview.NewForm(),
		NavigationBack: NewNavigationBack(NewNavigation()),
		app:            app,
		request: houseModel.CreateHouseRequest{
			UserId: app.AuthorizedUser.Id,
		},
	}
	f.enrichNavigation(app, parent)

	f.
		AddInputField("Name", "", 20, nil, func(text string) { f.request.Name = text }).
		AddDropDown("Country", f.app.CountriesCodes, -1, func(option string, optionIndex int) { f.request.Country = option }).
		AddInputField("City", "", 20, nil, func(text string) { f.request.City = text }).
		AddInputField("Street Line 1", "", 20, nil, func(text string) { f.request.StreetLine1 = text }).
		AddInputField("Street Line 2", "", 20, nil, func(text string) { f.request.StreetLine2 = text }).
		AddButton("Create", func() {
			if houseResponse, err := f.app.GetHouseService().Add(f.request); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.app.House = houseResponse
				f.ShowInfoReturnHome(fmt.Sprintf("House %s successfully added.", houseResponse.Name))
			}
		})
	f.SetBorder(true).SetTitle("Add House").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}
