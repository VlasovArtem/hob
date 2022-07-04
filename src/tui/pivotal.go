package tui

import (
	"errors"
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

	if p.App.House == nil {
		p.ShowErrorTo(errors.New("house is not set"))
	} else {

		var incomesSum, paymentsSum, total float64
		var err error

		if len(p.App.House.Groups) != 0 {
			for _, group := range p.App.House.Groups {
				pivotalDto, nestedError := p.pivotalService.FindByGroupId(group.Id)
				if nestedError != nil {
					err = nestedError
					break
				}
				incomesSum += pivotalDto.Income
				paymentsSum += pivotalDto.Payments
				total += pivotalDto.Total
			}
		} else {
			if pivotalByHouse, nestedError := p.pivotalService.FindByHouseId(p.App.House.Id); nestedError != nil {
				err = nestedError
			} else {
				incomesSum = pivotalByHouse.Income
				paymentsSum = pivotalByHouse.Payments
				total = pivotalByHouse.Total
			}
		}

		if err != nil {
			p.ShowErrorTo(err)
		} else {

			details += fmt.Sprintf("Income: %.2f\n", incomesSum)
			details += fmt.Sprintf("Payments: %.2f\n", paymentsSum)
			details += "\n"
			details += fmt.Sprintf("Total: %.2f\n", total)

			_, err = p.details.Write([]byte(details))

			if err != nil {
				p.ShowErrorTo(err)
			}
		}
	}
}
