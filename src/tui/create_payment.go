package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

// Object to interface interface{}(f).(NavigationPage)

const CreatePaymentPageName = "create-payment"

type CreatePayment struct {
	*tview.Form
	*NavigationBack
	app     *TerminalApp
	request model.CreatePaymentRequest
}

func (c *CreatePayment) my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo {
	return NewNavigationInfo(CreatePaymentPageName, func() tview.Primitive { return NewCreatePayment(app, parent) })
}

func (c *CreatePayment) enrichNavigation(app *TerminalApp, parent *NavigationInfo) {
	c.MyNavigation = interface{}(c).(MyNavigation)
	c.enrich(app, parent)
}

func NewCreatePayment(app *TerminalApp, parent *NavigationInfo) *CreatePayment {
	f := &CreatePayment{
		Form:           tview.NewForm(),
		app:            app,
		NavigationBack: NewNavigationBack(NewNavigation()),
		request: model.CreatePaymentRequest{
			UserId:  app.AuthorizedUser.Id,
			HouseId: app.House.Id,
			Date:    time.Now(),
		},
	}
	f.enrichNavigation(app, parent)

	var name, description, date, sum string

	f.
		AddInputField("Name", "", 20, nil, func(text string) { name = text }).
		AddInputField("Description", "", 20, nil, func(text string) { description = text }).
		AddInputField("Date (ex. 2006-01-02)", "", 20, nil, func(text string) { date = text }).
		AddInputField("Sum", "", 20, nil, func(text string) { sum = text }).
		AddButton("Create", func() {
			f.request.Name = name
			f.request.Description = description

			parse, err := time.Parse("2006-01-02", date)
			if err != nil {
				f.ShowErrorTo(errors.New("date is not valid"))
			} else {
				f.request.Date = parse
			}

			if sum, err := strconv.ParseFloat(sum, 32); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.request.Sum = float32(sum)
			}

			if _, err := f.app.GetPaymentService().Add(f.request); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.ShowInfoReturnBack(fmt.Sprintf("Payment for the house %s successfully added.", f.app.House.Name))
			}
		})

	f.SetBorder(true).SetTitle("Add Payment").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}
