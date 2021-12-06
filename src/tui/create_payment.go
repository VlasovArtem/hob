package tui

import (
	"context"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

const CreatePaymentPageName = "create-payment"

type createPaymentReq struct {
	name, description, date, sum string
}

type CreatePayment struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (c *CreatePayment) my(app *TerminalApp, ctx context.Context) *NavigationInfo {
	return NewNavigationInfo(CreatePaymentPageName, func() tview.Primitive { return NewCreatePayment(app) })
}

func (c *CreatePayment) enrichNavigation(app *TerminalApp, ctx context.Context) {
	c.MyNavigation = interface{}(c).(MyNavigation)
	c.enrich(app, ctx)
}

func NewCreatePayment(app *TerminalApp) *CreatePayment {
	f := &CreatePayment{
		app:        app,
		FlexApp:    NewFlexApp(),
		Navigation: NewNavigation(),
	}
	f.bindKeys()
	f.InitFlexApp(app)
	f.enrichNavigation(app, nil)

	var request createPaymentReq

	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", "", 20, nil, func(text string) { request.description = text }).
		AddInputField("Date (ex. 2006-01-02)", "", 20, nil, func(text string) { request.date = text }).
		AddInputField("Sum", "", 20, nil, func(text string) { request.sum = text }).
		AddButton("Create", f.create(request)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Add Payment").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreatePayment) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.backToParent),
	}
}

func (c *CreatePayment) backToParent(key *tcell.EventKey) *tcell.EventKey {
	c.Back()
	return key
}

func (c *CreatePayment) create(request createPaymentReq) func() {
	return func() {
		newDate := time.Now()
		if request.date != "" {
			parsedDate, err := time.Parse("2006-01-02", request.date)
			if err != nil {
				c.ShowErrorTo(errors.New("date is not valid"))
				return
			}
			newDate = parsedDate
		}

		sum, err := strconv.ParseFloat(request.sum, 32)

		if err != nil {
			c.ShowErrorTo(err)
			return
		}

		paymentRequest := model.CreatePaymentRequest{
			UserId:      c.app.AuthorizedUser.Id,
			HouseId:     c.app.House.Id,
			Date:        newDate,
			Name:        request.name,
			Description: request.description,
			Sum:         float32(sum),
		}

		if _, err := c.app.GetPaymentService().Add(paymentRequest); err != nil {
			c.ShowErrorTo(err)
		} else {
			c.ShowInfoReturnBack(fmt.Sprintf("Payment for the house %s successfully added.", c.app.House.Name))
		}
	}
}
