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

func (u *UpdateIncome) my(app *TerminalApp, ctx context.Context) *NavigationInfo {
	return NewNavigationInfo(UpdateIncomePageName, func() tview.Primitive { return NewUpdateIncome(app, ctx) })
}

func (u *UpdateIncome) enrichNavigation(app *TerminalApp, ctx context.Context) {
	u.MyNavigation = interface{}(u).(MyNavigation)
	u.enrich(app, ctx)
}

func NewUpdateIncome(app *TerminalApp, ctx context.Context) *UpdateIncome {
	f := &UpdateIncome{
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

	incomeDto, err := app.GetIncomeService().FindById(paymentId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	var request updateIncomeReq

	form := tview.NewForm().
		AddInputField("Name", incomeDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Description", incomeDto.Description, 20, nil, func(text string) { request.description = text }).
		AddInputField("Date (ex. 2006-01-02)", incomeDto.Date.Format("2006-01-02"), 20, nil, func(text string) { request.date = text }).
		AddInputField("Sum", fmt.Sprintf("%.2f", incomeDto.Sum), 20, nil, func(text string) { request.sum = text }).
		AddButton("Update", f.update(request, paymentId)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update a House").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateIncome) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", func(key *tcell.EventKey) *tcell.EventKey {
			u.Back()
			return key
		}),
	}
}

func (u *UpdateIncome) getId() (uuid.UUID, error) {
	idString := u.updateContent.Value(UpdateIncomePageName)

	houseId, err := uuid.Parse(idString.(string))

	if err != nil {
		return uuid.UUID{}, err
	}

	return houseId, nil
}

func (u *UpdateIncome) update(update updateIncomeReq, id uuid.UUID) func() {
	return func() {
		request := model.UpdateIncomeRequest{
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

		if err := u.app.GetIncomeService().Update(request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack(fmt.Sprintf("Income %s successfully updated.", request.Name))
		}
	}
}
