package requirements

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var Install bool
var DryRun bool

func printInstallerDescription() {
	fmt.Println("Renderctl Installer")
	fmt.Println("-------------------")
	fmt.Println("This installer builds the renderctl binary and can optionally")
	fmt.Println("install external tools required for streaming mode.")
	fmt.Println()
	fmt.Println("What this does:")
	fmt.Println(" - Builds the renderctl binary using Go")
	fmt.Println(" - Optionally installs streaming dependencies (ffmpeg, yt-dlp)")
	fmt.Println()
	fmt.Println("What this does NOT do:")
	fmt.Println(" - Modify system configuration")
	fmt.Println(" - Change runtime behavior")
	fmt.Println(" - Enable streaming without your consent")
	fmt.Println()
	fmt.Println("Note:")
	fmt.Println(" - Streaming dependencies are optional")
	fmt.Println(" - All other modes work without them")
}

func RunInstaller() error {
	printInstallerDescription()
	fmt.Println("===================")
	if DryRun {
		fmt.Println("\nDry-run mode enabled.")
		fmt.Println("No commands will be executed.")
	}
	osType := runtime.GOOS
	switch osType {
	case "linux", "darwin":
		fmt.Printf("Detected OS: %s\n", osType)
	default:
		return fmt.Errorf("unsupported OS: %s", osType)
	}

	// ---- STEP 2: build ----
	fmt.Println("\nBinary setup")

	if binaryExists() {
		fmt.Println("renderctl binary already exists.")
		if DryRun {
			fmt.Println("Would ask: rebuild binary?")
		} else if confirm("Rebuild renderctl binary?") {
			if !commandExists("go") {
				return fmt.Errorf("go is required to rebuild renderctl")
			}
			if err := buildBinary(); err != nil {
				return err
			}
			fmt.Println("✔ Binary rebuilt")
		}
	} else {
		fmt.Println("renderctl binary not found.")
		if DryRun {
			fmt.Println("Would require: go")
			fmt.Println("Would run: go build -o renderctl main.go")
		} else if confirm("Build renderctl binary now?") {
			if !commandExists("go") {
				return fmt.Errorf("go is required to rebuild renderctl")
			}
			if err := buildBinary(); err != nil {
				return err
			}
			fmt.Println("✔ Build complete")
		}
	}

	// ---- STEP 3: streaming deps explanation ----
	printStreamingNotice()

	// ---- STEP 4: consent ----
	if DryRun {
		fmt.Println("\nSkipping confirmation prompt to install streaming dependencies(dry-run).")
	} else if !confirm("Install streaming dependencies now?") {
		printStreamingDisabled()
		return nil
	}

	// ---- STEP 5: install deps ----
	if err := installStreamingDeps(osType); err != nil {
		return err
	}

	fmt.Println("✔ Streaming dependencies installed")
	printStreamingEnabled()

	return nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func buildBinary() error {
	cmd := exec.Command("go", "build", "-o", "renderctl", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func binaryExists() bool {
	_, err := os.Stat("renderctl")
	return err == nil
}

func printStreamingNotice() {
	fmt.Println("\nStreaming mode requirements")
	fmt.Println("---------------------------")
	fmt.Println("Streaming mode requires external tools:")
	fmt.Println(" - yt-dlp")
	fmt.Println(" - ffmpeg")
	fmt.Println()
	fmt.Println("These are NOT required for:")
	fmt.Println(" - scan")
	fmt.Println(" - probe")
	fmt.Println(" - cache")
	fmt.Println(" - control")
	fmt.Println()
	fmt.Println("If missing:")
	fmt.Println(" - streaming mode will be unavailable")
}

func printStreamingDisabled() {
	fmt.Println("\nStreaming mode: DISABLED")
	fmt.Println("You can install dependencies later to enable it.")
}

func printStreamingEnabled() {
	fmt.Println("\nStreaming mode: ENABLED")
}

func confirm(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n%s [y/N]: ", msg)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

func installStreamingDeps(osType string) error {
	switch osType {
	case "linux":
		return installLinuxDeps()
	case "darwin":
		return installMacDeps()
	default:
		return fmt.Errorf("unsupported OS for dependency install")
	}
}

func installLinuxDeps() error {
	if DryRun {
		switch {
		case commandExists("apt"):
			fmt.Println("Would run: sudo apt install -y ffmpeg yt-dlp")
		case commandExists("dnf"):
			fmt.Println("Would run: sudo dnf install -y ffmpeg yt-dlp")
		case commandExists("pacman"):
			fmt.Println("Would run: sudo pacman -S --noconfirm ffmpeg yt-dlp")
		default:
			fmt.Println("Would install: ffmpeg yt-dlp (unknown package manager)")
		}
		return nil
	}

	switch {
	case commandExists("apt"):
		return runCmd("sudo", "apt", "install", "-y", "ffmpeg", "yt-dlp")
	case commandExists("dnf"):
		return runCmd("sudo", "dnf", "install", "-y", "ffmpeg", "yt-dlp")
	case commandExists("pacman"):
		return runCmd("sudo", "pacman", "-S", "--noconfirm", "ffmpeg", "yt-dlp")
	default:
		return fmt.Errorf("unsupported Linux package manager")
	}
}

func installMacDeps() error {
	if !commandExists("brew") {
		return fmt.Errorf("homebrew not found (brew required)")
	}

	if DryRun {
		fmt.Println("Would run: brew install ffmpeg yt-dlp")
		return nil
	}

	return runCmd("brew", "install", "ffmpeg", "yt-dlp")
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
