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

type Navigator interface {
	NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo
}

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
	App *TerminalApp
}

func NewNavigation(app *TerminalApp, my *NavigationInfo) *Navigation {
	return &Navigation{
		App:             app,
		NavigationRules: NewNavigationRules(app, my),
	}
}

func (n *Navigation) AddCustomPageWithVars(navigator Navigator, variables map[string]interface{}) *Navigation {
	navigationInfo := navigator.NavigationInfo(n.App, variables)
	n.additional[navigationInfo.pageName] = navigationInfo
	return n
}

func (n *Navigation) AddCustomPage(navigator Navigator) *Navigation {
	return n.AddCustomPageWithVars(navigator, nil)
}

func (n *Navigation) OnError() {
	n.NavigateTo(OnError)
}

func (n *Navigation) Refresh() {
	n.NavigateTo(Refresh)
}

func (n *Navigation) Back() {
	n.App.Main.RemovePage(n.my.pageName.String())
}

func (n *Navigation) Home() {
	n.NavigateTo(homePage)
}

func (n *Navigation) NavigateTo(name PageName) {
	if info, ok := n.NavigationRules.additional[name]; !ok {
		log.Error().Msgf("Page with Name %s not exits", name)
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
	n.App.Main.AddAndSwitchToPage(pageName.String(), provider(), true)
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

func (n *Navigation) KeyRefresh(key *tcell.EventKey) *tcell.EventKey {
	n.Refresh()
	return key
}

func (n *Navigation) KeyBack(key *tcell.EventKey) *tcell.EventKey {
	n.Back()
	return key
}

func (n *Navigation) KeyHome(key *tcell.EventKey) *tcell.EventKey {
	n.Home()
	return key
}

func (n *Navigation) KeyError(key *tcell.EventKey) *tcell.EventKey {
	n.OnError()
	return key
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
