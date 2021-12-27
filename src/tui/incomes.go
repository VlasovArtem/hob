package tui

import (
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

func (i *Incomes) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(IncomesPageName, func() tview.Primitive { return NewIncomes(app) })
}

func (i *Incomes) enrichNavigation(app *TerminalApp) {
	i.Navigation = NewNavigation(app, i.NavigationInfo(app, nil))
	i.AddCustomPage(&CreateIncome{})
}

func NewIncomes(app *TerminalApp) *Incomes {
	p := &Incomes{
		FlexApp: NewFlexApp(),
		incomes: NewTableFiller(incomesTableHeader),
	}
	p.enrichNavigation(app)

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
	content := i.App.GetIncomeService().FindByHouseId(i.App.House.Id)
	i.incomes.Fill(content)
	return i.incomes
}

func (i *Incomes) bindKeys() {
	i.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create Income", i.createIncome),
		tcell.KeyCtrlD:  NewKeyAction("Delete Income", i.deleteIncome),
		tcell.KeyCtrlU:  NewKeyAction("Update Income", i.updateIncome),
		tcell.KeyEscape: NewKeyAction("Back Home", i.KeyHome),
	}
}

func (i *Incomes) createIncome(key *tcell.EventKey) *tcell.EventKey {
	i.NavigateTo(CreateIncomePageName)
	return key
}

func (i *Incomes) deleteIncome(key *tcell.EventKey) *tcell.EventKey {
	err := i.incomes.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		name := i.incomes.GetCell(row, 2).Text
		ShowModal(i.App.Main, fmt.Sprintf("Do you want to delete income %s (%s)?", id, name), []ModalButton{
			i.createDeleteModalButton(name, id),
		})
	})

	if err != nil {
		i.ShowErrorTo(err)
	}

	return key
}

func (i *Incomes) updateIncome(key *tcell.EventKey) *tcell.EventKey {
	err := i.incomes.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		i.Navigate(NewNavigationInfo(UpdateIncomePageName, func() tview.Primitive {
			return NewUpdateIncome(i.App, id)
		}))
	})

	if err != nil {
		i.ShowErrorTo(err)
	}

	return key
}

func (i *Incomes) createDeleteModalButton(name string, id uuid.UUID) ModalButton {
	return ModalButton{
		Name: "Delete",
		Action: func() {
			if err := i.App.GetIncomeService().DeleteById(id); err != nil {
				i.ShowErrorTo(err)
			} else {
				i.ShowInfoRefresh(fmt.Sprintf("Income %s (%s) successfully deleted.", name, id))
			}
		},
	}
}
