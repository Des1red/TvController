package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func drawBox(s tcell.Screen, style tcell.Style, x, y, w, h int) {
	s.SetContent(x, y, '╔', nil, style)
	s.SetContent(x+w-1, y, '╗', nil, style)
	s.SetContent(x, y+h-1, '╚', nil, style)
	s.SetContent(x+w-1, y+h-1, '╝', nil, style)

	for i := 1; i < w-1; i++ {
		s.SetContent(x+i, y, '═', nil, style)
		s.SetContent(x+i, y+h-1, '═', nil, style)
	}

	for i := 1; i < h-1; i++ {
		s.SetContent(x, y+i, '║', nil, style)
		//s.SetContent(x+w-1, y+i, '║', nil, style)
	}
}

func drawText(s tcell.Screen, style tcell.Style, x, y int, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func clipText(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}

func drawFieldRow(
	s tcell.Screen,
	styles UIStyles,
	f Field,
	x int,
	y int,
	isSelected bool,
	editMode bool,
	editBuffer string,
	style tcell.Style,
) {
	// label
	labelText := fmt.Sprintf("%-22s", f.Label)

	// value
	value := ""
	switch f.Type {
	case FieldBool:
		if *f.Bool {
			value = "Yes"
		} else {
			value = "No"
		}
	case FieldString:
		value = *f.String
	case FieldInt:
		value = fmt.Sprintf("%d", *f.Int)
	}

	if editMode && isSelected {
		value = editBuffer
	}

	// draw label (dark red)
	drawText(s, styles.Label, x, y, labelText)

	// draw value
	drawText(
		s,
		style,
		x+len(labelText),
		y,
		" : "+clipText(value, maxTextWidth-len(labelText)-3),
	)
}
