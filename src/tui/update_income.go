package tui

import (
	"context"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"strconv"
	"time"
)

const UpdateIncomePageName = "income-update-page"

type updateIncomeReq struct {
	name, description, date, sum string
}

type UpdateIncome struct {
	*FlexApp
	*Navigation
	app           *TerminalApp
	updateContent context.Context
}

func (u *UpdateIncome) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(UpdateIncomePageName, func() tview.Primitive { return NewUpdateIncome(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdateIncome) enrichNavigation(app *TerminalApp, incomeId uuid.UUID) {
	u.Navigation = NewNavigation(app, u.NavigationInfo(app, map[string]interface{}{"id": incomeId}))
}

func NewUpdateIncome(app *TerminalApp, incomeId uuid.UUID) *UpdateIncome {
	f := &UpdateIncome{
		FlexApp: NewFlexApp(),
		app:     app,
	}
	f.enrichNavigation(app, incomeId)
	f.InitFlexApp(app)
	f.bindKeys()

	incomeDto, err := app.GetIncomeService().FindById(incomeId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	var request updateIncomeReq

	form := tview.NewForm().
		AddInputField("Name", incomeDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", incomeDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Date (ex. 2006-01-02)", incomeDto.Date.Format("2006-01-02"), 20, nil, func(text string) { request.date = text }).
		AddInputField("Sum", fmt.Sprintf("%.2f", incomeDto.Sum), 20, nil, func(text string) { request.sum = text }).
		AddButton("Update", f.update(request, incomeId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update a House").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateIncome) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}

func (u *UpdateIncome) update(update updateIncomeReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdateIncomeRequest{
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

		if err := u.app.GetIncomeService().Update(id, request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack("Income %s successfully updated.", request.Name)
		}
	}
}
