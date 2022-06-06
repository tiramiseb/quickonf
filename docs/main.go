package main

import (
	"os"
)

func main() {
	var makeDoc, makeUiFiles bool

	if len(os.Args) == 1 {
		makeDoc = true
		makeUiFiles = true
	} else {
		for _, arg := range os.Args[1:] {
			switch arg {
			case "doc":
				makeDoc = true
			case "ui":
				makeUiFiles = true
			}
		}
	}

	if makeDoc {
		makeDocCommands()
	}

	if makeUiFiles {
		makeUIFiles()
	}
}
