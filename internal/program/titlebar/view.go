package titlebar

import "strings"

func (m *Model) draw(width int) {
	title := " Quickonf "
	titleWidth := len(title)
	m.quitStart = -1
	m.quitEnd = -1
	m.helpStart = -1
	m.helpEnd = -1
	m.filterStart = -1
	m.filterEnd = -1
	// m.toggleStart = -1
	// m.toggleEnd = -1
	availableWidth := width - titleWidth

	// No place even for the title, cut it
	if availableWidth <= 0 {
		view := style.Render(title[:width])
		m.view = func() string {
			return view
		}
		m.helpView = m.view
		return
	}

	leftPart := func(availableWidth int) string {
		return style.Render(title + strings.Repeat(" ", availableWidth))
	}

	// No place for the quit button, only include the title
	if availableWidth <= m.quit.Width {
		view := leftPart(availableWidth)
		m.view = func() string {
			return view
		}
		m.helpView = m.view
		return
	}
	space := style.Render(" ")
	availableWidth = availableWidth - m.quit.Width - 1
	m.quitEnd = width - 2
	m.quitStart = m.quitEnd - m.quit.Width + 1

	// No place for the help button, include title & quit
	if availableWidth <= m.help.Width {
		view := leftPart(availableWidth) + m.quit.View + space
		m.view = func() string {
			return view
		}
		m.helpView = m.view
		return
	}
	availableWidth = availableWidth - m.help.Width - 1
	m.helpEnd = m.quitStart - 2
	m.helpStart = m.helpEnd - m.help.Width + 1

	endOfTitle := space + m.quit.View + space

	// From now on, when help is displayed, need only help & quit
	helpLeft := leftPart(availableWidth)
	m.helpView = func() string {
		return helpLeft + m.help.View() + endOfTitle
	}

	// No place for the filter button, include title & help & quit
	if availableWidth <= m.filter.Width {
		m.view = func() string {
			return helpLeft + m.help.View() + endOfTitle
		}
		return
	}
	availableWidth = availableWidth - m.filter.Width - 1
	m.filterEnd = m.helpStart - 2
	m.filterStart = m.filterEnd - m.filter.Width + 1

	left := leftPart(availableWidth)
	m.view = func() string {
		return left + m.filter.View() + space + m.help.View() + endOfTitle
	}

	// No place for the toggle button, include title & filter & help & quit
	// buttonWidth = lipgloss.Width(m.toggle.View())
	// if availableWidth <= m.filter.Width {
	// 	leftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
	// 	m.view = func() string {
	// 		return leftAndSpace + m.filter.View() + space + m.help.View() + space + m.quit.View() + space
	// 	}
	// 	return
	// }
	// availableWidth = availableWidth - buttonWidth - 1
	// m.toggleEnd = m.filterStart - 2
	// m.toggleStart = m.toggleEnd - buttonWidth + 1

	// leftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
	// m.view = func() string {
	// 	return leftAndSpace + m.toggle.View() + space + m.filter.View() + space + m.help.View() + space + m.quit.View() + space
	// }
}

func (m *Model) View() string {
	if m.isInHelp {
		return m.helpView()
	}
	return m.view()
}
