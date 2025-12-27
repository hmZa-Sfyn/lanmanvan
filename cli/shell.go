package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"lanmanvan/core"
)

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

	err := cmd.Run()
	duration := time.Since(startTime)

	fmt.Println()
	if err == nil {
		core.PrintSuccess(fmt.Sprintf("Command completed in %s", duration.String()))
	} else {
		core.PrintError(fmt.Sprintf("Command failed: %v (%s)", err, duration.String()))
	}
	fmt.Println()
}
