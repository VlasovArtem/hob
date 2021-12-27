package tui

import (
	"context"
	"fmt"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const UpdateHousePageName = "house-update-page"

type UpdateHouse struct {
	*FlexApp
	*Navigation
	app           *TerminalApp
	updateContent context.Context
}

func (u *UpdateHouse) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(UpdateHousePageName, func() tview.Primitive { return NewUpdateHouse(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdateHouse) enrichNavigation(app *TerminalApp, variables map[string]interface{}) {
	u.Navigation = NewNavigation(
		app,
		u.NavigationInfo(app, variables),
	)
}

func NewUpdateHouse(app *TerminalApp, houseId uuid.UUID) *UpdateHouse {
	f := &UpdateHouse{
		FlexApp: NewFlexApp(),
		app:     app,
	}
	f.enrichNavigation(app, map[string]interface{}{"id": houseId})
	f.InitFlexApp(app)
	f.bindKeys()

	houseDto, err := app.GetHouseService().FindById(houseId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	currentCountryCodeIndex := -1

	for i, code := range f.app.CountriesCodes {
		if code == houseDto.CountryCode {
			currentCountryCodeIndex = i
		}
	}

	request := houseModel.UpdateHouseRequest{}

	form := tview.NewForm().
		AddInputField("Name", houseDto.Name, 20, nil, func(text string) { request.Name = text }).
		AddDropDown("Country", f.app.CountriesCodes, currentCountryCodeIndex, func(option string, optionIndex int) { request.Country = option }).
		AddInputField("City", houseDto.City, 20, nil, func(text string) { request.City = text }).
		AddInputField("Street Line 1", houseDto.StreetLine1, 20, nil, func(text string) { request.StreetLine1 = text }).
		AddInputField("Street Line 2", houseDto.StreetLine2, 20, nil, func(text string) { request.StreetLine2 = text }).
		AddButton("Update", func() {
			if err := f.app.GetHouseService().Update(houseId, request); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.ShowInfoReturnBack(fmt.Sprintf("House %s successfully updated.", request.Name))
			}
		}).AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update a House").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateHouse) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}
