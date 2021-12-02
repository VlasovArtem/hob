package tui

import (
	"fmt"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const HomePageName = "home"

type Home struct {
	*tview.Flex
	*Navigation
	*Keyboard
	currentHouse *CurrentHouse
}

func (h *Home) my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo {
	return NewNavigationInfo(HomePageName, func() tview.Primitive { return NewHome(app) })
}

func (h *Home) enrichNavigation(app *TerminalApp, parent *NavigationInfo) {
	h.MyNavigation = interface{}(h).(MyNavigation)
	h.enrich(app, parent).
		addCustomPageToMe(app, &CreateHouse{}).
		addCustomPageToMe(app, &CreateIncome{}).
		addCustomPageToMe(app, &CreatePayment{}).
		addCustomPage(app, nil, &Payments{})
}

func NewHome(app *TerminalApp) *Home {

	h := &Home{
		Flex:         tview.NewFlex(),
		Navigation:   NewNavigation(),
		Keyboard:     NewKeyboard(),
		currentHouse: NewCurrentHouse(app.House),
	}
	h.enrichNavigation(app, nil)

	h.bindKeys()

	houseList := tview.NewList().ShowSecondaryText(false)
	houseList.SetTitle("Houses").SetBorder(true)
	paymentsInfo := tview.NewTextView()
	paymentsInfo.SetTitle("Payments").SetBorder(true)
	incomesInfo := tview.NewTextView()
	incomesInfo.SetTitle("Incomes").SetBorder(true)

	info := tview.NewFlex().
		AddItem(houseList, 0, 1, true).
		AddItem(paymentsInfo, 0, 2, false).
		AddItem(incomesInfo, 0, 2, false)

	flex := h.
		SetDirection(tview.FlexRow).
		AddItem(NewMenuBlock(h.Keyboard), 0, 1, false)

	h.currentHouse.enrich(flex)

	flex.AddItem(info, 0, 6, true)

	h.showHouses(houseList, func(dto houseModel.HouseDto) {
		paymentsInfo.Clear()
		incomesInfo.Clear()

		h.showPayments(paymentsInfo, dto.Id)
		h.showIncomes(incomesInfo, dto.Id)
		h.currentHouse.Clear()
		h.currentHouse.house = dto
		h.currentHouse.setText()
	})

	h.SetInputCapture(h.KeyboardFunc)

	return h
}

func (h *Home) showHouses(houseList *tview.List, onSelect func(dto houseModel.HouseDto)) {
	housesData := h.app.GetHouseService().FindByUserId(h.app.AuthorizedUser.Id)

	hasCurrentHouse := h.app.House != houseModel.HouseDto{}

	var currentIndex int

	for index, house := range housesData {
		if hasCurrentHouse && h.app.House.Id == house.Id {
			currentIndex = index
		}
		houseList.AddItem(house.Name, house.Id.String(), 0, nil)
	}

	houseList.AddItem("Add New House", "", 0, func() {
		h.NavigateTo(CreateHousePageName)
	})

	houseList.SetSelectedFunc(h.setHouse(onSelect))
	houseList.SetChangedFunc(h.setHouse(onSelect))

	houseList.SetCurrentItem(1)
	houseList.SetCurrentItem(currentIndex)
}

func (h *Home) setHouse(onSelect func(dto houseModel.HouseDto)) func(index int, mainText string, secondaryText string, shortcut rune) {
	return func(index int, mainText string, secondaryText string, shortcut rune) {
		if secondaryText != "" {
			id, err := uuid.Parse(secondaryText)
			if err != nil {
				log.Fatal().Err(err)
			}
			if houseDto, err := h.app.GetHouseService().FindById(id); err != nil {
				h.ShowErrorTo(err)
			} else {
				h.app.House = houseDto
				h.currentHouse.house = houseDto
				onSelect(houseDto)
			}
		}
	}
}

func (h *Home) showPayments(paymentsInfo *tview.TextView, houseId uuid.UUID) {
	payments := h.app.GetPaymentService().FindByHouseId(houseId)
	if len(payments) == 0 {
		fmt.Fprintf(paymentsInfo, "No payments")
	} else {
		for _, paymentDto := range payments {
			fmt.Fprintf(paymentsInfo, "%s - %v\n", paymentDto.Name, paymentDto.Sum)
		}
	}
}

func (h *Home) showIncomes(incomesInfo *tview.TextView, houseId uuid.UUID) {
	incomes := h.app.GetIncomeService().FindByHouseId(houseId)
	if len(incomes) == 0 {
		fmt.Fprintf(incomesInfo, "No incomes")
	} else {
		for _, incomeDto := range incomes {
			fmt.Fprintf(incomesInfo, "%s - %v\n", incomeDto.Name, incomeDto.Sum)
		}
	}
}

func (h *Home) createIncome(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(CreateIncomePageName)
	return key
}

func (h *Home) bindKeys() {
	h.Keyboard.Actions = KeyActions{
		tcell.KeyCtrlP: NewKeyAction("Create Payment", h.createPayment),
		tcell.KeyCtrlJ: NewKeyAction("Create Income", h.createIncome),
		tcell.KeyCtrlA: NewKeyAction("Show Payments", h.paymentsPage),
	}
}

func (h *Home) homePage(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateHome()
	return key
}

func (h *Home) createPayment(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(CreatePaymentPageName)
	return key
}

func (h *Home) paymentsPage(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(PaymentsPageName)
	return key
}
