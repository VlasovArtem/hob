package tui

import (
	"context"
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
}

type UpdatePayment struct {
	*FlexApp
	*Navigation
	app           *TerminalApp
	updateContent context.Context
}

func (u *UpdatePayment) my(app *TerminalApp, ctx context.Context) *NavigationInfo {
	return NewNavigationInfo(UpdatePaymentPageName, func() tview.Primitive { return NewUpdatePayment(app, ctx) })
}

func (u *UpdatePayment) enrichNavigation(app *TerminalApp, ctx context.Context) {
	u.MyNavigation = interface{}(u).(MyNavigation)
	u.enrich(app, ctx)
}

func NewUpdatePayment(app *TerminalApp, ctx context.Context) *UpdatePayment {
	f := &UpdatePayment{
		FlexApp:       NewFlexApp(),
		Navigation:    NewNavigation(),
		app:           app,
		updateContent: ctx,
	}
	f.enrichNavigation(app, ctx)
	f.InitFlexApp(app)
	f.bindKeys()

	paymentId, err := f.getId()

	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	paymentDto, err := app.GetPaymentService().FindById(paymentId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	var request updatePaymentReq

	form := tview.NewForm().
		AddInputField("Name", paymentDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", paymentDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Date (ex. 2006-01-02)", paymentDto.Date.Format("2006-01-02"), 20, nil, func(text string) { request.date = text }).
		AddInputField("Sum", fmt.Sprintf("%.2f", paymentDto.Sum), 20, nil, func(text string) { request.sum = text }).
		AddButton("Update", f.update(request, paymentId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update a House").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdatePayment) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", func(key *tcell.EventKey) *tcell.EventKey {
			u.Back()
			return key
		}),
	}
}

func (u *UpdatePayment) getId() (uuid.UUID, error) {
	idString := u.updateContent.Value(UpdatePaymentPageName)

	houseId, err := uuid.Parse(idString.(string))

	if err != nil {
		return uuid.UUID{}, err
	}

	return houseId, nil
}

func (u *UpdatePayment) update(update updatePaymentReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdatePaymentRequest{
			Id:          id,
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

		if err := u.app.GetPaymentService().Update(request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack(fmt.Sprintf("Payment %s successfully updated.", request.Name))
		}
	}
}
