package checks

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/commands"
)

var (
	styles = map[commands.Status]lipgloss.Style{
		commands.StatusNotRun:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#555555")),
		commands.StatusInfo:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#000099")),
		commands.StatusRunning: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#555500")),
		commands.StatusSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#005500")),
		commands.StatusError:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#ff0000")),
	}
	selectedStyles = map[commands.Status]lipgloss.Style{
		commands.StatusNotRun:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#777777")),
		commands.StatusInfo:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#3333bb")),
		commands.StatusRunning: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#777722")),
		commands.StatusSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#227722")),
		commands.StatusError:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#ff2222")),
	}
)
