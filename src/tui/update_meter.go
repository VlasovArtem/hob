package tui

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"strings"
)

const UpdateMeterPageName = "meter-update-page"

type updateMeterReq struct {
	name, description, details string
}

type UpdateMeter struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (u *UpdateMeter) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(UpdateMeterPageName, func() tview.Primitive { return NewUpdateMeter(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdateMeter) enrichNavigation(app *TerminalApp, meterId uuid.UUID) {
	u.Navigation = NewNavigation(app, u.NavigationInfo(app, map[string]interface{}{"id": meterId}))
}

func NewUpdateMeter(app *TerminalApp, meterId uuid.UUID) *UpdateMeter {
	f := &UpdateMeter{
		FlexApp: NewFlexApp(),
		app:     app,
	}
	f.enrichNavigation(app, meterId)
	f.InitFlexApp(app)
	f.bindKeys()

	meterDto, err := app.GetMeterService().FindById(meterId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}
	var detailsArray []string

	if meterDto.Details != nil {
		for key, value := range meterDto.Details {
			detailsArray = append(detailsArray, fmt.Sprintf("%s : %.2f", key, value))
		}
	}

	var request updateMeterReq

	form := tview.NewForm().
		AddInputField("Name", meterDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", meterDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Details", strings.Join(detailsArray, "; "), 20, nil, func(text string) { request.details = text }).
		AddButton("Update", f.update(request, meterId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update Meter").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateMeter) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}

func (u *UpdateMeter) update(update updateMeterReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdateMeterRequest{
			Name:        update.name,
			Description: update.description,
		}

		details, err := parseDetails(update.details)
		if err != nil {
			u.ShowErrorTo(err)
			return
		}
		request.Details = details

		if err := u.app.GetMeterService().Update(id, request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack("Payment %s successfully updated.", request.Name)
		}
	}
}
