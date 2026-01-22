package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var reportLines []string
var reportFilename string

var reportStart = time.Now()

func SetFileName(name string) {
	reportFilename = name
}

// report appends a single line to the in-memory report.
// It must NEVER fail or affect runtime behavior.
func report(msg string) {
	if reportFilename == "" {
		return
	}
	if msg == "" {
		return
	}

	elapsed := time.Since(reportStart).Truncate(100 * time.Millisecond)
	wall := time.Now().Format("15:04:05")

	line := fmt.Sprintf("+%s (%s)  %s", elapsed, wall, msg)
	reportLines = append(reportLines, line)

}

// CreateReport writes the collected report to report.txt.
// It is safe to call multiple times (writes once).
func CreateReport() {
	if len(reportLines) == 0 {
		return
	}
	if reportFilename == "" {
		return
	}
	f, err := os.Create(reportFilename)
	if err != nil {
		// never affect runtime behavior
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "renderctl execution report\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	for _, line := range reportLines {
		fmt.Fprintln(f, strings.TrimSpace(line))
	}
}
