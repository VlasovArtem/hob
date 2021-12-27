package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type MenuBlock struct {
	*tview.Flex
	app  *TerminalApp
	info *SessionInfo
}

func NewMenuBlock(app *TerminalApp, keyboard *Keyboard) *MenuBlock {
	menuBlock := &MenuBlock{tview.NewFlex(), app, NewSessionInfo(app)}
	menuBlock.
		AddItem(NewLogo(), 0, 1, false).
		AddItem(menuBlock.info, 0, 2, false).
		AddItem(NewMenu(keyboard.Actions), 0, 5, false)

	return menuBlock
}

func (m *MenuBlock) refreshSessionInfo() {
	m.info.refresh()
}

type Menu struct {
	*tview.Table
}

func NewMenu(actions KeyActions) *Menu {
	menu := Menu{tview.NewTable()}

	for i, hint := range actions.Hints() {
		menu.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("<%s>", strings.ToLower(hint.Mnemonic))).SetAlign(tview.AlignLeft).SetTextColor(tcell.ColorLightBlue))
		menu.SetCell(i, 1, tview.NewTableCell(hint.Description))
	}

	menu.SetTitle("Menu")

	return &menu
}
