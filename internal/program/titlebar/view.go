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
	title := " Quickonf "
	titleWidth := len(title)
	m.quitStart = -1
	m.quitEnd = -1
	m.helpStart = -1
	m.helpEnd = -1
	m.detailsStart = -1
	m.detailsEnd = -1
	m.filterStart = -1
	m.filterEnd = -1
	availableWidth := width - titleWidth

	// No place even for the title, cut it
	if availableWidth <= 0 {
		view := style.Render(title[:width])
		m.View = func() string {
			return view
		}
		m.HelpView = m.View
		return
	}

	leftPart := func(availableWidth int) string {
		return style.Render(title + strings.Repeat(" ", availableWidth))
	}

	// No place for the quit button, only include the title
	if availableWidth <= m.quit.Width {
		view := leftPart(availableWidth)
		m.View = func() string {
			return view
		}
		return
	}
	availableWidth = availableWidth - m.quit.Width - 1
	m.quitEnd = width - 2
	m.quitStart = m.quitEnd - m.quit.Width + 1

	// No place for the help button, include title & quit
	if availableWidth <= m.help.Width {
		view := leftPart(availableWidth) + m.quit.View + space
		m.View = func() string {
			return view
		}
		return
	}
	availableWidth = availableWidth - m.help.Width - 1
	m.helpEnd = m.quitStart - 2
	m.helpStart = m.helpEnd - m.help.Width + 1

	// No place for the filter button, include title & help & quit
	if availableWidth <= m.filter.Width {
		view := leftPart(availableWidth) + m.help.View + space + m.quit.View + space
		m.View = func() string {
			return view
		}
		return
	}
	availableWidth = availableWidth - m.filter.Width - 1

	endOfTitle := space + m.help.View + space + m.quit.View + space

	// No place for the details button, include title & filter & help & quit
	if availableWidth <= m.details.Width {
		m.filterEnd = m.helpStart - 2
		m.filterStart = m.filterEnd - m.filter.Width + 1
		m.View = func() string {
			return leftPart(availableWidth) + m.filter.View() + endOfTitle
		}
		return
	}
	availableWidth = availableWidth - m.details.Width - 1
	m.detailsEnd = m.helpStart - 2
	m.detailsStart = m.detailsEnd - m.details.Width + 1
	m.filterEnd = m.detailsStart - 2
	m.filterStart = m.filterEnd - m.filter.Width + 1

	m.View = func() string {
		return leftPart(availableWidth) + m.filter.View() + space + m.details.View() + endOfTitle
	}
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
