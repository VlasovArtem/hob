package tui

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const ScheduledPaymentsPageName = "scheduled-payments"

var scheduledPaymentsTableHeader = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Id").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Name"),
	NewTableHeader("Description"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Spec").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Provider").SetContentModifier(AlignCenterExpansion()),
}

type ScheduledPayments struct {
	*FlexApp
	*Navigation
	payments *TableFiller
}

func (p *ScheduledPayments) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(ScheduledPaymentsPageName, func() tview.Primitive { return NewScheduledPayments(app) })
}

func NewScheduledPayments(app *TerminalApp) *ScheduledPayments {
	p := &ScheduledPayments{
		FlexApp:  NewFlexApp(),
		payments: NewTableFiller(scheduledPaymentsTableHeader),
	}
	p.enrichNavigation(app)

	p.bindKeys()
	p.InitFlexApp(app)

	p.
		AddItem(p.fillTable(), 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (p *ScheduledPayments) fillTable() *TableFiller {
	p.payments.SetSelectable(true, false)
	p.payments.SetTitle("Scheduled Payments")
	p.payments.AddContentProvider("Provider", p.findProviderName)
	content := p.App.GetPaymentService().FindByHouseId(p.App.House.Id, 100, 0, nil, nil)
	p.payments.Fill(content)
	return p.payments
}

func (p *ScheduledPayments) findProviderName(payment any) any {
	providerId := payment.(model.PaymentSchedulerDto).ProviderId

	providerDto, err := p.App.GetProviderService().FindById(providerId)

	if err != nil {
		p.ShowErrorTo(err)

		return nil
	} else {
		return providerDto.Name
	}
}

func (p *ScheduledPayments) enrichNavigation(app *TerminalApp) {
	p.Navigation = NewNavigation(app, p.NavigationInfo(app, nil))
	p.AddCustomPage(&CreateScheduledIncome{})
}

func (p *ScheduledPayments) bindKeys() {
	p.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create", p.createScheduledPayment),
		tcell.KeyCtrlD:  NewKeyAction("Delete", p.deleteScheduledPayment),
		tcell.KeyCtrlU:  NewKeyAction("Update", p.updateScheduledPayment),
		tcell.KeyEscape: NewKeyAction("Back", p.KeyBack),
	}
}

func (p *ScheduledPayments) createScheduledPayment(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateTo(CreateScheduledPaymentPageName)
	return key
}

func (p *ScheduledPayments) deleteScheduledPayment(key *tcell.EventKey) *tcell.EventKey {
	row, _ := p.payments.GetSelection()
	if row == 1 {
		return key
	}
	paymentIdString := p.payments.GetCell(row, 1).Text
	paymentId, err := uuid.Parse(paymentIdString)
	if err != nil {
		p.ShowErrorTo(err)
	} else {
		paymentName := p.payments.GetCell(row, 2).Text
		ShowModal(p.App.Main, fmt.Sprintf("Do you want to delete payment %s (%s)?", paymentId, paymentName), []ModalButton{
			p.createDeleteModalButton(paymentName, paymentId),
		})
	}

	return key
}

func (p *ScheduledPayments) createDeleteModalButton(paymentName string, paymentId uuid.UUID) ModalButton {
	return ModalButton{
		Name: "Delete",
		Action: func() {
			if err := p.App.GetPaymentSchedulerService().Remove(paymentId); err != nil {
				p.ShowErrorTo(err)
			} else {
				p.ShowInfoRefresh("Payment scheduler %s (%s) successfully deleted.", paymentName, paymentId)
			}
		},
	}
}

func (p *ScheduledPayments) updateScheduledPayment(key *tcell.EventKey) *tcell.EventKey {
	err := p.payments.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(UpdateScheduledPaymentPageName, func() tview.Primitive {
			return NewUpdateScheduledIncome(p.App, id)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}
