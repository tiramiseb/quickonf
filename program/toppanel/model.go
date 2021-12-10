package toppanel

import (
	_ "embed"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/program/style"
)

var (
	//go:embed help.txt
	helpSrc  string
	help     = strings.TrimRight(helpSrc, "\n")
	helpSize = strings.Count(help, "\n")
)

type displaying int

const (
	displayingNone displaying = iota
	displayingHelp
)

type Model struct {
	width      int
	displaying displaying

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
				cmd = m.hideHelp
			} else {
				cmd = m.showHelp
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return cmd
}

func (m *Model) hideHelp() tea.Msg {
	m.displaying = displayingNone
	m.Size = 0
	m.View = ""
	return nil
}

func (m *Model) showHelp() tea.Msg {
	m.displaying = displayingHelp
	view := style.TopPanel.
		Width(m.width - 2).
		Render(help)
	m.Size = helpSize + 2
	m.View = view
	return nil
}
