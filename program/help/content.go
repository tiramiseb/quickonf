package help

import (
	"embed"
)

var (
	//go:embed content/intro.dark.msg
	introDark string

	//go:embed content/intro.light.msg
	introLight string

	//go:embed content/language.dark.msg
	languageDark string

	//go:embed content/language.light.msg
	languageLight string

	//go:embed content/ui.dark.msg
	uiDark string

	//go:embed content/ui.light.msg
	uiLight string

	//go:embed content/commands
	commandsFS embed.FS

	//go:embed content/cookbook
	recipesFS embed.FS
)
