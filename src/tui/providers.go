package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const ProvidersPageName = "providers"

var providersTableHeader = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Id").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Name"),
	NewTableHeader("Description"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion())}

type Providers struct {
	*FlexApp
	*Navigation
	providers *TableFiller
}

func (p *Providers) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(ProvidersPageName, func() tview.Primitive { return NewProviders(app) })
}

func NewProviders(app *TerminalApp) *Providers {
	p := &Providers{
		FlexApp:   NewFlexApp(),
		providers: NewTableFiller(providersTableHeader),
	}
	p.enrichNavigation(app)

	p.bindKeys()
	p.InitFlexApp(app)

	p.
		AddItem(p.fillTable(), 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (p *Providers) fillTable() *TableFiller {
	p.providers.SetSelectable(true, false)
	p.providers.SetTitle("Providers")
	content := p.App.GetProviderService().FindByUserId(p.App.AuthorizedUser.Id)
	p.providers.Fill(content)
	return p.providers
}

func (p *Providers) enrichNavigation(app *TerminalApp) {
	p.Navigation = NewNavigation(app, p.NavigationInfo(app, nil))
}

func (p *Providers) bindKeys() {
	p.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create Provider", p.createProviderWithAuthorizedUser),
		tcell.KeyCtrlO:  NewKeyAction("Create Provider as default", p.createProviderWithDefaultUser),
		tcell.KeyCtrlD:  NewKeyAction("Delete Provider", p.deleteProvider),
		tcell.KeyCtrlU:  NewKeyAction("Update Payment", p.updatePayment),
		tcell.KeyEscape: NewKeyAction("Back Home", p.KeyHome),
	}
}

func (p *Providers) createProviderWithAuthorizedUser(key *tcell.EventKey) *tcell.EventKey {
	p.createProvider(p.App.AuthorizedUser.Id)
	return key
}

func (p *Providers) createProviderWithDefaultUser(key *tcell.EventKey) *tcell.EventKey {
	p.createProvider(uuid.UUID{})
	return key
}

func (p *Providers) createProvider(userId uuid.UUID) {
	err := p.providers.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(CreateProviderPageName, func() tview.Primitive {
			return NewCreateProvider(p.App, userId)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
}

func (p *Providers) deleteProvider(key *tcell.EventKey) *tcell.EventKey {
	row, _ := p.providers.GetSelection()
	if row == 1 {
		return key
	}
	providerIdString := p.providers.GetCell(row, 1).Text
	providerId, err := uuid.Parse(providerIdString)
	if err != nil {
		p.ShowErrorTo(err)
	} else {
		providerName := p.providers.GetCell(row, 2).Text
		ShowModal(p.App.Main, fmt.Sprintf("Do you want to delete provider %s (%s)?", providerId, providerName), []ModalButton{
			p.createDeleteModalButton(providerName, providerId),
		})
	}

	return key
}

func (p *Providers) createDeleteModalButton(paymentName string, paymentId uuid.UUID) ModalButton {
	return ModalButton{
		Name: "Delete",
		Action: func() {
			if err := p.App.GetProviderService().Delete(paymentId); err != nil {
				p.ShowErrorTo(err)
			} else {
				p.ShowInfoRefresh("Provider %s (%s) successfully deleted.", paymentName, paymentId)
			}
		},
	}
}

func (p *Providers) updatePayment(key *tcell.EventKey) *tcell.EventKey {
	err := p.providers.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		p.Navigate(NewNavigationInfo(UpdateProviderPageName, func() tview.Primitive {
			return NewUpdatePayment(p.App, id)
		}))
	})

	if err != nil {
		p.ShowErrorTo(err)
	}
	return key
}
