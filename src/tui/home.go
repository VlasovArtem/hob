package tui

import (
	"fmt"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"os"
)

const HomePageName = "home"

var homePaymentFields = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Name"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion())}
var homeIncomeFields = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Name"),
	NewTableHeader("Date").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Sum").SetContentModifier(AlignCenterExpansion())}

type Home struct {
	*FlexApp
	*Navigation
	payments *TableFiller
	incomes  *TableFiller
}

func (h *Home) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(HomePageName, func() tview.Primitive { return NewHome(app) })
}

func (h *Home) enrichNavigation(app *TerminalApp) {
	h.Navigation = NewNavigation(app, h.NavigationInfo(app, nil))
	h.
		AddCustomPage(&Payments{}).
		AddCustomPage(&Houses{}).
		AddCustomPage(&Incomes{}).
		AddCustomPage(&Providers{})
}

func NewHome(app *TerminalApp) *Home {
	h := &Home{
		FlexApp:  NewFlexApp(),
		payments: NewTableFiller(homePaymentFields),
		incomes:  NewTableFiller(homeIncomeFields),
	}

	h.Init(app)

	return h
}

func (h *Home) Init(app *TerminalApp) {
	h.bindKeys()

	h.InitFlexApp(app)

	h.enrichNavigation(app)

	houseList := tview.NewList().ShowSecondaryText(false)
	houseList.
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Houses").
		SetBorder(true)

	h.payments.
		SetSelectable(false, false).
		SetTitle("Payments for Current Month").
		SetBorder(true)
	h.incomes.
		SetSelectable(false, false).
		SetTitle("Incomes for Current Month").
		SetBorder(true)

	info := tview.NewFlex().
		AddItem(houseList, 0, 1, true).
		AddItem(h.fillPaymentsTable(), 0, 2, false).
		AddItem(h.fillIncomesTable(), 0, 2, false)

	h.AddItem(info, 0, 8, true)

	h.showHouses(houseList, func(dto houseModel.HouseDto) {
		h.fillPaymentsTable()
		h.fillIncomesTable()
		h.menu.refreshSessionInfo()
	})

	h.SetInputCapture(h.KeyboardFunc)
}

func (h *Home) showHouses(houseList *tview.List, onSelect func(dto houseModel.HouseDto)) {
	housesData := h.App.GetHouseService().FindByUserId(h.App.AuthorizedUser.Id)

	hasCurrentHouse := h.App.House != nil

	for _, house := range housesData {
		if hasCurrentHouse && h.App.House.Id == house.Id {
			onSelect(house)
		}
		houseList.AddItem(house.Name, house.Id.String(), 0, nil)
	}

	houseList.AddItem("Add New House", "", 0, func() {
		h.NavigateTo(CreateHousePageName)
	})

	houseList.SetSelectedFunc(h.setHouse(onSelect))
}

func (h *Home) setHouse(onSelect func(dto houseModel.HouseDto)) func(index int, mainText string, secondaryText string, shortcut rune) {
	return func(index int, mainText string, secondaryText string, shortcut rune) {
		if secondaryText != "" {
			id, err := uuid.Parse(secondaryText)
			if err != nil {
				log.Fatal().Err(err)
			}
			if houseDto, err := h.App.GetHouseService().FindById(id); err != nil {
				h.ShowErrorTo(err)
			} else {
				h.App.House = &houseDto
				onSelect(houseDto)
			}
		}
	}
}

func (h *Home) bindKeys() {
	h.Keyboard.Actions = KeyActions{
		tcell.KeyCtrlA: NewKeyAction("Show Payments", h.paymentsPage),
		tcell.KeyCtrlE: NewKeyAction("Show Houses", h.housesPage),
		tcell.KeyCtrlF: NewKeyAction("Show Incomes", h.incomesPage),
		tcell.KeyCtrlP: NewKeyAction("Show Providers", h.providersPage),
		tcell.KeyCtrlQ: NewKeyAction("Quit", func(key *tcell.EventKey) *tcell.EventKey {
			os.Exit(0)
			return key
		}),
	}
}

func (h *Home) paymentsPage(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(PaymentsPageName)
	return key
}

func (h *Home) housesPage(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(HousesPageName)
	return key
}

func (h *Home) incomesPage(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(IncomesPageName)
	return key
}

func (h *Home) providersPage(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(ProvidersPageName)
	return key
}

func (h *Home) fillPaymentsTable() *TableFiller {
	if h.App.House == nil {
		return h.payments
	}
	payments := h.App.GetPaymentService().FindByHouseId(h.App.House.Id, 100, 0, nil, nil)

	h.payments.Fill(payments)
	var sum float64
	for _, payment := range payments {
		sum += float64(payment.Sum)
	}
	h.payments.addResultRow(fmt.Sprintf("%v", sum))
	return h.payments
}

func (h *Home) fillIncomesTable() *TableFiller {
	if h.App.House == nil {
		return h.incomes
	}
	incomes := h.App.GetIncomeService().FindByHouseId(h.App.House.Id, 100, 0, nil, nil)

	h.incomes.Fill(incomes)
	var sum float64
	for _, income := range incomes {
		sum += float64(income.Sum)
	}
	h.incomes.addResultRow(fmt.Sprintf("%v", sum))
	return h.incomes
}
