package header

import "github.com/charmbracelet/lipgloss"

const Height = 3

var headerStyle = lipgloss.NewStyle().
	PaddingTop(1).
	PaddingBottom(1).
	Background(lipgloss.Color("#666666")).
	Foreground(lipgloss.Color("#ffffff")).
	Bold(true).
	Align(lipgloss.Center)
