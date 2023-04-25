package details

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/commands"
	"github.com/tiramiseb/quickonf/instructions"
	"github.com/tiramiseb/quickonf/program/styles"
)

func (m *Model) View() string {
	var view string
	if m.group == nil {
		return ""
	}
	if len(m.group.Reports) == 0 {
		for _, ins := range m.group.Instructions {
			for _, rep := range ins.NotRunReports(0) {
				view += m.reportView(rep)
			}
		}
	} else {
		for _, rep := range m.group.Reports {
			view += m.reportView(rep)
		}
	}
	m.viewport.SetContent(view)
	return m.viewport.View()
}

func (m *Model) reportView(rep *instructions.CheckReport) string {
	status, message, level := rep.GetStatusAndMessage()
	indent := m.width / 20 * level
	content := fmt.Sprintf("[%s] %s", rep.Name, message)
	result := strings.Repeat(" ", indent) + styles.Styles[status].Render(
		styles.MakeWidth(content, m.width-indent),
	) + "\n"
	if status == commands.StatusInfo && m.showDetails {
		result += m.detailsView(rep)
	}
	return result
}

func (m *Model) detailsView(rep *instructions.CheckReport) string {
	switch {
	case rep.Before == "" && rep.After == "":
		return ""
	case rep.Before != "" && rep.After != "" && rep.Before != rep.After:
		leftWidth := (m.width - 1) / 2
		rightWidth := m.width - leftWidth - 1
		left := lipgloss.NewStyle().Width(leftWidth).Render(rep.Before)
		right := lipgloss.NewStyle().Width(rightWidth).Render(rep.After)
		height := lipgloss.Height(left)
		if rightHeight := lipgloss.Height(right); rightHeight > height {
			height = rightHeight
		}
		var separator string
		for i := 0; i < height-1; i++ {
			separator += "│\n"
		}
		separator += "│"
		return lipgloss.JoinHorizontal(lipgloss.Top, left, separator, right) + "\n"
	case rep.Before == "":
		return lipgloss.NewStyle().Width(m.width).Render(rep.After) + "\n"
	default:
		return lipgloss.NewStyle().Width(m.width).Render(rep.Before) + "\n"
	}
}
