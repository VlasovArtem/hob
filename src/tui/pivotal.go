package tui

import (
	"fmt"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const PivotalPageName = "pivotal"

type Pivotal struct {
	*FlexApp
	*Navigation
	details        *tview.TextView
	pivotalService pivotalService.PivotalService
}

func (p *Pivotal) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(PivotalPageName, func() tview.Primitive { return NewPivotal(app) })
}

func NewPivotal(app *TerminalApp) *Pivotal {
	p := &Pivotal{
		FlexApp:        NewFlexApp(),
		details:        tview.NewTextView(),
		pivotalService: app.GetPivotalService(),
	}
	p.enrichNavigation(app)

	p.bindKeys()
	p.InitFlexApp(app)

	p.initDetails()

	p.
		AddItem(p.details, 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (p *Pivotal) enrichNavigation(app *TerminalApp) {
	p.Navigation = NewNavigation(app, p.NavigationInfo(app, nil))
}

func (p *Pivotal) bindKeys() {
	p.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back Home", p.KeyHome),
	}
}

func (p *Pivotal) initDetails() {
	var details string

	pivotalByHouse, err := p.pivotalService.FindByHouseId(p.App.House.Id)

	if err != nil {
		p.ShowErrorTo(err)
	}

	details += fmt.Sprintf("Income: %.2f\n", pivotalByHouse.Income)
	details += fmt.Sprintf("Payments: %.2f\n", pivotalByHouse.Payments)
	details += "\n"
	details += fmt.Sprintf("Total: %.2f\n", pivotalByHouse.Total)

	_, err = p.details.Write([]byte(details))

	if err != nil {
		p.ShowErrorTo(err)
	}
}
