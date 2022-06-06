package help

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/styles"
)

var (
	subtitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	subtitleWithFocusStyle = subtitleStyle.Copy().
				Background(styles.FgColor).
				Foreground(styles.BgColor).
				Bold(true)
)
