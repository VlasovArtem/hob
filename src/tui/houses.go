package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const HousesPageName = "houses"

var housesTableHeader = []*TableHeader{
	NewIndexHeader(),
	NewTableHeader("Id").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("Name"),
	NewTableHeader("CountryCode").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("City").SetContentModifier(AlignCenterExpansion()),
	NewTableHeader("StreetLine1"),
	NewTableHeader("StreetLine2")}

type Houses struct {
	*FlexApp
	*Navigation
	houses *TableFiller
}

func (h *Houses) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(HousesPageName, func() tview.Primitive { return NewHouses(app) })
}

func (h *Houses) enrichNavigation(app *TerminalApp) {
	h.Navigation = NewNavigation(
		app,
		h.NavigationInfo(app, nil),
	)
	h.AddCustomPage(&CreateHouse{})
}

func NewHouses(app *TerminalApp) *Houses {
	p := &Houses{
		FlexApp: NewFlexApp(),
		houses:  NewTableFiller(housesTableHeader),
	}
	p.enrichNavigation(app)

	p.bindKeys()
	p.InitFlexApp(app)

	p.
		AddItem(p.fillTable(), 0, 8, true).
		SetInputCapture(p.KeyboardFunc)

	return p
}

func (h *Houses) fillTable() *TableFiller {
	h.houses.SetSelectable(true, false)
	h.houses.SetTitle("Houses")
	content := h.App.GetHouseService().FindByUserId(h.App.AuthorizedUser.Id)
	h.houses.Fill(content)
	return h.houses
}

func (h *Houses) bindKeys() {
	h.Actions = KeyActions{
		tcell.KeyCtrlE:  NewKeyAction("Create House", h.createHouse),
		tcell.KeyCtrlU:  NewKeyAction("Update House", h.updateHouse),
		tcell.KeyCtrlD:  NewKeyAction("Delete House", h.deleteHouse),
		tcell.KeyEscape: NewKeyAction("Back Home", h.KeyHome),
	}
}

func (h *Houses) createHouse(key *tcell.EventKey) *tcell.EventKey {
	h.NavigateTo(CreateHousePageName)
	return key
}

func (h *Houses) updateHouse(key *tcell.EventKey) *tcell.EventKey {
	err := h.houses.PerformWithSelectedId(1, func(row int, id uuid.UUID) {
		h.Navigate(NewNavigationInfo(UpdateHousePageName, func() tview.Primitive {
			return NewUpdateHouse(h.App, id)
		}))
	})

	if err != nil {
		h.ShowErrorTo(err)
	}

	return key
}

func (h *Houses) deleteHouse(key *tcell.EventKey) *tcell.EventKey {
	err := h.houses.PerformWithSelectedId(1, func(row int, houseId uuid.UUID) {
		name := h.houses.GetCell(row, 2).Text
		ShowModal(h.App.Main, fmt.Sprintf("Do you want to delete house %s (%s)?", houseId, name), []ModalButton{
			h.createDeleteModalButton(name, houseId),
		})
	})

	if err != nil {
		h.ShowErrorTo(err)
	}

	return key
}

func (h *Houses) createDeleteModalButton(houseName string, houseId uuid.UUID) ModalButton {
	return ModalButton{
		Name: "Delete",
		Action: func() {
			if err := h.App.GetHouseService().DeleteById(houseId); err != nil {
				h.ShowErrorTo(err)
			} else {
				h.ShowInfoRefresh(fmt.Sprintf("House %s (%s) successfully deleted.", houseName, houseId))
			}
		},
	}
}
