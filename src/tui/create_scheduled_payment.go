package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"strconv"
)

var specs = []string{string(scheduler.HOURLY), string(scheduler.DAILY), string(scheduler.WEEKLY), string(scheduler.MONTHLY), string(scheduler.ANNUALLY)}

const CreateScheduledPaymentPageName = "create-scheduled-payment"

type createScheduledPaymentReq struct {
	name, description, sum, spec string
	providerId                   uuid.UUID
}

type CreateScheduledPayment struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (c *CreateScheduledPayment) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(CreateScheduledPaymentPageName, func() tview.Primitive { return NewCreateScheduledPayment(app) })
}

func (c *CreateScheduledPayment) enrichNavigation(app *TerminalApp) {
	c.Navigation = NewNavigation(app, c.NavigationInfo(app, nil))
}

func NewCreateScheduledPayment(app *TerminalApp) *CreateScheduledPayment {
	f := &CreateScheduledPayment{
		app:     app,
		FlexApp: NewFlexApp(),
	}
	f.bindKeys()
	f.InitFlexApp(app)
	f.enrichNavigation(app)

	providers, providerOptions := GetProviders(app)

	var request createScheduledPaymentReq

	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", "", 20, nil, func(text string) { request.description = text }).
		AddInputField("Sum", "", 20, nil, func(text string) { request.sum = text }).
		AddDropDown("Spec", specs, 1, func(option string, optionIndex int) {
			request.spec = option
		}).
		AddDropDown("Provider", providerOptions, -1, func(option string, optionIndex int) {
			request.providerId = providers[optionIndex].Id
		}).
		AddButton("Create", f.create(request)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Add Payment").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreateScheduledPayment) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.KeyBack),
	}
}

func (c *CreateScheduledPayment) create(request createScheduledPaymentReq) func() {
	return func() {
		sum, err := strconv.ParseFloat(request.sum, 32)

		if err != nil {
			c.ShowErrorTo(err)
			return
		}

		if request.providerId == DefaultUUID {
			c.ShowErrorTo(errors.New("provider id is not valid"))
			return
		}

		paymentRequest := model.CreatePaymentSchedulerRequest{
			UserId:      c.app.AuthorizedUser.Id,
			HouseId:     c.app.House.Id,
			ProviderId:  request.providerId,
			Name:        request.name,
			Description: request.description,
			Sum:         float32(sum),
			Spec:        scheduler.SchedulingSpecification(request.spec),
		}

		if _, err := c.app.GetPaymentSchedulerService().Add(paymentRequest); err != nil {
			c.ShowErrorTo(err)
		} else {
			c.ShowInfoReturnBack(fmt.Sprintf("Scheduled Payment for the house %s successfully added.", c.app.House.Name))
		}
	}
}
