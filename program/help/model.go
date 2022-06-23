package help

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

const filterMaxLength = 20

type Model struct {
	viewport viewport.Model

	height int
	width  int

	introTitle             string
	introTitleWithFocus    string
	introStart             int
	introEnd               int
	languageTitle          string
	languageTitleWithFocus string
	languageStart          int
	languageEnd            int
	commandsTitle          string
	commandsTitleWithFocus string
	commandsStart          int
	commandsEnd            int
	uiTitle                string
	uiTitleWithFocus       string
	uiStart                int
	uiEnd                  int
	subtitleSeparator      string

	activeSection int
	commandFilter string
}

func New() *Model {
	v := viewport.Model{Width: 1, Height: 1}
	v.Style = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	return &Model{
		viewport:      v,
		activeSection: 0,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "left":
			if m.activeSection > 0 {
				m.activeSection--
				m.setContent()
			}
		case "right":
			if m.activeSection < 3 {
				m.activeSection++
				m.setContent()
			}
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
			"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			".":
			if m.activeSection == 2 && len(m.commandFilter) < filterMaxLength {
				m.commandFilter += strings.ToLower(key)
				m.setContent()
			}
		case "backspace":
			if m.activeSection == 2 && m.commandFilter != "" {
				m.commandFilter = m.commandFilter[:len(m.commandFilter)-1]
				m.setContent()
			}
		default:
			m.viewport, cmd = m.viewport.Update(msg)
		}
	case tea.MouseMsg:
		if msg.Y == 0 && msg.Type == tea.MouseRelease {
			switch {
			case msg.X >= m.introStart && msg.X <= m.introEnd && m.activeSection != 0:
				m.activeSection = 0
				m.setContent()
			case msg.X >= m.languageStart && msg.X <= m.languageEnd && m.activeSection != 1:
				m.activeSection = 1
				m.setContent()
			case msg.X >= m.commandsStart && msg.X <= m.commandsEnd && m.activeSection != 2:
				m.activeSection = 2
				m.setContent()
			case msg.X >= m.uiStart && msg.X <= m.uiEnd && m.activeSection != 3:
				m.activeSection = 3
				m.setContent()
			}
		} else {
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}
	return m, cmd
}

func (m *Model) setContent() {
	m.viewport.Width = m.width
	var text string
	if lipgloss.HasDarkBackground() {
		switch m.activeSection {
		case 0:
			text = gunzip(introDark)
			m.viewport.Height = m.height - 2
		case 1:
			text = gunzip(languageDark)
			m.viewport.Height = m.height - 2
		case 2:
			m.viewport.Height = m.height - 4
			text = m.commandsDoc(true)
		case 3:
			text = gunzip(uiDark)
			m.viewport.Height = m.height - 2
		}
	} else {
		switch m.activeSection {
		case 0:
			text = gunzip(introLight)
			m.viewport.Height = m.height - 2
		case 1:
			text = gunzip(languageLight)
			m.viewport.Height = m.height - 2
		case 2:
			m.viewport.Height = m.height - 4
			text = m.commandsDoc(false)
		case 3:
			text = gunzip(uiLight)
			m.viewport.Height = m.height - 2
		}
	}
	m.viewport.GotoTop()
	m.viewport.SetContent(
		wordwrap.String(text, m.width-4),
	)
}

func gunzip(gzipped []byte) string {
	buf := bytes.NewReader(gzipped)
	r, _ := gzip.NewReader(buf)
	defer r.Close()
	unzipped, _ := io.ReadAll(r)
	return string(unzipped)

}
