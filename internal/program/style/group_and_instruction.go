package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	baseGroup = lipgloss.NewStyle().
			Foreground(text).
			PaddingLeft(1)
	GroupWaiting = baseGroup.Copy().
			Background(waitingBg)
	GroupRunning = baseGroup.Copy().
			Background(runningBg)
	GroupFail = baseGroup.Copy().
			Background(failedBg)
	GroupSuccess = baseGroup.Copy().
			Background(successBg)

	baseInstruction = lipgloss.NewStyle().
			Foreground(text)
	InstructionInfo = baseInstruction.Copy().
			Background(runningBg)
	InstructionError = baseInstruction.Copy().
				Background(failedBg)
	InstructionSuccess = baseInstruction.Copy().
				Background(successBg)

	InstructionMessage = BoxContent.Copy()
)
