package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

type PageName string

type NavigationRules struct {
	owner      PageName
	home       *NavigationInfo
	onError    *NavigationInfo
	refresh    *NavigationInfo
	back       *NavigationInfo
	additional map[PageName]*NavigationInfo
}

func NewNavigationRules(app *TerminalApp, my *NavigationInfo) *NavigationRules {
	n := &NavigationRules{additional: make(map[PageName]*NavigationInfo)}
	n.home = &NavigationInfo{
		pageName: HomePageName,
		provider: func() tview.Primitive { return NewHome(app) },
	}
	n.refresh = my
	n.onError = my
	n.owner = my.pageName
	n.additional[my.pageName] = my
	return n
}

type MyNavigation interface {
	my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo
	enrichNavigation(app *TerminalApp, parent *NavigationInfo)
}

type Navigation struct {
	MyNavigation
	*NavigationRules
	app *TerminalApp
}

func NewNavigation() *Navigation {
	return &Navigation{}
}

func (n *Navigation) enrich(app *TerminalApp, parent *NavigationInfo) *Navigation {
	n.enrichMe(app, n.MyNavigation.my(app, nil))
	n.NavigationRules.back = parent
	return n
}

func (n *Navigation) enrichMe(app *TerminalApp, my *NavigationInfo) *Navigation {
	n.app = app
	n.NavigationRules = NewNavigationRules(app, my)
	return n
}

func (n *Navigation) addCustomPage(app *TerminalApp, parent *NavigationInfo, nav MyNavigation) *Navigation {
	navInfo := nav.my(app, parent)
	n.NavigationRules.additional[navInfo.pageName] = navInfo
	return n
}

func (n *Navigation) addCustomPageToMe(app *TerminalApp, nav MyNavigation) *Navigation {
	navInfo := nav.my(app, n.MyNavigation.my(app, nil))
	n.NavigationRules.additional[navInfo.pageName] = navInfo
	return n
}

func (n *Navigation) OnError() {
	n.Navigate(n.NavigationRules.onError)
}

func (n *Navigation) Refresh() {
	n.Navigate(n.NavigationRules.refresh)
}

func (n *Navigation) Back() {
	n.Navigate(n.NavigationRules.back)
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

func (n *Navigation) NavigateHome() {
	n.Navigate(n.NavigationRules.home)
}

func (n *Navigation) Navigate(info *NavigationInfo) {
	n.Show(info.pageName, info.provider)
}

func (n *Navigation) ShowOnMe(primitive tview.Primitive) {
	n.Show(n.NavigationRules.owner, func() tview.Primitive { return primitive })
}

func (n *Navigation) Show(pageName PageName, provider NavigationProvider) {
	n.app.Main.AddAndSwitchToPage(pageName.String(), provider(), true)
}

func (n *Navigation) NavigateToMe() {
	n.Navigate(n.NavigationRules.refresh)
}

func (n *Navigation) DoneFuncRefresh(key tcell.Key) {
	n.Navigate(n.NavigationRules.refresh)
}

func (n *Navigation) DoneFuncBack(key tcell.Key) {
	n.Navigate(n.NavigationRules.back)
}

func (n *Navigation) DoneFuncHome(key tcell.Key) {
	n.Navigate(n.NavigationRules.home)
}

func (n *Navigation) DoneFuncError(key tcell.Key) {
	n.Navigate(n.NavigationRules.onError)
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

type NavigationBack struct {
	*Navigation
	*Keyboard
}

func NewNavigationBack(navigation *Navigation) *NavigationBack {
	back := &NavigationBack{
		Keyboard:   NewKeyboard(),
		Navigation: navigation,
	}

	back.bindKeys()

	return back
}

func (c *NavigationBack) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.backToParent),
	}
}

func (c *NavigationBack) backToParent(key *tcell.EventKey) *tcell.EventKey {
	c.Back()
	return key
}
