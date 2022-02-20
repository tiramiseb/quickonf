package program

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

var (
	subtitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	subtitleWithFocusStyle = subtitleStyle.Copy().
				Background(global.FgColor).
				Foreground(global.BgColor).
				Bold(true)
)
