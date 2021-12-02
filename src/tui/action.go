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
	menuHints := make(MenuHints, 0)
	var keys []tcell.Key
	for key, _ := range a {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return i < j
	})

	for _, key := range keys {
		if name, ok := tcell.KeyNames[key]; ok {
			menuHints = append(menuHints,
				MenuHint{
					Mnemonic:    name,
					Description: a[key].Description,
				},
			)
		} else {
			log.Error().Msgf("Unable to locate KeyName for %#v", key)
		}
	}
	return menuHints
}
