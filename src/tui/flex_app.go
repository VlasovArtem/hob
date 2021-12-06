package tui

import (
	"github.com/rivo/tview"
)

type FlexApp struct {
	*tview.Flex
	*Keyboard
	menu *MenuBlock
}

func NewFlexApp() *FlexApp {
	return &FlexApp{
		Flex:     tview.NewFlex(),
		Keyboard: NewKeyboard(),
	}
}

func (f *FlexApp) InitFlexApp(app *TerminalApp) {
	f.menu = NewMenuBlock(app, f.Keyboard)

	f.SetDirection(tview.FlexRow).
		AddItem(f.menu, 0, 1, false)
}
