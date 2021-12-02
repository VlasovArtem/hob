package tui

import (
	"github.com/gdamore/tcell/v2"
)

type Keyboard struct {
	Actions KeyActions
}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) KeyboardFunc(event *tcell.EventKey) *tcell.EventKey {
	if action, ok := k.Actions[event.Key()]; ok {
		action.Action(event)
		return nil
	}
	return event
}
