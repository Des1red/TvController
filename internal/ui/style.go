package ui

import "github.com/gdamore/tcell/v2"

type UIStyles struct {
	Title  tcell.Style
	Normal tcell.Style
	Label  tcell.Style
	Select tcell.Style
	Edit   tcell.Style
	Dim    tcell.Style
	Border tcell.Style
}

func defaultStyles() UIStyles {
	base := tcell.StyleDefault

	return UIStyles{
		Title:  base.Foreground(tcell.NewRGBColor(120, 200, 255)).Bold(true),
		Normal: base.Foreground(tcell.NewRGBColor(220, 220, 220)),
		Select: base.Foreground(tcell.NewRGBColor(0, 255, 180)).Bold(true),
		Label:  base.Foreground(tcell.NewRGBColor(180, 180, 250)),
		Edit: base.
			Foreground(tcell.NewRGBColor(190, 120, 255)). // purple
			Background(tcell.NewRGBColor(0, 120, 140)).   // cyan-ish bg
			Bold(true),
		Dim:    base.Foreground(tcell.NewRGBColor(130, 130, 130)),
		Border: base.Foreground(tcell.NewRGBColor(90, 140, 180)),
	}
}

func fieldSection(label string) string {
	switch label {
	case "TV IP", "TV Port", "TV Path", "TV Vendor":
		return "TV"

	case "Local file", "Local IP", "Local dir", "Serve port":
		return "Local"

	case "Subnet", "Deep search", "SSDP discovery":
		return "Scan"

	case "Probe only", "Auto cache", "Use cache", "Select cache index":
		return "Cache"

	case "List cache", "Forget cache":
		return "Cache"

	default:
		return ""
	}
}
