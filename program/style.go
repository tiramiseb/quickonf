package program

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/program/styles"
)

var (
	subtitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	subtitleWithFocusStyle = subtitleStyle.Copy().
				Background(styles.FgColor).
				Foreground(styles.BgColor).
				Bold(true)
)
