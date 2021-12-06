package tui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const IncomesPageName = "incomes"

var incomesTableHeader = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Id").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Name"),
	NewTableHeader("Description"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion())}

type Incomes struct {
	*FlexApp
	*Navigation
	incomes *TableFiller
}

func (i *Incomes) my(app *TerminalApp, ctx context.Context) *NavigationInfo {
	return NewNavigationInfo(IncomesPageName, func() tview.Primitive { return NewIncomes(app) })
}

func (i *Incomes) enrichNavigation(app *TerminalApp, ctx context.Context) {
	i.MyNavigation = interface{}(i).(MyNavigation)
	i.enrich(app, ctx).
		addCustomPage(ctx, &CreateIncome{})
}

func NewIncomes(app *TerminalApp) *Incomes {
	p := &Incomes{
		FlexApp:    NewFlexApp(),
		Navigation: NewNavigation(),
		incomes:    NewTableFiller(incomesTableHeader),
	}
	p.enrichNavigation(app, nil)

	p.bindKeys()
	p.InitFlexApp(app)

	p.
		AddItem(p.fillTable(), 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (i *Incomes) fillTable() *TableFiller {
	i.incomes.SetSelectable(true, false)
	i.incomes.SetTitle("Incomes")
	content := i.app.GetIncomeService().FindByHouseId(i.app.House.Id)
	i.incomes.fill(content)
	return i.incomes
}

func (i *Incomes) bindKeys() {
	i.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create Income", i.createIncome),
		tcell.KeyCtrlD:  NewKeyAction("Delete Income", i.deleteIncome),
		tcell.KeyCtrlU:  NewKeyAction("Update Income", i.updateIncome),
		tcell.KeyEscape: NewKeyAction("Back Home", i.homePage),
	}
}

func (i *Incomes) createIncome(key *tcell.EventKey) *tcell.EventKey {
	i.NavigateTo(CreateIncomePageName)
	return key
}

func (i *Incomes) homePage(key *tcell.EventKey) *tcell.EventKey {
	i.NavigateHome()
	return key
}

func (i *Incomes) deleteIncome(key *tcell.EventKey) *tcell.EventKey {
	err := i.incomes.performWithSelectedId(1, func(row int, id uuid.UUID) {
		name := i.incomes.GetCell(row, 2).Text
		showModal(i.app.Main, fmt.Sprintf("Do you want to delete income %s (%s)?", id, name), []modalButton{
			i.createDeleteModalButton(name, id),
		})
	})

	if err != nil {
		i.ShowErrorTo(err)
	}

	return key
}

func (i *Incomes) updateIncome(key *tcell.EventKey) *tcell.EventKey {
	err := i.incomes.performWithSelectedId(1, func(row int, id uuid.UUID) {
		i.Navigate(NewNavigationInfo(UpdateIncomePageName, func() tview.Primitive {
			updateContext := context.WithValue(context.Background(), UpdateIncomePageName, id.String())
			return NewUpdateIncome(i.app, updateContext)
		}))
	})

	if err != nil {
		i.ShowErrorTo(err)
	}

	return key
}

func (i *Incomes) createDeleteModalButton(name string, id uuid.UUID) modalButton {
	return modalButton{
		name: "Delete",
		action: func() {
			if err := i.app.GetIncomeService().DeleteById(id); err != nil {
				i.ShowErrorTo(err)
			} else {
				i.ShowInfoRefresh(fmt.Sprintf("Income %s (%s) successfully deleted.", name, id))
			}
		},
	}
}
