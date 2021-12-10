package group

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			PaddingLeft(2)
	runningStyle = baseStyle.Copy().
			Background(lipgloss.Color("#0000ff"))
	waitingStyle = baseStyle.Copy().
			Background(lipgloss.Color("#555555"))
	failedStyle = baseStyle.Copy().
			Background(lipgloss.Color("#ff0000"))
	succeededStyle = baseStyle.Copy().
			Background(lipgloss.Color("#005500"))
	styleMap = map[status]lipgloss.Style{
		statusWaiting:   waitingStyle,
		statusRunning:   runningStyle,
		statusFailed:    failedStyle,
		statusSucceeded: succeededStyle,
	}

	baseInstructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				PaddingLeft(2).
				PaddingRight(2)
	infoStyle = baseInstructionStyle.Copy().
			Background(lipgloss.Color("#0000ff"))
	successStyle = baseInstructionStyle.Copy().
			Background(lipgloss.Color("#005500"))
	errorStyle = baseInstructionStyle.Copy().
			Background(lipgloss.Color("#ff0000"))
	instructionStyleMap = map[instructionStatus]lipgloss.Style{
		instructionInfo:    infoStyle,
		instructionSuccess: successStyle,
		instructionError:   errorStyle,
	}

	instructionMessageStyle = lipgloss.NewStyle().
				PaddingLeft(2)
)
