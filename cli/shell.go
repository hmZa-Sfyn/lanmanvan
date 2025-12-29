package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"lanmanvan/core"
)

var CurrentDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		if u, err := user.Current(); err == nil {
			homeDir = u.HomeDir
		} else {
			homeDir = "."
		}
	}
	CurrentDir = homeDir
}

// ExecuteShellCommand executes a shell command
func (cli *CLI) ExecuteShellCommand(input string) {
	input = strings.TrimPrefix(input, "$")
	input = strings.TrimSpace(input)

	if input == "" {
		core.PrintWarning("Empty command")
		return
	}

	var shell string
	var cmd *exec.Cmd

	// Determine which shell to use
	if strings.HasPrefix(input, "bash ") {
		shell = "bash"
		input = strings.TrimPrefix(input, "bash ")
		cmd = exec.Command("bash", "-c", input)
	} else if strings.HasPrefix(input, "zsh ") {
		shell = "zsh"
		input = strings.TrimPrefix(input, "zsh ")
		cmd = exec.Command("zsh", "-c", input)
	} else {
		// Default to zsh
		shell = "zsh"
		cmd = exec.Command("zsh", "-c", input)
	}

	startTime := time.Now()
	fmt.Println()
	core.PrintInfo(fmt.Sprintf("Executing in %s", core.Color("cyan", shell)))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = CurrentDir // Set working directory

	err := cmd.Run()
	duration := time.Since(startTime)

	// Update current directory if cd command
	if strings.HasPrefix(strings.TrimSpace(input), "cd ") {
		parts := strings.Fields(input)
		if len(parts) > 1 {
			newDir := strings.Join(parts[1:], " ")
			var updateDir *exec.Cmd
			if newDir == "-" {
				// Handle cd - (previous directory)
				updateDir = exec.Command(shell, "-c", "pwd")
			} else {
				// Resolve the directory path
				updateDir = exec.Command(shell, "-c", fmt.Sprintf("cd %s && pwd", newDir))
			}
			updateDir.Dir = CurrentDir
			if output, err := updateDir.Output(); err == nil {
				CurrentDir = strings.TrimSpace(string(output))
			}
		}
	}

	fmt.Println()
	if err == nil {
		core.PrintSuccess(fmt.Sprintf("Command completed in %s", duration.String()))
	} else {
		core.PrintError(fmt.Sprintf("Command failed: %v (%s)", err, duration.String()))
	}
	fmt.Println()
}
