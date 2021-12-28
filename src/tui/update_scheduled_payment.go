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

const UpdateScheduledPaymentPageName = "payment-update-page"

type updateScheduledPaymentReq struct {
	name, description, sum, spec string
	providerId                   uuid.UUID
}

type UpdateScheduledPayment struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (u *UpdateScheduledPayment) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(UpdateScheduledPaymentPageName, func() tview.Primitive { return NewUpdateScheduledPayment(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdateScheduledPayment) enrichNavigation(app *TerminalApp, variables map[string]any) {
	u.Navigation = NewNavigation(app, u.NavigationInfo(app, variables))
}

func NewUpdateScheduledPayment(app *TerminalApp, scheduledPaymentId uuid.UUID) *UpdateScheduledPayment {
	f := &UpdateScheduledPayment{
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

	providers, providerOptions := GetProviders(app)

	var request updateScheduledPaymentReq

	form := tview.NewForm().
		AddInputField("Name", paymentDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", paymentDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Sum", fmt.Sprintf("%.2f", paymentDto.Sum), 20, nil, func(text string) { request.sum = text }).
		AddDropDown("Spec", specs, 1, func(option string, optionIndex int) {
			request.spec = option
		}).
		AddDropDown("Provider", providerOptions, -1, func(option string, optionIndex int) {
			request.providerId = providers[optionIndex].Id
		}).
		AddButton("Update", f.update(request, scheduledPaymentId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update Scheduled Payment").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateScheduledPayment) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}

func (u *UpdateScheduledPayment) update(update updateScheduledPaymentReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdatePaymentSchedulerRequest{
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

		if request.ProviderId == DefaultUUID {
			u.ShowErrorTo(errors.New("provider id is not valid"))

			return
		}
		request.ProviderId = update.providerId

		if err := u.app.GetPaymentSchedulerService().Update(id, request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack(fmt.Sprintf("Scheduled Payment %s successfully updated.", request.Name))
		}
	}
}
