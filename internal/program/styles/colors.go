package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	BgColor = lipgloss.Color(fmt.Sprintf("%v", termenv.BackgroundColor()))
	FgColor = lipgloss.Color(fmt.Sprintf("%v", termenv.ForegroundColor()))
)
