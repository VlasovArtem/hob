package tui

import (
	"fmt"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"reflect"
	"strconv"
	"time"
)

const PaymentsPageName = "payments"

var paymentsTableHeader = []string{"#", "Name", "Description", "Date", "Sum"}

type Payments struct {
	*tview.Flex
	*Navigation
	*Keyboard
	search       *Search
	currentHouse *CurrentHouse
}

func (p *Payments) my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo {
	return NewNavigationInfo(PaymentsPageName, func() tview.Primitive { return NewPayments(app) })
}

func (p *Payments) enrichNavigation(app *TerminalApp, parent *NavigationInfo) {
	p.MyNavigation = interface{}(p).(MyNavigation)
	p.enrich(app, parent).
		addCustomPage(app, nil, &CreatePayment{})
}

func NewPayments(app *TerminalApp) *Payments {
	p := &Payments{
		Flex:         tview.NewFlex(),
		Navigation:   NewNavigation(),
		Keyboard:     NewKeyboard(),
		search:       NewSearch(app),
		currentHouse: NewCurrentHouse(app.House),
	}
	p.enrichNavigation(app, nil)

	p.bindKeys()

	flex := p.
		SetDirection(tview.FlexRow).
		AddItem(NewMenuBlock(p.Keyboard), 0, 2, false)
	p.currentHouse.enrich(flex)
	flex.
		AddItem(p.search, 3, 0, false).
		AddItem(p.createPaymentTable(), 0, 10, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (p *Payments) createPaymentTable() *tview.Table {
	table := tview.NewTable()
	table.SetBorderPadding(1, 1, 1, 1)
	table.SetBorders(true)
	table.SetTitle("Payments")

	for i, header := range paymentsTableHeader {
		table.SetCell(0, i, tview.NewTableCell(header).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorYellow))
	}

	for i, dto := range p.app.GetPaymentService().FindByHouseId(p.currentHouse.house.Id) {
		fillRow(i+1, dto, table)
	}

	return table
}

func fillRow(currentRow int, dto paymentModel.PaymentDto, table *tview.Table) {
	for i, name := range paymentsTableHeader {
		var value string
		if i == 0 {
			value = strconv.Itoa(currentRow)
		} else {
			byName := reflect.ValueOf(dto).FieldByName(name)

			if byName.Kind() == reflect.Invalid {
				log.Fatal().Msgf("field with name %s not found in object %v", name, dto)
			}

			switch byName.Type().String() {
			case "time.Time":
				value = byName.Interface().(time.Time).Format("2006-01-02")
			default:
				value = fmt.Sprintf("%v", byName)
			}
		}
		table.SetCell(currentRow, i, tview.NewTableCell(value))
	}
}

func (p *Payments) bindKeys() {
	p.Keyboard.Actions = KeyActions{
		tcell.KeyCtrlP:  NewKeyAction("Create Payment", p.createPayment),
		tcell.KeyEscape: NewKeyAction("Back Home", p.homePage),
		tcell.KeyEnter:  NewKeyAction("Search", p.startSearch),
	}
}

func (p *Payments) startSearch(key *tcell.EventKey) *tcell.EventKey {
	p.search.Focus(nil)
	return key
}

func (p *Payments) createPayment(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateTo(CreatePaymentPageName)
	return key
}

func (p *Payments) homePage(key *tcell.EventKey) *tcell.EventKey {
	p.NavigateHome()
	return key
}
