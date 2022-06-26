package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/commands"
)

var (
	Styles = map[commands.Status]lipgloss.Style{
		commands.StatusNotRun:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#555555")),
		commands.StatusInfo:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#000099")),
		commands.StatusRunning: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#555500")),
		commands.StatusSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#005500")),
		commands.StatusError:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#990000")),
	}
	SelectedStyles = map[commands.Status]lipgloss.Style{
		commands.StatusNotRun:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#777777")).Bold(true),
		commands.StatusInfo:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#3333bb")).Bold(true),
		commands.StatusRunning: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#777722")).Bold(true),
		commands.StatusSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#227722")).Bold(true),
		commands.StatusError:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#bb3333")).Bold(true),
	}
)

func MakeWidth(str string, width int, selected ...bool) string {
	if width <= 0 {
		return ""
	}
	var line string
	if len(selected) > 0 && selected[0] {
		line = "> " + str
	} else {
		line = " " + str
	}
	remaining := width - lipgloss.Width(line)
	if remaining < 0 {
		return line[:width]
	}
	return line + strings.Repeat(" ", remaining)
}
