package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func renderState(
	screen tcell.Screen,
	styles UIStyles,
	state uiState,
	selectedMode int,
	fields []Field,
	selectedField int,
	editMode bool,
	editBuffer string,
	ctx *uiContext,
	confirmSelected int,
) {
	switch state {
	case stateModeSelect:
		renderModeScreen(screen, styles, selectedMode)

	case stateConfig:
		renderInputScreen(
			screen,
			styles,
			fields,
			selectedField,
			editMode,
			editBuffer,
		)

	case stateConfirm:
		renderConfirmScreen(
			screen,
			styles,
			ctx.working.Mode,
			fields,
			confirmSelected,
		)
	}
}

const (
	uiBoxWidth     = 80
	uiBoxMargin    = 2
	uiBoxMaxHeight = 24

	maxTextWidth = uiBoxWidth - 6 // left padding + safety
)

// MODE SELECTION SCREEN
func renderModeScreen(s tcell.Screen, styles UIStyles, selected int) {
	s.Clear()

	w, h := s.Size()

	boxW := uiBoxWidth
	boxH := h - (uiBoxMargin * 2)
	if boxH > uiBoxMaxHeight {
		boxH = uiBoxMaxHeight
	}

	x := (w - boxW) / 2
	y := (h - boxH) / 2
	if y < uiBoxMargin {
		y = uiBoxMargin
	}

	drawBox(s, styles.Border, x, y, boxW, boxH)

	drawText(s, styles.Title, x+4, y+2, "renderctl v2")
	drawText(s, styles.Normal, x+4, y+4, "Select execution mode")

	for i, m := range modes {
		style := styles.Normal
		prefix := "  "
		if i == selected {
			prefix = "> "
			style = styles.Select
		}
		drawText(s, style, x+4, y+6+i, prefix+m)
	}

	drawText(s, styles.Dim, x+4, y+boxH-3, "↑ ↓ move    Enter select    q quit")

	s.Show()
}

// INPUT HANDLING

func renderInputScreen(
	s tcell.Screen,
	styles UIStyles,
	fields []Field,
	selected int,
	editMode bool,
	editBuffer string,
) {
	s.Clear()
	w, h := s.Size()

	boxW := uiBoxWidth
	boxH := h - (uiBoxMargin * 2)
	if boxH > uiBoxMaxHeight {
		boxH = uiBoxMaxHeight
	}

	y := (h - boxH) / 2
	if y < uiBoxMargin {
		y = uiBoxMargin
	}

	x := (w - boxW) / 2

	drawBox(s, styles.Border, x, y, boxW, boxH)

	drawText(s, styles.Title, x+3, y+1, "renderctl configuration")
	drawText(s, styles.Dim, x+3, y+3, "↑ ↓ navigate   Enter edit/toggle   Esc back   q quit")

	rowY := y + 5
	selectIdx := 0 // counts ONLY selectable items
	lastSection := ""

	// ----- FIELDS -----
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		section := fieldSection(f.Label)

		if section != "" && section != lastSection {
			drawText(s, styles.Dim, x+3, rowY, "-- "+section+" --")
			rowY++
			lastSection = section
		}

		style := styles.Normal
		if selectIdx == selected {
			if editMode {
				style = styles.Edit
			} else {
				style = styles.Select
			}
		}

		drawFieldRow(
			s,
			styles,
			f,
			x+3,
			rowY,
			selectIdx == selected,
			editMode,
			editBuffer,
			style,
		)

		rowY++
		selectIdx++
	}

	// spacer
	rowY++
	// ----- EXECUTE -----
	style := styles.Normal
	if selectIdx == selected {
		style = styles.Select
	}
	drawText(s, style, x+3, rowY, "[ Execute ]")
	rowY++
	selectIdx++

	// ----- BACK -----
	style = styles.Normal
	if selectIdx == selected {
		style = styles.Select
	}
	drawText(s, style, x+3, rowY, "[ Back to mode selection ]")

	s.Show()
}

func renderConfirmScreen(
	s tcell.Screen,
	styles UIStyles,
	mode string,
	fields []Field,
	selected int,
) {
	s.Clear()

	w, h := s.Size()

	boxW := uiBoxWidth
	boxH := h - (uiBoxMargin * 2)
	if boxH > uiBoxMaxHeight {
		boxH = uiBoxMaxHeight
	}

	x := (w - boxW) / 2
	y := (h - boxH) / 2
	if y < uiBoxMargin {
		y = uiBoxMargin
	}

	drawBox(s, styles.Border, x, y, boxW, boxH)

	drawText(s, styles.Title, x+3, y+1, "Confirm execution")
	drawText(s, styles.Dim, x+3, y+3, fmt.Sprintf("Mode: %s", mode))

	rowY := y + 5

	// dump fields (simple version, ordered later)
	printKV := func(label, value string) {
		drawText(s, styles.Label, x+3, rowY, fmt.Sprintf("%-22s", label))
		drawText(s, styles.Dim, x+26, rowY, "→")
		drawText(s, styles.Normal, x+30, rowY, value)
		rowY++
	}

	// --- minimal explicit list (safe & clear) ---
	for _, f := range fields {
		var value string

		switch f.Type {
		case FieldBool:
			if f.Bool != nil && *f.Bool {
				value = "Yes"
			} else {
				value = "No"
			}

		case FieldString:
			if f.String != nil {
				value = *f.String
			}

		case FieldInt:
			if f.Int != nil {
				value = fmt.Sprintf("%d", *f.Int)
			}
		}

		printKV(f.Label, value)
	}

	styleConfirm := styles.Normal
	styleCancel := styles.Normal

	if selected == 0 {
		styleConfirm = styles.Select
	} else {
		styleCancel = styles.Select
	}

	rowY++
	drawText(s, styleConfirm, x+3, rowY, "[ Confirm & Execute ]")
	rowY++
	drawText(s, styleCancel, x+3, rowY, "[ Cancel ]")

	s.Show()
}
