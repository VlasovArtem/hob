package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

type PageName string

var (
	OnError  = PageName("on_error")
	Refresh  = PageName("refresh")
	homePage = PageName("home")
)

type NavigationRules struct {
	my         *NavigationInfo
	additional map[PageName]*NavigationInfo
}

func NewNavigationRules(app *TerminalApp, my *NavigationInfo) *NavigationRules {
	return &NavigationRules{
		my: my,
		additional: map[PageName]*NavigationInfo{
			homePage: {
				pageName: HomePageName,
				provider: func() tview.Primitive { return NewHome(app) },
			},
			Refresh:     my,
			OnError:     my,
			my.pageName: my,
		},
	}
}

type Navigation struct {
	*NavigationRules
	app *TerminalApp
}

func NewNavigation(app *TerminalApp, my *NavigationInfo) *Navigation {
	return &Navigation{
		NavigationRules: NewNavigationRules(app, my),
		app:             app,
	}
}

func (n *Navigation) addCustomPage(navigation *NavigationInfo) *Navigation {
	n.additional[navigation.pageName] = navigation
	return n
}

func (n *Navigation) OnError() {
	n.NavigateTo(OnError)
}

func (n *Navigation) Refresh() {
	n.NavigateTo(Refresh)
}

func (n *Navigation) Back() {
	n.app.Main.RemovePage(n.my.pageName.String())
}

func (n *Navigation) Home() {
	n.NavigateTo(homePage)
}

func (n *Navigation) NavigateTo(name PageName) {
	if info, ok := n.NavigationRules.additional[name]; !ok {
		log.Error().Msgf("Page with name %s not exits", name)
	} else {
		n.Navigate(info)
	}
}

func (n *Navigation) ShowInfoRefresh(msg string) {
	n.ShowInfo(msg, n.DoneFuncRefresh)
}

func (n *Navigation) ShowInfoReturnBack(msg string) {
	n.ShowInfo(msg, n.DoneFuncBack)
}

func (n *Navigation) ShowInfoReturnHome(msg string) {
	n.ShowInfo(msg, n.DoneFuncHome)
}

func (n *Navigation) ShowInfoReturnTo(msg string, name PageName) {
	n.ShowInfo(msg, func(key tcell.Key) {
		n.NavigateTo(name)
	})
}

func (n *Navigation) ShowInfo(msg string, doneFunc func(key tcell.Key)) {
	n.ShowOnMe(NewInfo(msg, doneFunc))
}

func (n *Navigation) ShowErrorTo(err error) {
	n.ShowError(err, n.DoneFuncError)
}

func (n *Navigation) ShowError(err error, doneFunc func(key tcell.Key)) {
	n.ShowOnMe(NewInfoWithError(err, doneFunc))
}

func (n *Navigation) Navigate(info *NavigationInfo) {
	n.Show(info.pageName, info.provider)
}

func (n *Navigation) ShowOnMe(primitive tview.Primitive) {
	n.Show(n.NavigationRules.my.pageName, func() tview.Primitive { return primitive })
}

func (n *Navigation) Show(pageName PageName, provider NavigationProvider) {
	n.app.Main.AddAndSwitchToPage(pageName.String(), provider(), true)
}

func (n *Navigation) BackFunc() func() {
	return func() {
		n.Back()
	}
}

func (n *Navigation) DoneFuncRefresh(key tcell.Key) {
	n.Refresh()
}

func (n *Navigation) DoneFuncBack(key tcell.Key) {
	n.Back()
}

func (n *Navigation) DoneFuncHome(key tcell.Key) {
	n.Home()
}

func (n *Navigation) DoneFuncError(key tcell.Key) {
	n.OnError()
}

func (p PageName) String() string {
	return string(p)
}

type NavigationProvider func() tview.Primitive

type NavigationInfo struct {
	pageName PageName
	provider NavigationProvider
}

func NewNavigationInfo(pageName PageName, provider func() tview.Primitive) *NavigationInfo {
	return &NavigationInfo{pageName: pageName, provider: provider}
}
