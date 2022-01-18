package help

import "github.com/tiramiseb/quickonf/internal/program/common/style"

func (m *model) updateNav(width int) {
	space := style.BoxContent.Render(" ")
	if width >= 97 {
		m.filterMaxLength = 20
	} else {
		m.filterMaxLength = 10
	}
	switch {
	case width >= 87:
		// " [Introduction] [Configuration language] [Commands] [User interface] Filter:XXXXXXXXXX ": 87 chars
		m.activeIntro = space + style.ClickedButton.Render("[Introduction]") +
			space + style.Button.Render("[Configuration language]") +
			space + style.Button.Render("[Commands]") +
			space + style.Button.Render("[User interface]")
		m.activeLang = space + style.Button.Render("[Introduction]") +
			space + style.ClickedButton.Render("[Configuration language]") +
			space + style.Button.Render("[Commands]") +
			space + style.Button.Render("[User interface]")
		m.activeCommands = space + style.Button.Render("[Introduction]") +
			space + style.Button.Render("[Configuration language]") +
			space + style.ClickedButton.Render("[Commands]") +
			space + style.Button.Render("[User interface]")
		m.activeUI = space + style.Button.Render("[Introduction]") +
			space + style.Button.Render("[Configuration language]") +
			space + style.Button.Render("[Commands]") +
			space + style.ClickedButton.Render("[User interface]")
		m.introStart = 1
		m.introEnd = 14
		m.langStart = 16
		m.langEnd = 39
		m.commandsStart = 41
		m.commandsEnd = 50
		m.uiStart = 52
		m.uiEnd = 67
		m.showFilter = true
		m.buttonsWidth = 68

	case width >= 73:
		// " [Intro] [Config language] [Commands] [User interface] Filter:XXXXXXXXXX ": 73 chars
		m.activeIntro = space + style.ClickedButton.Render("[Intro]") +
			space + style.Button.Render("[Config language]") +
			space + style.Button.Render("[Commands]") +
			space + style.Button.Render("[User interface]")
		m.activeLang = space + style.Button.Render("[Intro]") +
			space + style.ClickedButton.Render("[Config language]") +
			space + style.Button.Render("[Commands]") +
			space + style.Button.Render("[User interface]")
		m.activeCommands = space + style.Button.Render("[Intro]") +
			space + style.Button.Render("[Config language]") +
			space + style.ClickedButton.Render("[Commands]") +
			space + style.Button.Render("[User interface]")
		m.activeUI = space + style.Button.Render("[Intro]") +
			space + style.Button.Render("[Config language]") +
			space + style.Button.Render("[Commands]") +
			space + style.ClickedButton.Render("[User interface]")
		m.introStart = 1
		m.introEnd = 7
		m.langStart = 9
		m.langEnd = 25
		m.commandsStart = 27
		m.commandsEnd = 36
		m.uiStart = 38
		m.uiEnd = 53
		m.showFilter = true
		m.buttonsWidth = 54

	case width >= 62:
		// " [Intro] [Conf lang] [Commands] [Interface] Filter:XXXXXXXXXX ": 62 chars
		m.activeIntro = space + style.ClickedButton.Render("[Intro]") +
			space + style.Button.Render("[Conf lang]") +
			space + style.Button.Render("[Commands]") +
			space + style.Button.Render("[Interface]")
		m.activeLang = space + style.Button.Render("[Intro]") +
			space + style.ClickedButton.Render("[Conf lang]") +
			space + style.Button.Render("[Commands]") +
			space + style.Button.Render("[Interface]")
		m.activeCommands = space + style.Button.Render("[Intro]") +
			space + style.Button.Render("[Conf lang]") +
			space + style.ClickedButton.Render("[Commands]") +
			space + style.Button.Render("[Interface]")
		m.activeUI = space + style.Button.Render("[Intro]") +
			space + style.Button.Render("[Conf lang]") +
			space + style.Button.Render("[Commands]") +
			space + style.ClickedButton.Render("[Interface]")
		m.introStart = 1
		m.introEnd = 7
		m.langStart = 9
		m.langEnd = 19
		m.commandsStart = 21
		m.commandsEnd = 30
		m.uiStart = 32
		m.uiEnd = 42
		m.showFilter = true
		m.buttonsWidth = 43

	case width >= 28:
		// " [Intro] [Lang] [Cmds] [UI] Filter:XXXXXXXXXX ": 46 chars
		// " [Intro] [Lang] [Cmds] [UI] ": 28 chars
		m.activeIntro = space + style.ClickedButton.Render("[Intro]") +
			space + style.Button.Render("[Lang]") +
			space + style.Button.Render("[Cmds]") +
			space + style.Button.Render("[UI]")
		m.activeLang = space + style.Button.Render("[Intro]") +
			space + style.ClickedButton.Render("[Lang]") +
			space + style.Button.Render("[Cmds]") +
			space + style.Button.Render("[UI]")
		m.activeCommands = space + style.Button.Render("[Intro]") +
			space + style.Button.Render("[Lang]") +
			space + style.ClickedButton.Render("[Cmds]") +
			space + style.Button.Render("[UI]")
		m.activeUI = space + style.Button.Render("[Intro]") +
			space + style.Button.Render("[Lang]") +
			space + style.Button.Render("[Cmds]") +
			space + style.ClickedButton.Render("[UI]")
		m.introStart = 1
		m.introEnd = 7
		m.langStart = 9
		m.langEnd = 14
		m.commandsStart = 16
		m.commandsEnd = 21
		m.uiStart = 23
		m.uiEnd = 26
		m.showFilter = width >= 46
		m.buttonsWidth = 27

	case width >= 17:
		// " [I] [L] [C] [U] ": 17 chars
		m.activeIntro = space + style.ClickedButton.Render("[I]") +
			space + style.Button.Render("[L]") +
			space + style.Button.Render("[C]") +
			space + style.Button.Render("[U]")
		m.activeLang = space + style.Button.Render("[I]") +
			space + style.ClickedButton.Render("[L]") +
			space + style.Button.Render("[C]") +
			space + style.Button.Render("[U]")
		m.activeCommands = space + style.Button.Render("[I]") +
			space + style.Button.Render("[L]") +
			space + style.ClickedButton.Render("[C]") +
			space + style.Button.Render("[U]")
		m.activeUI = space + style.Button.Render("[I]") +
			space + style.Button.Render("[L]") +
			space + style.Button.Render("[C]") +
			space + style.ClickedButton.Render("[U]")
		m.introStart = 1
		m.introEnd = 3
		m.langStart = 5
		m.langEnd = 7
		m.commandsStart = 9
		m.commandsEnd = 11
		m.uiStart = 13
		m.uiEnd = 15
		m.showFilter = false
		m.buttonsWidth = 16
	case width >= 5:
		txt := style.BoxContent.Render(" ...")
		m.activeIntro = txt
		m.activeLang = txt
		m.activeCommands = txt
		m.activeUI = txt
		m.introStart = -1
		m.introEnd = -1
		m.langStart = -1
		m.langEnd = -1
		m.commandsStart = -1
		m.commandsEnd = -1
		m.uiStart = -1
		m.uiEnd = -1
		m.showFilter = false
		m.buttonsWidth = 4
	case width >= 3:
		txt := style.BoxContent.Render("...")
		m.activeIntro = txt
		m.activeLang = txt
		m.activeCommands = txt
		m.activeUI = txt
		m.introStart = -1
		m.introEnd = -1
		m.langStart = -1
		m.langEnd = -1
		m.commandsStart = -1
		m.commandsEnd = -1
		m.uiStart = -1
		m.uiEnd = -1
		m.showFilter = false
		m.buttonsWidth = 3
	default:
		m.activeIntro = ""
		m.activeLang = ""
		m.activeCommands = ""
		m.activeUI = ""
		m.introStart = -1
		m.introEnd = -1
		m.langStart = -1
		m.langEnd = -1
		m.commandsStart = -1
		m.commandsEnd = -1
		m.uiStart = -1
		m.uiEnd = -1
		m.showFilter = false
		m.buttonsWidth = 0

	}
}
