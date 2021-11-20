package output

import (
	"strings"
)

type titleData struct {
	text string

	toDraw string
}

var title = &titleData{}

func SetTitle(text string) {
	title = &titleData{text: text}
	title.preRedraw()
	redraw()
}

func (s *titleData) preRedraw() {
	b := strings.Builder{}
	b.WriteString(bgGrey)
	b.WriteString(bold)

	allSpaces := strings.Repeat(" ", width)
	title := "Quickonf: " + s.text
	titleSpaces := width - len(title)

	var before, after string
	if titleSpaces > 0 {
		halfSpaces := titleSpaces / 2
		before = strings.Repeat(" ", halfSpaces)
		after = strings.Repeat(" ", halfSpaces)
		if titleSpaces%2 == 1 {
			before = before + " "
		}
	} else if titleSpaces < 0 {
		title = title[:width]
	}

	b.WriteString(allSpaces)
	b.WriteString(before)
	b.WriteString(title)
	b.WriteString(after)
	b.WriteString(allSpaces)
	b.WriteRune('\n')
	b.WriteString(reset)

	s.toDraw = b.String()

}
