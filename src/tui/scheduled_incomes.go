package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const ScheduledIncomesPageName = "scheduled-incomes"

var scheduledIncomesTableHeader = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Id").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Name"),
	NewTableHeader("Description"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Spec").SetContentModifier(AlignCenterExpansion()),
}

type ScheduledIncomes struct {
	*FlexApp
	*Navigation
	scheduledIncomes *TableFiller
}

func (p *ScheduledIncomes) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(ScheduledIncomesPageName, func() tview.Primitive { return NewScheduledIncomes(app) })
}

func NewScheduledIncomes(app *TerminalApp) *ScheduledIncomes {
	p := &ScheduledIncomes{
		FlexApp:          NewFlexApp(),
		scheduledIncomes: NewTableFiller(scheduledIncomesTableHeader),
	}
	p.enrichNavigation(app)

	p.bindKeys()
	p.InitFlexApp(app)

	p.
		AddItem(p.fillTable(), 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (p *ScheduledIncomes) fillTable() *TableFiller {
	p.scheduledIncomes.SetSelectable(true, false)
	p.scheduledIncomes.SetTitle("Scheduled Incomes")
	content := p.App.GetPaymentService().FindByHouseId(p.App.House.Id, 100, 0, nil, nil)
	p.scheduledIncomes.Fill(content)
	return p.scheduledIncomes
}

func (p *ScheduledIncomes) enrichNavigation(app *TerminalApp) {
	p.Navigation = NewNavigation(app, p.NavigationInfo(app, nil))
	p.AddCustomPage(&CreateScheduledIncome{})
}

func (p *ScheduledIncomes) bindKeys() {
	p.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create", p.createScheduledIncome),
		tcell.KeyCtrlD:  NewKeyAction("Delete", p.deleteScheduledIncome),
		tcell.KeyCtrlU:  NewKeyAction("Update", p.updateScheduledIncome),
		tcell.KeyEscape: NewKeyAction("Back", p.KeyBack),
	}
}

func (p *ScheduledIncomes) createScheduledIncome(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateTo(CreateScheduledIncomePageName)
	return key
}

func (p *ScheduledIncomes) deleteScheduledIncome(key *tcell.EventKey) *tcell.EventKey {
	row, _ := p.scheduledIncomes.GetSelection()
	if row == 1 {
		return key
	}
	id, err := uuid.Parse(p.scheduledIncomes.GetCell(row, 1).Text)
	if err != nil {
		p.ShowErrorTo(err)
	} else {
		name := p.scheduledIncomes.GetCell(row, 2).Text
		ShowModal(p.App.Main, fmt.Sprintf("Do you want to delete scheduled income %s (%s)?", id, name), []ModalButton{
			p.createDeleteModalButton(name, id),
		})
	}

	return key
}

func (p *ScheduledIncomes) createDeleteModalButton(name string, id uuid.UUID) ModalButton {
	return ModalButton{
		Name: "Delete",
		Action: func() {
			if err := p.App.GetIncomeSchedulerService().DeleteById(id); err != nil {
				p.ShowErrorTo(err)
			} else {
				p.ShowInfoRefresh("Income scheduler %s (%s) successfully deleted.", name, id)
			}
		},
	}
}

func (p *ScheduledIncomes) updateScheduledIncome(key *tcell.EventKey) *tcell.EventKey {
	err := p.scheduledIncomes.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(UpdateScheduledIncomePageName, func() tview.Primitive {
			return NewUpdateScheduledIncome(p.App, id)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}
