package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog/log"
	"sort"
)

type (
	ActionHandler func(*tcell.EventKey) *tcell.EventKey

	KeyAction struct {
		Description string
		Action      ActionHandler
	}

	KeyActions map[tcell.Key]KeyAction
)

func NewKeyAction(description string, action ActionHandler) KeyAction {
	return KeyAction{Description: description, Action: action}
}

func (a KeyActions) Add(aa KeyActions) {
	for k, v := range aa {
		a[k] = v
	}
}

func (a KeyActions) Clear() {
	for k := range a {
		delete(a, k)
	}
}

func (a KeyActions) Set(aa KeyActions) {
	for k, v := range aa {
		a[k] = v
	}
}

func (a KeyActions) Delete(kk ...tcell.Key) {
	for _, k := range kk {
		delete(a, k)
	}
}

func (a KeyActions) Hints() MenuHints {
	kk := make([]int, 0, len(a))
	for k := range a {
		kk = append(kk, int(k))
	}
	sort.Ints(kk)

	menuHints := make(MenuHints, 0, len(kk))
	for _, key := range kk {
		if name, ok := tcell.KeyNames[tcell.Key(key)]; ok {
			menuHints = append(menuHints,
				MenuHint{
					Mnemonic:    name,
					Description: a[tcell.Key(key)].Description,
				},
			)
		} else {
			log.Error().Msgf("Unable to locate KeyName for %#v", key)
		}
	}
	return menuHints
}
