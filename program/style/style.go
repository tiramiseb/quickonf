package style

import "github.com/charmbracelet/lipgloss"

var (
	Main = lipgloss.NewStyle().
		Background(lipgloss.Color("#222222")).
		Foreground(lipgloss.Color("#eeeeee"))

	Header = lipgloss.NewStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Background(lipgloss.Color("#666666")).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Align(lipgloss.Center)

	TopPanel = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderBackground(lipgloss.Color("#000000")).
			BorderForeground(lipgloss.Color("#8888ff")).
			Background(lipgloss.Color("#000000")).
			Foreground(lipgloss.Color("#aaaaff"))

	Footer = lipgloss.NewStyle().
		Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#ffffff"))

	FooterLeft = Footer.Copy().
			Italic(true)

	FooterRight = Footer.Copy().
			Bold(true)
)
