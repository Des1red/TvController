package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func handleKeyEvent(
	ev *tcell.EventKey,
	screen tcell.Screen,
	styles UIStyles,
	ctx *uiContext,
	state *uiState,
	selectedMode *int,
	fields *[]Field,
	selectedField *int,
	editMode *bool,
	editBuffer *string,
	confirmSelected *int,
) {
	// confirm screen is special
	if *state == stateConfirm {
		switch ev.Key() {

		case tcell.KeyUp, tcell.KeyDown:
			*confirmSelected = (*confirmSelected + 1) % 2

		case tcell.KeyEnter:
			if *confirmSelected == 0 {
				ctx.commit()
				*state = stateExit
			} else {
				*selectedMode = 0
				*state = stateModeSelect
			}

		case tcell.KeyEscape:
			*selectedMode = 0
			*state = stateModeSelect
		}
		return
	}

	if *state == stateModeSelect {
		handleModeSelectKey(
			ev, ctx, modes,
			selectedMode, state,
			fields, selectedField,
			editMode, editBuffer,
			screen,
		)
		return
	}

	handleConfigKey(
		ev, ctx, *fields,
		selectedField,
		editMode,
		editBuffer,
		state,
		screen,
		confirmSelected,
	)
}

func handleModeSelectKey(
	ev *tcell.EventKey,
	ctx *uiContext,
	modes []string,
	selectedMode *int,
	state *uiState,
	fields *[]Field,
	selectedField *int,
	editMode *bool,
	editBuffer *string,
	screen tcell.Screen,
) {
	// hard quit
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'q' {
		screen.Fini()
		os.Exit(0)
	}

	switch ev.Key() {
	case tcell.KeyUp:
		*selectedMode = (*selectedMode - 1 + len(modes)) % len(modes)

	case tcell.KeyDown:
		*selectedMode = (*selectedMode + 1) % len(modes)

	case tcell.KeyEnter:
		selected := strings.ToLower(modes[*selectedMode])

		ctx.resetWorking()
		ctx.working.Mode = selected

		*fields = buildFieldsForMode(&ctx.working, ctx.working.Mode)
		*selectedField = 0
		*editMode = false
		*editBuffer = ""
		*state = stateConfig

	case tcell.KeyEscape:
		// do nothing, stay here
	}
}

func handleConfigKey(
	ev *tcell.EventKey,
	ctx *uiContext,
	fields []Field,
	selectedField *int,
	editMode *bool,
	editBuffer *string,
	state *uiState,
	screen tcell.Screen,
	confirmSelected *int,
) {
	// hard quit
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'q' {
		screen.Fini()
		os.Exit(0)
	}

	// ----- EDIT MODE -----
	if *editMode {
		switch ev.Key() {
		case tcell.KeyEnter:
			f := fields[*selectedField]
			if f.Type == FieldString {
				*f.String = *editBuffer
			} else if f.Type == FieldInt {
				if v, err := strconv.Atoi(*editBuffer); err == nil {
					*f.Int = v
				}
			}
			*editMode = false

		case tcell.KeyEscape:
			*editMode = false

		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(*editBuffer) > 0 {
				*editBuffer = (*editBuffer)[:len(*editBuffer)-1]
			}

		case tcell.KeyRune:
			*editBuffer += string(ev.Rune())
		}
		return
	}

	// ----- NORMAL MODE -----
	switch ev.Key() {
	case tcell.KeyUp:
		if *selectedField > 0 {
			*selectedField--
		}

	case tcell.KeyDown:
		if *selectedField < len(fields)+1 {
			*selectedField++
		}

	case tcell.KeyEnter:

		// Execute
		if *selectedField == len(fields) {
			*state = stateConfirm
			*selectedField = 0 // optional safety
			*confirmSelected = 0
			return
		}

		// Back
		if *selectedField == len(fields)+1 {
			*selectedField = 0
			*editMode = false
			*editBuffer = ""
			*state = stateModeSelect
			return
		}

		// Normal field
		f := fields[*selectedField]
		switch f.Type {
		case FieldBool:
			*f.Bool = !*f.Bool
		case FieldString:
			*editMode = true
			*editBuffer = *f.String
		case FieldInt:
			*editMode = true
			*editBuffer = fmt.Sprintf("%d", *f.Int)
		}

	case tcell.KeyEscape:
		*selectedField = 0
		*editMode = false
		*editBuffer = ""
		*state = stateModeSelect
	}
}
