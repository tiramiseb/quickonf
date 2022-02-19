package style

import "github.com/charmbracelet/lipgloss"

const (
	border = lipgloss.Color("#7777ff")

	waitingBg         = lipgloss.Color("#555555")
	hoveredWaitingBg  = lipgloss.Color("#888888")
	selectedWaitingBg = lipgloss.Color("#777777")
	runningBg         = lipgloss.Color("#555500")
	hoveredRunningBg  = lipgloss.Color("#888833")
	selectedRunningBg = lipgloss.Color("#777722")
	infoBg            = lipgloss.Color("#000099")
	hoveredInfoBg     = lipgloss.Color("#5555dd")
	selectedInfoBg    = lipgloss.Color("#3333bb")
	failedBg          = lipgloss.Color("#ff0000")
	hoveredFailedBg   = lipgloss.Color("#ff3333")
	selectedFailedBg  = lipgloss.Color("#ff2222")
	successBg         = lipgloss.Color("#005500")
	hoveredSuccessBg  = lipgloss.Color("#338833")
	selectedSuccessBg = lipgloss.Color("#227722")
)

var (
	BoxTitle = lipgloss.NewStyle().
			Background(subtitleBg).
			Foreground(text).
			Align(lipgloss.Center).
			Bold(true).
			Height(1)

	ActiveBoxTitle = BoxTitle.Copy().
			Background(text).
			Foreground(subtitleBg)

	BoxContent = lipgloss.NewStyle().
			Background(bg).
			Foreground(text)

	Box = BoxContent.Copy().
		Border(lipgloss.DoubleBorder()).
		BorderBackground(bg).
		BorderForeground(border)

	ActiveBox = Box.Copy().
			BorderForeground(text)
)
