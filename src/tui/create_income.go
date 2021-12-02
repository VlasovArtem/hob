package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

const CreateIncomePageName = "create-income"

type CreateIncome struct {
	*tview.Form
	*NavigationBack
	request model.CreateIncomeRequest
}

func (c *CreateIncome) my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo {
	return NewNavigationInfo(CreateIncomePageName, func() tview.Primitive { return NewCreateIncome(app, parent) })
}

func (c *CreateIncome) enrichNavigation(app *TerminalApp, parent *NavigationInfo) {
	c.MyNavigation = interface{}(c).(MyNavigation)
	c.enrich(app, parent)
}

func NewCreateIncome(app *TerminalApp, parent *NavigationInfo) *CreateIncome {
	f := &CreateIncome{
		Form:           tview.NewForm(),
		NavigationBack: NewNavigationBack(NewNavigation()),
		request: model.CreateIncomeRequest{
			HouseId: app.House.Id,
			Date:    time.Now(),
		},
	}
	f.enrich(app, parent)

	f.
		AddInputField("Name", "", 20, nil, func(text string) { f.request.Name = text }).
		AddInputField("Description", "", 20, nil, func(text string) { f.request.Description = text }).
		AddInputField("Date (ex. 2006-01-02)", "", 20, nil, func(text string) {
			parse, err := time.Parse("2006-01-02", text)
			if err != nil {
				f.ShowErrorTo(errors.New("date is not valid"))
			} else {
				f.request.Date = parse
			}
		}).
		AddInputField("Sum", "", 20, nil, func(text string) {
			if sum, err := strconv.ParseFloat(text, 32); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.request.Sum = float32(sum)
			}
		}).
		AddButton("Create", func() {
			if _, err := f.app.GetIncomeService().Add(f.request); err != nil {
				f.ShowErrorTo(err)
			} else {
				f.ShowInfoReturnBack(fmt.Sprintf("Income for the house %s successfully added.", f.app.House.Name))
			}
		})

	f.SetBorder(true).SetTitle("Add Income").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}
