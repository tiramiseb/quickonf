package titlebar

import (
	"strings"
)

var space = style.Render(" ")

func (m *Model) draw(width int) {
	m.drawRegular(width)
	m.drawHelp(width)
}

func (m *Model) drawRegular(width int) {
	availableWidth := width - len(m.title)

	neededWidth := m.quit.Width + 1
	m.showQuit = availableWidth >= neededWidth

	neededWidth += m.help.Width + 1
	m.showHelp = availableWidth >= neededWidth

	neededWidth += m.filter.Width + 1
	m.showFilter = availableWidth >= neededWidth

	neededWidth += m.apply.Width + 1
	m.showApply = availableWidth >= neededWidth

	neededWidth += m.applyAll.Width + 1
	m.showApplyAll = availableWidth >= neededWidth

	neededWidth += m.details.Width + 1
	m.showDetails = availableWidth >= neededWidth

	neededWidth += m.recheck.Width + 1
	m.showRecheck = availableWidth >= neededWidth

	neededWidth += m.reloadConfig.Width + 1
	m.showReloadConfig = availableWidth >= neededWidth

	rightBorder := width

	if m.showQuit {
		m.quitEnd = rightBorder - 2
		m.quitStart = m.quitEnd - m.quit.Width + 1
		rightBorder = m.quitStart
	} else {
		m.quitEnd = -1
		m.quitStart = -1
	}

	if m.showReloadConfig {
		m.reloadConfigEnd = rightBorder - 2
		m.reloadConfigStart = m.reloadConfigEnd - m.reloadConfig.Width + 1
		rightBorder = m.reloadConfigStart
	} else {
		m.reloadConfigEnd = -1
		m.reloadConfigStart = -1
	}

	if m.showHelp {
		m.helpEnd = rightBorder - 2
		m.helpStart = m.helpEnd - m.help.Width + 1
		rightBorder = m.helpStart
	} else {
		m.helpEnd = -1
		m.helpStart = -1
	}

	if m.showDetails {
		m.detailsEnd = rightBorder - 2
		m.detailsStart = m.detailsEnd - m.details.Width + 1
		rightBorder = m.detailsStart
	} else {
		m.detailsEnd = -1
		m.detailsStart = -1
	}

	if m.showFilter {
		m.filterEnd = rightBorder - 2
		m.filterStart = m.filterEnd - m.filter.Width + 1
		rightBorder = m.filterStart
	} else {
		m.filterEnd = -1
		m.filterStart = -1
	}

	if m.showApplyAll {
		m.applyAllEnd = rightBorder - 2
		m.applyAllStart = m.applyAllEnd - m.applyAll.Width + 1
		rightBorder = m.applyAllStart
	} else {
		m.applyAllEnd = -1
		m.applyAllStart = -1
	}

	if m.showRecheck {
		m.recheckEnd = rightBorder - 2
		m.recheckStart = m.recheckEnd - m.recheck.Width + 1
		rightBorder = m.recheckStart
	} else {
		m.recheckEnd = -1
		m.recheckStart = -1
	}

	if m.showApply {
		m.applyEnd = rightBorder - 2
		m.applyStart = m.applyEnd - m.apply.Width + 1
		rightBorder = m.applyStart
	} else {
		m.applyEnd = -1
		m.applyStart = -1
	}

	m.middleSpace = rightBorder - len(m.title)
}

func (m *Model) View() string {
	var title strings.Builder

	title.WriteString(style.Render(m.title + strings.Repeat(" ", m.middleSpace)))
	if m.showApply {
		title.WriteString(m.apply.View)
		title.WriteString(space)
	}
	if m.showRecheck {
		title.WriteString(m.recheck.View)
		title.WriteString(space)
	}
	if m.showApplyAll {
		title.WriteString(m.applyAll.View())
		title.WriteString(space)
	}
	if m.showFilter {
		title.WriteString(m.filter.View())
		title.WriteString(space)
	}
	if m.showDetails {
		title.WriteString(m.details.View())
		title.WriteString(space)
	}
	if m.showHelp {
		title.WriteString(m.help.View)
		title.WriteString(space)
	}
	if m.showReloadConfig {
		title.WriteString(m.reloadConfig.View())
		title.WriteString(space)
	}
	if m.showQuit {
		title.WriteString(m.quit.View)
		title.WriteString(space)
	}
	return title.String()
}

func (m *Model) drawHelp(width int) {
	title := " Quickonf - Help "
	titleWidth := len(title)
	m.helpBackStart = -1
	m.helpBackEnd = -1
	availableWidth := width - titleWidth

	// No place even for the title, cut it
	if availableWidth <= 0 {
		view := style.Render(title[:width])
		m.HelpView = func() string {
			return view
		}
		return
	}

	leftPart := func(availableWidth int) string {
		return style.Render(title + strings.Repeat(" ", availableWidth))
	}

	// No place for the back button, only include the title
	if availableWidth <= m.helpBack.Width {
		view := leftPart(availableWidth)
		m.HelpView = func() string {
			return view
		}
		return
	}
	availableWidth = availableWidth - m.helpBack.Width - 1
	m.helpBackEnd = width - 2
	m.helpBackStart = m.helpBackEnd - m.helpBack.Width + 1

	view := leftPart(availableWidth) + m.helpBack.View + space
	m.HelpView = func() string {
		return view
	}
}
