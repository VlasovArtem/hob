package tui

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/pivotal/calculator"
	"github.com/VlasovArtem/hob/src/pivotal/model"
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
	calculator     calculator.PivotalCalculatorService
}

func (p *Pivotal) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(PivotalPageName, func() tview.Primitive { return NewPivotal(app) })
}

func NewPivotal(app *TerminalApp) *Pivotal {
	p := &Pivotal{
		FlexApp:        NewFlexApp(),
		details:        tview.NewTextView(),
		pivotalService: app.GetPivotalService(),
		calculator:     app.GetPivotalCalculatorService(),
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
	if p.App.House == nil {
		p.ShowErrorTo(errors.New("house is not set"))
	} else {
		var pivotal model.PivotalResponseDto
		var err error

		if !p.pivotalService.ExistsByHouseId(p.App.House.Id) {
			pivotal, err = p.calculator.Calculate(p.App.House.Id)
		} else {
			pivotal, err = p.pivotalService.Find(p.App.House.Id)
		}

		if err != nil {
			p.ShowErrorTo(err)
		} else {
			details := fmt.Sprintf("House Pivotal: \nIncomes: %.2f\nExpenses: %.2f\nTotal: %.2f\n\n", pivotal.House.Income, pivotal.House.Payments, pivotal.House.Total)
			for _, group := range pivotal.Groups {
				details += fmt.Sprintf("Group '%s' Pivotal: \nIncomes: %.2f\nExpenses: %.2f\nTotal: %.2f\n", group.Group.Name, group.Income, group.Payments, group.Total)
			}
			details += fmt.Sprintf("\nTotal Pivotal: \nIncomes: %.2f\nExpenses: %.2f\nTotal: %.2f\n", pivotal.Total.Income, pivotal.Total.Payments, pivotal.Total.Total)

			_, err = p.details.Write([]byte(details))

			if err != nil {
				p.ShowErrorTo(err)
			}
		}
	}
}
