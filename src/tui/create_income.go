package tui

import (
	"errors"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

const CreateIncomePageName = "create-income"

type createIncome struct {
	name, description, date, sum string
}

type CreateIncome struct {
	*FlexApp
	*Navigation
	menu    *MenuBlock
	request model.CreateIncomeRequest
}

func (c *CreateIncome) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(CreateIncomePageName, func() tview.Primitive { return NewCreateIncome(app) })
}

func (c *CreateIncome) enrichNavigation(app *TerminalApp) {
	c.Navigation = NewNavigation(app, c.NavigationInfo(app, nil))
}

func NewCreateIncome(app *TerminalApp) *CreateIncome {
	f := &CreateIncome{
		FlexApp: NewFlexApp(),
		request: model.CreateIncomeRequest{
			HouseId: &app.House.Id,
			Date:    time.Now(),
		},
	}
	f.bindKeys()
	f.InitFlexApp(app)
	f.enrichNavigation(app)

	var create createIncome

	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) { create.name = text }).
		AddInputField("Description", "", 20, nil, func(text string) { create.description = text }).
		AddInputField("Date (ex. 2006-01-02)", time.Now().Format("2006-01-02"), 20, nil, func(text string) { create.date = text }).
		AddInputField("Sum", "", 20, nil, func(text string) { create.sum = text }).
		AddButton("Create", f.create(&create))

	form.SetBorder(true).SetTitle("Add Income").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreateIncome) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.KeyBack),
	}
}

func (c *CreateIncome) create(create *createIncome) func() {
	return func() {
		if newSum, err := strconv.ParseFloat(create.sum, 32); err != nil {
			c.ShowErrorTo(errors.New("sum is not valid"))
		} else {
			c.request.Sum = float32(newSum)
		}

		if newDate, err := time.Parse("2006-01-02", create.date); err != nil {
			c.ShowErrorTo(errors.New("date is not valid format. The valid format is 2006-01-02"))
		} else {
			c.request.Date = newDate
		}

		if _, err := c.App.GetIncomeService().Add(c.request); err != nil {
			c.ShowErrorTo(err)
		} else {
			c.ShowInfoReturnBack("Income for the house %s successfully added.", c.App.House.Name)
		}
	}
}
