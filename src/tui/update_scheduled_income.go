package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/scheduler"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"strconv"
)

const UpdateScheduledIncomePageName = "scheduled-income-update-page"

type updateScheduledIncomeReq struct {
	name, description, sum, spec string
}

type UpdateScheduledIncome struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (u *UpdateScheduledIncome) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(UpdateScheduledIncomePageName, func() tview.Primitive { return NewUpdateScheduledIncome(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdateScheduledIncome) enrichNavigation(app *TerminalApp, variables map[string]any) {
	u.Navigation = NewNavigation(app, u.NavigationInfo(app, variables))
}

func NewUpdateScheduledIncome(app *TerminalApp, scheduledPaymentId uuid.UUID) *UpdateScheduledIncome {
	f := &UpdateScheduledIncome{
		FlexApp: NewFlexApp(),
		app:     app,
	}
	f.enrichNavigation(app, map[string]any{"id": scheduledPaymentId})
	f.InitFlexApp(app)
	f.bindKeys()

	paymentDto, err := app.GetPaymentSchedulerService().FindById(scheduledPaymentId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	var request updateScheduledIncomeReq

	form := tview.NewForm().
		AddInputField("Name", paymentDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", paymentDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Sum", fmt.Sprintf("%.2f", paymentDto.Sum), 20, nil, func(text string) { request.sum = text }).
		AddDropDown("Spec", specs, 1, func(option string, optionIndex int) {
			request.spec = option
		}).
		AddButton("Update", f.update(request, scheduledPaymentId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update Scheduled Payment").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateScheduledIncome) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}

func (u *UpdateScheduledIncome) update(update updateScheduledIncomeReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdateIncomeSchedulerRequest{
			Name:        update.name,
			Description: update.description,
			Spec:        scheduler.SchedulingSpecification(update.spec),
		}

		if newSum, err := strconv.ParseFloat(update.sum, 32); err != nil {
			u.ShowErrorTo(errors.New("sum is not valid"))
			return
		} else {
			request.Sum = float32(newSum)
		}

		if err := u.app.GetIncomeSchedulerService().Update(id, request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack("Scheduled Income %s successfully updated.", request.Name)
		}
	}
}
