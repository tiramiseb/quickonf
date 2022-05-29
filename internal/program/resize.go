package program

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) resize(size tea.WindowSizeMsg) {
	width := size.Width - 1
	height := size.Height - 3
	left := tea.WindowSizeMsg{
		Width:  width / 2,
		Height: height,
	}
	if left.Width > m.groups.MaxNameLength()+2 {
		left.Width = m.groups.MaxNameLength() + 2
	}
	m.separatorXPos = left.Width
	right := tea.WindowSizeMsg{
		Width:  width - left.Width,
		Height: height,
	}

	m.leftTitle = subtitleStyle.Width(left.Width).Render("Groups")
	m.rightTitle = subtitleStyle.Width(right.Width).Render("Details")
	m.leftTitleWithFocus = subtitleWithFocusStyle.Width(left.Width).Render("Groups")
	m.rightTitleWithFocus = subtitleWithFocusStyle.Width(right.Width).Render("Details")

	m.subtitlesSeparator = strings.Repeat("─", left.Width) + "┼" + strings.Repeat("─", right.Width)

	var separator string
	for i := 0; i < height-1; i++ {
		separator += "│\n"
	}
	separator += "│"
	m.verticalSeparator = separator

	m.titlebar = m.titlebar.Resize(size)
	m.groupsview = m.groupsview.Resize(left)
	m.details = m.details.Resize(right)
	m.help = m.help.Resize(tea.WindowSizeMsg{
		Width:  size.Width,
		Height: size.Height - 1,
	})
}
