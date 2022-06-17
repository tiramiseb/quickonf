package button

import "github.com/charmbracelet/lipgloss"

const (
	bgColor = lipgloss.Color("#0000ff")
	fgColor = lipgloss.Color("#ffff00")
)

var (
	inactiveStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(fgColor)
	inactiveHintStyle = inactiveStyle.Copy().
				Underline(true).
				Bold(true)
	activeStyle = lipgloss.NewStyle().
			Background(fgColor).
			Foreground(bgColor)
	activeHintStyle = activeStyle.Copy().
			Underline(true).
			Bold(true)
)
