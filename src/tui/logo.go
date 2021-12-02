package tui

import (
	"fmt"
	"github.com/rivo/tview"
)

// LogoSmall http://patorjk.com/software/taag/#p=display&f=Big&t=HoB
var LogoSmall = []string{
	` _    _       ____  `,
	`| |  | |     |  _ \ `,
	`| |__| | ___ | |_) |`,
	`|  __  |/ _ \|  _ < `,
	`| |  | | (_) | |_) |`,
	`|_|  |_|\___/|____/ `,
}

type Logo struct {
	*tview.TextView
}

func NewLogo() *Logo {
	logo := Logo{
		TextView: tview.NewTextView(),
	}

	logo.SetWordWrap(false)
	logo.SetWrap(false)
	logo.SetTextAlign(tview.AlignLeft)
	logo.SetDynamicColors(true)

	logo.Clear()
	for i, s := range LogoSmall {
		fmt.Fprintf(logo, "[%s::b]%s", "green", s)
		if i+1 < len(LogoSmall) {
			fmt.Fprintf(logo, "\n")
		}
	}

	return &logo
}
