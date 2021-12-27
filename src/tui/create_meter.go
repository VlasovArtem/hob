package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

const CreateMeterPageName = "create-meter"

type createMeterReq struct {
	name, description, details string
}

type CreateMeter struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (c *CreateMeter) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(CreateMeterPageName, func() tview.Primitive { return NewCreateMeter(app, variables["paymentId"].(uuid.UUID)) })
}

func (c *CreateMeter) enrichNavigation(app *TerminalApp, paymentId uuid.UUID) {
	c.Navigation = NewNavigation(app, c.NavigationInfo(app, map[string]interface{}{"paymentId": paymentId}))
}

func NewCreateMeter(app *TerminalApp, paymentId uuid.UUID) *CreateMeter {
	f := &CreateMeter{
		app:     app,
		FlexApp: NewFlexApp(),
	}
	f.bindKeys()
	f.InitFlexApp(app)
	f.enrichNavigation(app, paymentId)

	var request createMeterReq

	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", "", 20, nil, func(text string) { request.description = text }).
		AddInputField("Details", "", 20, nil, func(text string) { request.details = text }).
		AddButton("Create", f.create(paymentId, request)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Add Meter").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreateMeter) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.KeyBack),
	}
}

func (c *CreateMeter) create(paymentId uuid.UUID, request createMeterReq) func() {
	return func() {
		details, err := parseDetails(request.details)
		if err != nil {
			c.ShowErrorTo(err)
		}

		if paymentId == DefaultUUID {
			c.ShowErrorTo(errors.New("payment id is not valid"))
			return
		}

		paymentRequest := model.CreateMeterRequest{
			PaymentId:   paymentId,
			Name:        request.name,
			Description: request.description,
			Details:     details,
		}

		if _, err := c.app.GetMeterService().Add(paymentRequest); err != nil {
			c.ShowErrorTo(err)
		} else {
			c.ShowInfoReturnBack(fmt.Sprintf("Payment for the house %s successfully added.", c.app.House.Name))
		}
	}
}

func parseDetails(requestDetails string) (details map[string]float64, err error) {
	if requestDetails != "" {
		detailsInfo := strings.Split(requestDetails, ";")

		details = make(map[string]float64)
		for _, detail := range detailsInfo {
			if detail != "" {
				detailContent := strings.Split(strings.Trim(detail, " "), ":")

				if len(detailContent) != 2 {
					return details, errors.New("details content is not valid. The meter details should have the next style 'first:12,95; second:11;'")
				}

				if float, err := strconv.ParseFloat(strings.Trim(detailContent[1], " "), 2); err != nil {
					return details, errors.New("details content is not valid. The meter details should have the next style 'first:12,95; second:11;'")
				} else {
					details[strings.Trim(detailContent[0], " ")] = float
				}
			}
		}
	}
	return details, err
}
