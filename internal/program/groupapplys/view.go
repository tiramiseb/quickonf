package groupapplys

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/apply"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

func (m *model) windowSize(msg tea.WindowSizeMsg) tea.Cmd {
	newMsg := tea.WindowSizeMsg{
		Height: msg.Height - 1,
		Width:  msg.Width - 2,
	}
	m.width = msg.Width
	m.boxHeight = newMsg.Height - 2
	if m.boxHeight < 0 {
		m.boxHeight = 0
	}
	m.updateActive()
	apply.GroupStyles = map[apply.Status]lipgloss.Style{
		apply.StatusWaiting:   style.GroupWaiting.Copy().Width(newMsg.Width),
		apply.StatusRunning:   style.GroupRunning.Copy().Width(newMsg.Width),
		apply.StatusFailed:    style.GroupFail.Copy().Width(newMsg.Width),
		apply.StatusSucceeded: style.GroupSuccess.Copy().Width(newMsg.Width),
	}
	cmds := make([]tea.Cmd, len(m.groups))
	for i, g := range m.groups {
		m.groups[i], cmds[i] = g.Update(newMsg)
	}
	return tea.Batch(cmds...)
}

// redrawContent draws content for all groups, even if they are not displayed
func (m *model) redrawContent() {
	result := make([]string, 0, len(m.groups)*2)
	lineToGroup := make([]groupLine, 0, len(m.groups)*2)
	for i, g := range m.groups {
		groupView := strings.Split(g.View(), "\n")
		thisLineToGroup := make([]groupLine, len(groupView))
		for j := range thisLineToGroup {
			thisLineToGroup[j] = groupLine{i, j}
		}
		result = append(result, groupView...)
		lineToGroup = append(lineToGroup, thisLineToGroup...)
	}
	m.allGroupsView = result
	m.allLineToGroup = lineToGroup
}

// updateView extracts lines to display from content for all groups
func (m *model) updateView() {

	// Find where is the selected group
	startOfSelectedGroup := -1
	sizeOfSelectedGroup := -1
findGroup:
	for i, groupline := range m.allLineToGroup {
		switch {
		case startOfSelectedGroup < 0 && groupline.gidx == m.selectedGroup:
			startOfSelectedGroup = i
		case startOfSelectedGroup >= 0 && groupline.gidx != m.selectedGroup:
			sizeOfSelectedGroup = i - startOfSelectedGroup
			break findGroup
		}
	}

	// Put the selected group at the top of the view
	firstLineInView := startOfSelectedGroup
	lastLineInView := firstLineInView + m.boxHeight

	// Slide up if at the bottom of the list
	if lastLineInView >= len(m.allLineToGroup) {
		firstLineInView -= (lastLineInView - len(m.allLineToGroup))
		if firstLineInView < 0 {
			firstLineInView = 0
		}
		lastLineInView = firstLineInView + m.boxHeight
		if lastLineInView >= len(m.allLineToGroup) {
			lastLineInView = len(m.allLineToGroup)
		}
	}

	// Check if the view could be slided up to put the group in the middle
	idealLinesAbove := (m.boxHeight - sizeOfSelectedGroup) / 2
	idealFirstLineInVew := startOfSelectedGroup - idealLinesAbove
	if idealFirstLineInVew < 0 {
		idealFirstLineInVew = 0
	}
	if idealFirstLineInVew < firstLineInView {
		delta := firstLineInView - idealFirstLineInVew
		firstLineInView = idealFirstLineInVew
		lastLineInView -= delta
	}

	// Extract the view
	m.view = strings.Join(m.allGroupsView[firstLineInView:lastLineInView], "\n")
	m.viewLineToGroup = m.allLineToGroup[firstLineInView:lastLineInView]
	m.selectedGroupFirstLine = startOfSelectedGroup - firstLineInView

}

func (m *model) updateActive() {
	if m.active {
		m.subtitleStyle = style.ActiveBoxTitle.Copy().Width(m.width)
		m.boxStyle = style.ActiveBox.Copy().Width(m.width - 2).Height(m.boxHeight)
	} else {
		m.subtitleStyle = style.BoxTitle.Copy().Width(m.width)
		m.boxStyle = style.Box.Copy().Width(m.width - 2).Height(m.boxHeight)
	}
}

func (m *model) View() string {
	return m.subtitleStyle.Render("Check") + "\n" + m.boxStyle.Render(m.view)
}
