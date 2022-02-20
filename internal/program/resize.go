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
	right := tea.WindowSizeMsg{
		Width:  width - left.Width,
		Height: height,
	}

	leftTitle := subtitleStyle.Width(left.Width).Render("Checks")
	rightTitle := subtitleStyle.Width(right.Width).Render("Details")

	m.subtitles = leftTitle + "│" + rightTitle + "\n" + strings.Repeat("─", left.Width) + "┼" + strings.Repeat("─", right.Width) + "\n"

	var separator string
	for i := 0; i < height-1; i++ {
		separator += "│\n"
	}
	separator += "│"
	m.separator = separator

	m.titlebar = m.titlebar.Resize(size)
	m.checks.Resize(left)
	m.details.Resize(right)
}
