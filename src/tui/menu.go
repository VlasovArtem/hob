package tui

import (
	"fmt"
	"github.com/rivo/tview"
)

type MenuBlock struct {
	*tview.Flex
}

func NewMenuBlock(keyboard *Keyboard) *MenuBlock {
	menuBlock := &MenuBlock{tview.NewFlex()}
	menuBlock.
		AddItem(NewLogo(), 0, 1, false).
		AddItem(NewMenu(keyboard.Actions), 0, 5, false)

	return menuBlock
}

type Menu struct {
	*tview.Table
}

func NewMenu(actions KeyActions) *Menu {
	menu := Menu{tview.NewTable()}

	for i, hint := range actions.Hints() {
		cell := tview.NewTableCell(fmt.Sprintf("%-10s - %s", hint.Mnemonic, hint.Description))
		menu.SetCell(i, 0, cell)
	}

	menu.SetTitle("Menu")

	return &menu
}
