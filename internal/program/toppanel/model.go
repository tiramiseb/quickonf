package toppanel

import (
	_ "embed"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

const showHideDuration = 10 * time.Millisecond

type showMessage struct{}

type hideMessage struct{}

var (
	//go:embed help.txt
	helpSrc string
	help    = strings.Split(strings.TrimSpace(helpSrc), "\n")
)

type displaying int

const (
	displayingNone displaying = iota
	displayingHelp
)

type Model struct {
	width int

	displaying    displaying
	targetContent []string
	moving        bool
	partialLines  int

	Size int
	View string
}

func New(width int) *Model {
	return &Model{width: width}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if m.displaying == displayingHelp {
				cmd = m.startHide
			} else {
				cmd = m.startShowHelp
			}
		}
	case showMessage:
		cmd = m.waitShow
	case hideMessage:
		cmd = m.waitHide
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	m.update()
	return cmd
}

func (m *Model) waitHide() tea.Msg {
	time.Sleep(showHideDuration)
	return m.hide()
}

func (m *Model) waitShow() tea.Msg {
	time.Sleep(showHideDuration)
	return m.show()
}

func (m *Model) startHide() tea.Msg {
	if m.moving {
		return nil
	}
	m.partialLines = len(m.targetContent) + 1
	m.moving = true
	return m.hide()
}

func (m *Model) startShowHelp() tea.Msg {
	if m.moving {
		return nil
	}
	m.partialLines = 0
	m.displaying = displayingHelp
	m.targetContent = help
	m.moving = true
	return m.show()
}

func (m *Model) hide() tea.Msg {
	m.partialLines--
	if m.partialLines == 0 {
		m.moving = false
		m.displaying = displayingNone
		return nil
	}
	return hideMessage{}
}

func (m *Model) show() tea.Msg {
	m.partialLines++
	if m.partialLines > len(m.targetContent) {
		m.moving = false
		m.partialLines = 0
		return nil
	}
	return showMessage{}
}

func (m *Model) update() {
	if m.displaying == displayingNone {
		m.View = ""
		m.Size = 0
	} else if m.moving {
		content := m.targetContent[len(m.targetContent)-m.partialLines:]
		m.View = style.TopPanelMoving.Width(m.width - 2).Render(
			strings.Join(content, "\n"),
		)
		m.Size = len(content)
	} else {
		m.View = style.TopPanel.Width(m.width - 2).Render(
			strings.Join(m.targetContent, "\n"),
		)
		m.Size = len(m.targetContent) + 1
	}
}

// func (m *Model) hideHelp() tea.Msg {
// 	m.displaying = displayingNone
// 	return nil
// }

// func (m *Model) showHelp() tea.Msg {
// 	m.displaying = displayingHelp
// 	return nil
// }

// func (m *Model) Size() int {
// 	if m.displaying == displayingNone {
// 		return 0
// 	}
// }

// func (m *Model) View() string {
// 	lines := m.content()
// 	if m.partialLines != 0 {
// 		lines = lines[len(lines)-m.partialLines:]
// 	}
// 	return strings.Join(lines, "\n")
// }

// func (m *Model) content() []string {
// 	var content string
// 	if m.displaying == displayingHelp {
// 		content = help
// 	}
// 	return strings.Split(
// 		style.TopPanel.Width(m.width-2).Render(content),
// 		"\n",
// 	)

// }
