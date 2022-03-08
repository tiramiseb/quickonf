package help

import "strings"

func (m *Model) tabsView() string {
	var view string
	switch m.activeSection {
	case 0:
		view = m.introTitleWithFocus + "│" + m.languageTitle + "│" + m.commandsTitle + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 1:
		view = m.introTitle + "│" + m.languageTitleWithFocus + "│" + m.commandsTitle + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 2:
		view = m.introTitle + "│" + m.languageTitle + "│" + m.commandsTitleWithFocus + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 3:
		view = m.introTitle + "│" + m.languageTitle + "│" + m.commandsTitle + "│" + m.uiTitleWithFocus + "\n" + m.subtitleSeparator
	}
	return view
}

func (m *Model) filterView() string {
	var prefix string
	if m.width > 20 {
		prefix = "Filter: "
	}
	return prefix + m.commandFilter + "\n" + strings.Repeat("─", m.width)
}

func (m *Model) View() string {
	if m.activeSection == 2 {
		return m.tabsView() + "\n" + m.filterView() + "\n" + m.viewport.View()
	} else {
		return m.tabsView() + "\n" + m.viewport.View()
	}
}
