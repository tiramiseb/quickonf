package help

import (
	_ "embed"

	"github.com/charmbracelet/glamour"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

var (
	//go:embed intro.md
	introMd string

	//go:embed language.md
	languageMd string

	//go:embed ui.md
	uiMd string

	backgroundColor      = style.BgHTML
	margin          uint = 1
)

func init() {
	glamour.DarkStyleConfig.Document.BackgroundColor = &backgroundColor
	glamour.DarkStyleConfig.Document.Margin = &margin
}

func (m *model) render(content string) string {
	result, err := m.mdRenderer.Render(content)
	if err != nil {
		result = err.Error() + "\n" + result
	}
	return result
}

func (m *model) reRenderContents() {
	// Reset the filtered commands doc when the window is resized
	m.filteredCommandsDoc = map[string]string{}
	m.contents[sectionIntro] = m.render(introMd)
	m.contents[sectionLang] = m.render(languageMd)
	m.contents[sectionCommands] = m.commandsDoc()
	m.contents[sectionUI] = m.render(uiMd)
	m.viewport.SetContent(m.contents[m.active])
}
