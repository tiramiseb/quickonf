package help

import _ "embed"

//go:generate go run makemd/main.go

var (
	//go:embed intro.dark.msg
	introDark string

	//go:embed intro.light.msg
	introLight string

	//go:embed language.dark.msg
	languageDark string

	//go:embed language.light.msg
	languageLight string

	//go:embed ui.dark.msg
	uiDark string

	//go:embed ui.light.msg
	uiLight string
)
