package footer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/program/group"
	"github.com/tiramiseb/quickonf/program/style"
)

const Height = 1

type Model struct {
	nbWaiting   int
	nbRunning   int
	nbSucceeded int
	nbFailed    int

	width int
	View  string
}

func New(nb int) *Model {
	return &Model{nbWaiting: nb}
}

func (m *Model) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case group.RunningMessage:
		m.nbWaiting--
		m.nbRunning++
	case group.SucceededMessage:
		m.nbRunning--
		m.nbSucceeded++
	case group.FailedMessage:
		m.nbRunning--
		m.nbFailed++
	default:
		return
	}
	m.update()
}

func (m *Model) update() {
	status := fmt.Sprintf(" %d remaining -> %d running -> %d finished / %d failed", m.nbWaiting, m.nbRunning, m.nbSucceeded, m.nbFailed)
	if len(status) > m.width {
		status = fmt.Sprintf("%d rem -> %d run -> %d fin / %d fail", m.nbWaiting, m.nbRunning, m.nbSucceeded, m.nbFailed)
		if len(status) > m.width {
			status = fmt.Sprintf("%d->%d->%d->%d", m.nbWaiting, m.nbRunning, m.nbSucceeded, m.nbFailed)
			if len(status) > m.width {
				status = status[:m.width]
			}
		}
	}

	left := style.FooterLeft.Render(status)
	var (
		spacing = m.width - len(status)
		right   string
	)

	if m.width-len(status) > 11 {
		right = style.FooterRight.Render(" [h: help] ")
		spacing = m.width - len(status) - 11
	}
	center := style.Footer.Render(strings.Repeat(" ", spacing))
	m.View = left + center + right
}
