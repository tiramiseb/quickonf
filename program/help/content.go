package help

import (
	"embed"
)

//go:generate go run ../../docs ui

var (
	//go:embed intro.dark.msg
	introDark []byte

	//go:embed intro.light.msg
	introLight []byte

	//go:embed language.dark.msg
	languageDark []byte

	//go:embed language.light.msg
	languageLight []byte

	//go:embed ui.dark.msg
	uiDark []byte

	//go:embed ui.light.msg
	uiLight []byte

	//go:embed commands
	commandsFS embed.FS
)
