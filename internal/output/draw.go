package output

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// TODO Maybe move to https://github.com/gdamore/tcell, I don't know if it would be useful

const (
	// Global stuff
	clear      = "\033[H\033[2J"
	reset      = "\033[0m"
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"

	// Text modification
	bold   = "\033[1m" // Situation title
	italic = "\033[3m" // Situation title

	// Background colors
	bgBlack = "\033[40m"  // Situation title
	bgRed   = "\033[41m"  // Group title in report if error
	bgBlue  = "\033[44m"  // Group title
	bgGrey  = "\033[100m" // Situation title

	// Foreground colors
	fgWhite = "\033[37m" // Status line
	fgRed   = "\033[31m" // Instruction name (error)
	fgGreen = "\033[32m" // Instruction name (success)
	fgBlue  = "\033[34m" // Instruction name (running)

	prefixRunning = "»"
	prefixSuccess = "✓"
	prefixError   = "!"
)

var width = 80

func Start(nbGroups int) {
	nbWaitingGroups = nbGroups
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	go func() {
		for {
			<-sig
			redraw()
		}
	}()
	fmt.Print(hideCursor)
	redraw()
}

func End() {
	fmt.Print(statusline())
	fmt.Print(report())
	fmt.Println(showCursor)
}

func redraw() {
	size, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		size = &unix.Winsize{Col: 80, Row: 24}
	}
	cols := int(size.Col)
	if cols != width {
		width = cols
		title.preRedraw()
	}
	fmt.Print(clear)
	fmt.Print(title.toDraw)
	if nbWaitingGroups == 0 && nbRunningGroups == 0 {
		return
	}

	var toDraw []string
	rows := int(size.Row) - 4 // situation = 3, statusline = 1, + 1 empty line
	for _, group := range groups {
		toDraw = append(toDraw, group.toDraw...)
	}
	if len(toDraw) > rows {
		toDraw = toDraw[len(toDraw)-rows:]
	} else {
		for len(toDraw) < rows {
			toDraw = append([]string{""}, toDraw...)
		}
	}
	fmt.Println(strings.Join(toDraw, "\n"))
	fmt.Print(statusline())
}
