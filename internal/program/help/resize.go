package help

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.height = size.Height
	m.width = size.Width
	m.updateTabs()
	m.setContent()
	return m
}

func (m *Model) updateTabs() {
	if m.width < 15 {
		m.introStart = -1
		m.introEnd = -1
		m.languageStart = -1
		m.languageEnd = -1
		m.commandsStart = -1
		m.commandsEnd = -1
		m.uiStart = -1
		m.uiEnd = -1
	}
	// 3 chars for separators, 2 chars par section for spacing, total = 11
	introWidth := (m.width - 11) / 4
	languageWidth := introWidth
	commandsWidth := introWidth
	uiWidth := introWidth
	switch m.width - introWidth*4 - 11 {
	case 3:
		languageWidth++
		commandsWidth++
		uiWidth++
	case 2:
		languageWidth++
		uiWidth++
	case 1:
		languageWidth++
	}

	var (
		introText    string
		languageText string
		commandsText string
		uiText       string
	)
	switch {
	case introWidth >= 12:
		introText = "Introduction"
	case introWidth >= 5:
		introText = "Intro"
	default:
		introText = "I"
	}

	switch {
	case languageWidth >= 22:
		languageText = "Configuration language"
	case languageWidth >= 13:
		languageText = "Conf language"
	case languageWidth >= 8:
		languageText = "Language"
	case languageWidth >= 4:
		languageText = "Lang"
	default:
		languageText = "L"
	}

	switch {
	case commandsWidth >= 8:
		commandsText = "Commands"
	case commandsWidth >= 4:
		commandsText = "Cmds"
	default:
		commandsText = "C"
	}

	switch {
	case uiWidth >= 14:
		uiText = "User interface"
	case uiWidth >= 9:
		uiText = "Interface"
	case uiWidth >= 5:
		uiText = "Iface"
	case uiWidth >= 2:
		uiText = "UI"
	default:
		uiText = "U"
	}

	m.introStart = 0
	m.introEnd = m.introStart + introWidth + 1
	m.introTitle = subtitleStyle.Width(introWidth + 2).Render(introText)
	m.introTitleWithFocus = subtitleWithFocusStyle.Width(introWidth + 2).Render(introText)

	m.languageStart = m.introEnd + 2
	m.languageEnd = m.languageStart + languageWidth + 1
	m.languageTitle = subtitleStyle.Width(languageWidth + 2).Render(languageText)
	m.languageTitleWithFocus = subtitleWithFocusStyle.Width(languageWidth + 2).Render(languageText)

	m.commandsStart = m.languageEnd + 2
	m.commandsEnd = m.commandsStart + commandsWidth + 1
	m.commandsTitle = subtitleStyle.Width(commandsWidth + 2).Render(commandsText)
	m.commandsTitleWithFocus = subtitleWithFocusStyle.Width(commandsWidth + 2).Render(commandsText)

	m.uiStart = m.commandsEnd + 2
	m.uiEnd = m.uiStart + uiWidth + 1
	m.uiTitle = subtitleStyle.Width(uiWidth + 2).Render(uiText)
	m.uiTitleWithFocus = subtitleWithFocusStyle.Width(uiWidth + 2).Render(uiText)

	m.subtitleSeparator = strings.Repeat("─", introWidth+2) + "┴" + strings.Repeat("─", languageWidth+2) + "┴" + strings.Repeat("─", commandsWidth+2) + "┴" + strings.Repeat("─", uiWidth+2)
}
