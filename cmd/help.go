package cmd

import (
	"fmt"
)

type helpFlag struct {
	name string
	arg  string
	desc string
}

func calcWidths(flags []helpFlag) (flagW, argW int) {
	for _, f := range flags {
		if len(f.name) > flagW {
			flagW = len(f.name)
		}
		if len(f.arg) > argW {
			argW = len(f.arg)
		}
	}
	return
}

func printFlags(flags []helpFlag) {
	flagW, argW := calcWidths(flags)
	for _, f := range flags {
		fmt.Printf(
			"  %-*s  %-*s  %s\n",
			flagW,
			f.name,
			argW,
			f.arg,
			f.desc,
		)
	}
}

func printHelp() {
	fmt.Println("renderctl - Simple TV controller using AVTransport")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  renderctl [flags]")
	fmt.Println()

	// ─── Execution ───────────────────────────────────────────
	fmt.Println("Execution:")
	printFlags([]helpFlag{
		{"--tui", "", "Start program as TUI"},
		{"--probe-only", "", "Probe AVTransport only"},
		{"--mode", "string", "Execution mode (auto/manual/scan/stream)"},
	})
	fmt.Println()

	// ─── Cache ───────────────────────────────────────────────
	fmt.Println("Cache:")
	printFlags([]helpFlag{
		{"--auto-cache", "", "Skip cache save confirmation"},
		{"--no-cache", "", "Disable cache usage"},
		{"--list-cache", "", "List cached devices"},
		{"--forget-cache", "string", "Forget cache (interactive | IP | all)"},
		{"--select-cache", "int", "Select cached device by index"},
		{"--details-cache", "int", "List cached device with details"},
		{"--show-actions", "", "Show supported actions from cached devices"},
		{"--show-media", "", "Show media information from cached devices"},
		{"--show-media-all", "", "Show all media information from cached devices"},
	})
	fmt.Println()

	// ─── Scan ────────────────────────────────────────────────
	fmt.Println("Scan:")
	printFlags([]helpFlag{
		{"--subnet", "string", "Subnet to scan (e.g. 192.168.1.0/24)"},
		{"--deep-search", "", "Extended endpoint probing"},
		{"--ssdp", "", "Enable SSDP discovery"},
		{"--ssdp-timeout", "duration", "SSDP discovery timeout duration"},
	})
	fmt.Println()

	// ─── TV ──────────────────────────────────────────────────
	fmt.Println("TV:")
	printFlags([]helpFlag{
		{"--Tip", "string", "TV IP"},
		{"--Tport", "string", "SOAP port"},
		{"--Tpath", "string", "SOAP path"},
		{"--type", "string", "Vendor"},
	})
	fmt.Println()

	// ─── Media ───────────────────────────────────────────────
	fmt.Println("Media:")
	printFlags([]helpFlag{
		{"--Lf", "string", "Local media file or url (url is stream explicit)"},
		{"--Lip", "string", "Local IP"},
		{"--Ldir", "string", "Local directory"},
		{"--LPort", "string", "Local port"},
	})
	fmt.Println()

	// ─── Output ──────────────────────────────────────────────
	fmt.Println("Output:")
	printFlags([]helpFlag{
		{"--verbose", "", "Enables verbose output"},
		{"--report-file", "string", "Report output file name"},
	})

	fmt.Println()

	// ─── Installer ───────────────────────────────────────────
	fmt.Println("Installer:")
	printFlags([]helpFlag{
		{"--install", "", "Run installer (build binary and optional dependencies)"},
		{"--dry-run", "", "Show installer actions without executing"},
	})
}
