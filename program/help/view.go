package help

import "strings"

func (m *Model) tabsView() string {
	var view string
	switch m.activeSection {
	case 0:
		view = m.introTitleWithFocus + "│" + m.languageTitle + "│" + m.commandsTitle + "│" + m.recipesTitle + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 1:
		view = m.introTitle + "│" + m.languageTitleWithFocus + "│" + m.commandsTitle + "│" + m.recipesTitle + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 2:
		view = m.introTitle + "│" + m.languageTitle + "│" + m.commandsTitleWithFocus + "│" + m.recipesTitle + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 3:
		view = m.introTitle + "│" + m.languageTitle + "│" + m.commandsTitle + "│" + m.recipesTitleWithFocus + "│" + m.uiTitle + "\n" + m.subtitleSeparator
	case 4:
		view = m.introTitle + "│" + m.languageTitle + "│" + m.commandsTitle + "│" + m.recipesTitle + "│" + m.uiTitleWithFocus + "\n" + m.subtitleSeparator
	}
	return view
}

func (m *Model) commandFilterView() string {
	var prefix string
	if m.width > 20 {
		prefix = "Filter: "
	}
	return prefix + m.commandFilter + "\n" + strings.Repeat("─", m.width)
}

func (m *Model) recipeFilterView() string {
	var prefix string
	if m.width > 20 {
		prefix = "Filter: "
	}
	return prefix + m.recipeFilter + "\n" + strings.Repeat("─", m.width)
}

func (m *Model) View() string {
	switch m.activeSection {
	case 2:
		return m.tabsView() + "\n" + m.commandFilterView() + "\n" + m.viewport.View()
	case 3:
		return m.tabsView() + "\n" + m.recipeFilterView() + "\n" + m.viewport.View()
	default:
		return m.tabsView() + "\n" + m.viewport.View()
	}
}
