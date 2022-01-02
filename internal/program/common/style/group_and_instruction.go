package style

import "github.com/charmbracelet/lipgloss"

var (
	baseGroup = lipgloss.NewStyle().
			Foreground(text).
			PaddingLeft(1)

	GroupWaiting = baseGroup.Copy().
			Background(waitingBg)
	GroupInfo = baseGroup.Copy().
			Background(infoBg)
	GroupRunning = baseGroup.Copy().
			Background(runningBg)
	GroupFail = baseGroup.Copy().
			Background(failedBg)
	GroupSuccess = baseGroup.Copy().
			Background(successBg)

	HoveredGroupWaiting = baseGroup.Copy().
				Background(hoveredWaitingBg)
	HoveredGroupInfo = baseGroup.Copy().
				Background(hoveredInfoBg)
	HoveredGroupRunning = baseGroup.Copy().
				Background(hoveredRunningBg)
	HoveredGroupFail = baseGroup.Copy().
				Background(hoveredFailedBg)
	HoveredGroupSuccess = baseGroup.Copy().
				Background(hoveredSuccessBg)

	SelectedGroupWaiting = baseGroup.Copy().
				Background(selectedWaitingBg)
	SelectedGroupInfo = baseGroup.Copy().
				Background(selectedInfoBg)
	SelectedGroupRunning = baseGroup.Copy().
				Background(selectedRunningBg)
	SelectedGroupFail = baseGroup.Copy().
				Background(selectedFailedBg)
	SelectedGroupSuccess = baseGroup.Copy().
				Background(selectedSuccessBg)

	baseInstruction = lipgloss.NewStyle().
			Foreground(text)
	InstructionInfo = baseInstruction.Copy().
			Background(infoBg)
	InstructionError = baseInstruction.Copy().
				Background(failedBg)
	InstructionSuccess = baseInstruction.Copy().
				Background(successBg)

	InstructionMessage = BoxContent.Copy()
)
