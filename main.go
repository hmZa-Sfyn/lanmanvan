package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"lanmanvan/cli"
)

func main() {
	var modulesDir string
	var version bool
	var exec string
	var console bool

	flag.StringVar(&modulesDir, "modules", "./modules", "Path to modules directory")
	flag.BoolVar(&version, "version", false, "Show version")
	flag.StringVar(&exec, "exec", "", "Execute a command without starting the interactive shell")
	flag.BoolVar(&console, "console", false, "Start the interactive console")
	flag.Parse()

	if version {
		fmt.Println("LanManVan v2.0.0 - Advanced Metasploit-like Framework in Go")
		os.Exit(0)
	}

	// Expand home directory if needed
	if modulesDir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not determine home directory: %v\n", err)
			os.Exit(1)
		}
		modulesDir = home
	}

	// Make absolute path
	absPath, err := filepath.Abs(modulesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid modules path: %v\n", err)
		os.Exit(1)
	}

	//exec a command and exit
	if exec != "" {
		cliInstance := cli.NewCLI(absPath)

		// This executes the command and sets running = false so the loop exits after one command
		cliInstance.ExecuteCommandAndExit(exec)

		// Always start the CLI — it will either run interactively or exit immediately
		if err := cliInstance.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// Start() returns only after the loop ends cleanly
		os.Exit(0)
	}

	if console != false {
		// Create and start CLI
		cliInstance := cli.NewCLI(absPath)
		if err := cliInstance.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
