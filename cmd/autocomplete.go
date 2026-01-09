package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var installCompletionCmd = &cobra.Command{
	Use:   "install-completion",
	Short: "Install shell auto-completion for tvctrl",
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := filepath.Base(os.Getenv("SHELL"))

		switch shell {
		case "bash":
			return run(`tvctrl completion bash > ~/.tvctrl.bash && echo "source ~/.tvctrl.bash" >> ~/.bashrc`)
		case "zsh":
			return run(`mkdir -p ~/.zsh/completions && tvctrl completion zsh > ~/.zsh/completions/_tvctrl && echo "fpath+=(~/.zsh/completions)" >> ~/.zshrc`)
		case "fish":
			return run(`tvctrl completion fish > ~/.config/fish/completions/tvctrl.fish`)
		default:
			return fmt.Errorf("unsupported shell: %s", shell)
		}
	},
}

func run(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
