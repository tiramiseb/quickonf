package program

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func apply(i int) tea.Cmd {
	return func() tea.Msg {
		go global.DisplayedGroups[i].Apply()
		return nil
	}
}

func (m *model) listenSignal() tea.Msg {
	log.Print("WAITING SIGNAL")
	<-m.signalTarget
	log.Print("RECEIVED SIGNAL")
	return newSignal{}
}

type newSignal struct{}
