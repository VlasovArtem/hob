package tui

import (
	"fmt"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const CreateHousePageName = "create-house"

type CreateHouse struct {
	*FlexApp
	*Navigation
	app     *TerminalApp
	request houseModel.CreateHouseRequest
}

func (c *CreateHouse) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(CreateHousePageName, func() tview.Primitive { return NewCreateHouse(app) })
}

func (c *CreateHouse) enrichNavigation(app *TerminalApp) {
	c.Navigation = NewNavigation(app, c.NavigationInfo(app, nil))
}

func NewCreateHouse(app *TerminalApp) *CreateHouse {
	f := &CreateHouse{
		FlexApp: NewFlexApp(),
		app:     app,
		request: houseModel.CreateHouseRequest{
			UserId: app.AuthorizedUser.Id,
		},
	}
	f.enrichNavigation(app)
	f.bindKeys()

	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) { f.request.Name = text }).
		AddDropDown("CountryCode", f.app.CountriesCodes, -1, func(option string, optionIndex int) { f.request.CountryCode = option }).
		AddInputField("City", "", 20, nil, func(text string) { f.request.City = text }).
		AddInputField("Street Line 1", "", 20, nil, func(text string) { f.request.StreetLine1 = text }).
		AddInputField("Street Line 2", "", 20, nil, func(text string) { f.request.StreetLine2 = text }).
		AddButton("Create", func() {
			if houseResponse, err := f.app.GetHouseService().Add(f.request); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.app.House = &houseResponse
				f.ShowInfoReturnHome(fmt.Sprintf("House %s successfully added.", houseResponse.Name))
			}
		})
	form.SetBorder(true).SetTitle("Add House").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreateHouse) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.KeyBack),
	}
}
