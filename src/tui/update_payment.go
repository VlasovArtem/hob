package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

const UpdatePaymentPageName = "payment-update-page"

type updatePaymentReq struct {
	name, description, date, sum string
	providerId                   uuid.UUID
}

type UpdatePayment struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (u *UpdatePayment) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(UpdatePaymentPageName, func() tview.Primitive { return NewUpdatePayment(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdatePayment) enrichNavigation(app *TerminalApp, variables map[string]interface{}) {
	u.Navigation = NewNavigation(app, u.NavigationInfo(app, variables))
}

func NewUpdatePayment(app *TerminalApp, paymentId uuid.UUID) *UpdatePayment {
	f := &UpdatePayment{
		FlexApp: NewFlexApp(),
		app:     app,
	}
	f.enrichNavigation(app, map[string]interface{}{"id": paymentId})
	f.InitFlexApp(app)
	f.bindKeys()

	paymentDto, err := app.GetPaymentService().FindById(paymentId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	providers, providerOptions := GetProviders(app)

	var request updatePaymentReq

	form := tview.NewForm().
		AddInputField("Name", paymentDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", paymentDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Date (ex. 2006-01-02)", paymentDto.Date.Format("2006-01-02"), 20, nil, func(text string) { request.date = text }).
		AddInputField("Sum", fmt.Sprintf("%.2f", paymentDto.Sum), 20, nil, func(text string) { request.sum = text }).
		AddDropDown("Provider", providerOptions, -1, func(option string, optionIndex int) {
			request.providerId = providers[optionIndex].Id
		}).
		AddButton("Update", f.update(request, paymentId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update Payment").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdatePayment) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}

func (u *UpdatePayment) update(update updatePaymentReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdatePaymentRequest{
			Name:        update.name,
			Description: update.description,
		}

		if newSum, err := strconv.ParseFloat(update.sum, 32); err != nil {
			u.ShowErrorTo(errors.New("sum is not valid"))
			return
		} else {
			request.Sum = float32(newSum)
		}

		if newDate, err := time.Parse("2006-01-02", update.date); err != nil {
			u.ShowErrorTo(errors.New("date is not valid format. The valid format is 2006-01-25"))
			return
		} else {
			request.Date = newDate
		}

		if request.ProviderId == DefaultUUID {
			u.ShowErrorTo(errors.New("provider id is not valid"))

			return
		}
		request.ProviderId = update.providerId

		if err := u.app.GetPaymentService().Update(id, request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack("Payment %s successfully updated.", request.Name)
		}
	}
}
