package tui

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/payment/model"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
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
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Provider").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Meter Id").SetContentModifier(AlignCenterExpansion()),
}

type Payments struct {
	*FlexApp
	*Navigation
	payments *TableFiller
}

func (p *Payments) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(PaymentsPageName, func() tview.Primitive { return NewPayments(app) })
}

func NewPayments(app *TerminalApp) *Payments {
	p := &Payments{
		FlexApp:  NewFlexApp(),
		payments: NewTableFiller(paymentsTableHeader),
	}
	p.enrichNavigation(app)

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
	p.payments.TableHeaders[6].SetContentProvider(p.findProviderName)
	p.payments.TableHeaders[7].SetContentProvider(p.findMeterId)
	content := p.App.GetPaymentService().FindByHouseId(p.App.House.Id)
	p.payments.Fill(content)
	return p.payments
}

func (p *Payments) findProviderName(payment interface{}) interface{} {
	providerId := payment.(model.PaymentDto).ProviderId

	providerDto, err := p.App.GetProviderService().FindById(providerId)

	if err != nil {
		p.ShowErrorTo(err)

		return nil
	} else {
		return providerDto.Name
	}
}

func (p *Payments) findMeterId(payment interface{}) interface{} {
	paymentId := payment.(model.PaymentDto).Id

	meterDto, err := p.App.GetMeterService().FindByPaymentId(paymentId)

	if err != nil {
		return nil
	}
	return meterDto.Id
}

func (p *Payments) enrichNavigation(app *TerminalApp) {
	p.Navigation = NewNavigation(app, p.NavigationInfo(app, nil))
	p.AddCustomPage(&CreatePayment{})
}

func (p *Payments) bindKeys() {
	p.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create Payment", p.createPayment),
		tcell.KeyCtrlD:  NewKeyAction("Delete Payment", p.deletePayment),
		tcell.KeyCtrlU:  NewKeyAction("Update Payment", p.updatePayment),
		tcell.KeyCtrlJ:  NewKeyAction("Add Meter", p.createMeter),
		tcell.KeyCtrlF:  NewKeyAction("Update Meter", p.updateMeter),
		tcell.KeyCtrlM:  NewKeyAction("Show Meter", p.showMeter),
		tcell.KeyCtrlS:  NewKeyAction("Show Scheduled", p.showScheduled),
		tcell.KeyEscape: NewKeyAction("Back Home", p.KeyHome),
	}
}

func (p *Payments) createPayment(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateTo(CreatePaymentPageName)
	return key
}

func (p *Payments) createMeter(key *tcell.EventKey) *tcell.EventKey {
	err := p.payments.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(CreateMeterPageName, func() tview.Primitive {
			return NewCreateMeter(p.App, id)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}

func (p *Payments) updateMeter(key *tcell.EventKey) *tcell.EventKey {
	err := p.payments.PerformWithSelectedId(7, func(row int, meterId uuid.UUID) {
		p.Navigate(NewNavigationInfo(UpdateMeterPageName, func() tview.Primitive {
			return NewUpdateMeter(p.App, meterId)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}

func (p *Payments) showMeter(key *tcell.EventKey) *tcell.EventKey {
	row, _ := p.payments.GetSelection()
	if row == 1 {
		return key
	}
	meterIdString := p.payments.GetCell(row, 7).Text
	meterId, err := uuid.Parse(meterIdString)
	if err != nil {
		p.ShowErrorTo(err)
	} else {
		meterDto, err := p.App.GetMeterService().FindById(meterId)

		if err != nil {
			p.ShowErrorTo(err)
		} else {

			ShowModal(
				p.App.Main,
				fmt.Sprintf("Name: %s\nDescription: %s\nDetails: %v", meterDto.Name, meterDto.Description, meterDto.Details),
				[]ModalButton{},
			)
		}
	}

	return key
}

func (p *Payments) showScheduled(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateTo(CreatePaymentPageName)
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
		ShowModal(p.App.Main, fmt.Sprintf("Do you want to delete payment %s (%s)?", paymentId, paymentName), []ModalButton{
			p.createDeleteModalButton(paymentName, paymentId),
		})
	}

	return key
}

func (p *Payments) createDeleteModalButton(paymentName string, paymentId uuid.UUID) ModalButton {
	return ModalButton{
		Name: "Delete",
		Action: func() {
			if err := p.App.GetPaymentService().DeleteById(paymentId); err != nil {
				p.ShowErrorTo(err)
			} else {
				p.ShowInfoRefresh(fmt.Sprintf("Payment %s (%s) successfully deleted.", paymentName, paymentId))
			}
		},
	}
}

func (p *Payments) updatePayment(key *tcell.EventKey) *tcell.EventKey {
	err := p.payments.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(UpdatePaymentPageName, func() tview.Primitive {
			return NewUpdatePayment(p.App, id)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}

func GetProviders(t *TerminalApp) ([]providerModel.ProviderDto, []string) {
	providerDtos := t.GetProviderService().FindByUserId(t.AuthorizedUser.Id)
	var providerOptions []string
	for _, provider := range providerDtos {
		providerOptions = append(providerOptions, provider.Name)
	}

	return providerDtos, providerOptions
}
