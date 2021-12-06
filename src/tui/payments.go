package tui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const PaymentsPageName = "payments"

var paymentsTableHeader = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Id").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Name"),
	NewTableHeader("Description"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion())}

type Payments struct {
	*FlexApp
	*Navigation
	payments *TableFiller
}

func (p *Payments) my(app *TerminalApp, ctx context.Context) *NavigationInfo {
	return NewNavigationInfo(PaymentsPageName, func() tview.Primitive { return NewPayments(app) })
}

func (p *Payments) enrichNavigation(app *TerminalApp, ctx context.Context) {
	p.MyNavigation = interface{}(p).(MyNavigation)
	p.enrich(app, ctx).
		addCustomPage(ctx, &CreatePayment{})
}

func NewPayments(app *TerminalApp) *Payments {
	p := &Payments{
		FlexApp:    NewFlexApp(),
		Navigation: NewNavigation(),
		payments:   NewTableFiller(paymentsTableHeader),
	}
	p.enrichNavigation(app, nil)

	p.bindKeys()
	p.InitFlexApp(app)

	p.
		AddItem(p.fillTable(), 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (p *Payments) fillTable() *TableFiller {
	p.payments.SetSelectable(true, false)
	p.payments.SetTitle("Payments")
	content := p.app.GetPaymentService().FindByHouseId(p.app.House.Id)
	p.payments.fill(content)
	return p.payments
}

func (p *Payments) bindKeys() {
	p.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create Payment", p.createPayment),
		tcell.KeyCtrlD:  NewKeyAction("Delete Payment", p.deletePayment),
		tcell.KeyCtrlU:  NewKeyAction("Update Payment", p.updatePayment),
		tcell.KeyEscape: NewKeyAction("Back Home", p.homePage),
	}
}

func (p *Payments) createPayment(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateTo(CreatePaymentPageName)
	return key
}

func (p *Payments) homePage(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateHome()
	return key
}

func (p *Payments) deletePayment(key *tcell.EventKey) *tcell.EventKey {
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
		showModal(p.app.Main, fmt.Sprintf("Do you want to delete payment %s (%s)?", paymentId, paymentName), []modalButton{
			p.createDeleteModalButton(paymentName, paymentId),
		})
	}

	return key
}

func (p *Payments) createDeleteModalButton(paymentName string, paymentId uuid.UUID) modalButton {
	return modalButton{
		name: "Delete",
		action: func() {
			if err := p.app.GetPaymentService().DeleteById(paymentId); err != nil {
				p.ShowErrorTo(err)
			} else {
				p.ShowInfoRefresh(fmt.Sprintf("Payment %s (%s) successfully deleted.", paymentName, paymentId))
			}
		},
	}
}

func (p *Payments) updatePayment(key *tcell.EventKey) *tcell.EventKey {
	err := p.payments.performWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(UpdatePaymentPageName, func() tview.Primitive {
			updateContext := context.WithValue(context.Background(), UpdatePaymentPageName, id.String())
			return NewUpdatePayment(p.app, updateContext)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}
