package tui

import (
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

const CreateScheduledIncomePageName = "create-scheduled-payment"

type createScheduledIncomeReq struct {
	name, description, sum, spec string
}

type CreateScheduledIncome struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (c *CreateScheduledIncome) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(CreateScheduledIncomePageName, func() tview.Primitive { return NewCreateScheduledIncome(app) })
}

func (c *CreateScheduledIncome) enrichNavigation(app *TerminalApp) {
	c.Navigation = NewNavigation(app, c.NavigationInfo(app, nil))
}

func NewCreateScheduledIncome(app *TerminalApp) *CreateScheduledIncome {
	f := &CreateScheduledIncome{
		app:     app,
		FlexApp: NewFlexApp(),
	}
	f.bindKeys()
	f.InitFlexApp(app)
	f.enrichNavigation(app)

	var request createScheduledIncomeReq

	form := tview.NewForm().
		AddInputField("Name", "", DefaultInputFieldWidth, nil, func(text string) { request.name = text }).
		AddInputField("Description", "", DefaultInputFieldWidth, nil, func(text string) { request.description = text }).
		AddInputField("Sum", "", DefaultInputFieldWidth, nil, func(text string) { request.sum = text }).
		AddDropDown("Spec", specs, 1, func(option string, optionIndex int) {
			request.spec = option
		}).
		AddButton("Create", f.create(request)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Add Scheduled Income").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreateScheduledIncome) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.KeyBack),
	}
}

func (c *CreateScheduledIncome) create(request createScheduledIncomeReq) func() {
	return func() {
		sum, err := strconv.ParseFloat(request.sum, 32)

		if err != nil {
			c.ShowErrorTo(err)
			return
		}

		paymentRequest := model.CreateIncomeSchedulerRequest{
			HouseId:     c.app.House.Id,
			Name:        request.name,
			Description: request.description,
			Sum:         float32(sum),
			Spec:        scheduler.SchedulingSpecification(request.spec),
		}

		if _, err := c.app.GetIncomeSchedulerService().Add(paymentRequest); err != nil {
			c.ShowErrorTo(err)
		} else {
			c.ShowInfoReturnBack("Scheduled Income for the house %s successfully added.", c.app.House.Name)
		}
	}
}
