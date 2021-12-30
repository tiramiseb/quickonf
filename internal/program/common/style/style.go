package style

import "github.com/charmbracelet/lipgloss"

const (
	bg     = lipgloss.Color("#222222")
	border = lipgloss.Color("#7777ff")
	text   = lipgloss.Color("#ffffff")

	titleBg    = lipgloss.Color("#555555")
	subtitleBg = lipgloss.Color("#5555ff")

	buttonBg = lipgloss.Color("#0000ff")
	buttonFg = lipgloss.Color("#ffff00")

	waitingBg         = lipgloss.Color("#555555")
	hoveredWaitingBg  = lipgloss.Color("#888888")
	selectedWaitingBg = lipgloss.Color("#777777")
	runningBg         = lipgloss.Color("#000099")
	hoveredRunningBg  = lipgloss.Color("#5555dd")
	selectedRunningBg = lipgloss.Color("#3333bb")
	failedBg          = lipgloss.Color("#ff0000")
	hoveredFailedBg   = lipgloss.Color("#ff3333")
	selectedFailedBg  = lipgloss.Color("#ff2222")
	successBg         = lipgloss.Color("#005500")
	hoveredSuccessBg  = lipgloss.Color("#338833")
	selectedSuccessBg = lipgloss.Color("#227722")
)

var (
	Title = lipgloss.NewStyle().
		Background(titleBg).
		Foreground(text).
		Bold(true)

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

	Button = lipgloss.NewStyle().
		Background(buttonBg).
		Foreground(buttonFg)

	ButtonKey = Button.Copy().
			Underline(true).
			Bold(true)

	ClickedButton = lipgloss.NewStyle().
			Background(buttonFg).
			Foreground(buttonBg)

	ClickedButtonKey = ClickedButton.Copy().
				Underline(true).
				Bold(true)
)
