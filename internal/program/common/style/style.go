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

	waitingBg = lipgloss.Color("#555555")
	runningBg = lipgloss.Color("#0000ff")
	failedBg  = lipgloss.Color("#ff0000")
	successBg = lipgloss.Color("#005500")
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
