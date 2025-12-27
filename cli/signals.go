package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// ModuleExecutor tracks the currently running module process
type ModuleExecutor struct {
	running bool
	pid     int
}

// moduleExecutor is a global instance tracking module execution
var moduleExecutor = &ModuleExecutor{running: false, pid: 0}

// startModuleExecution marks the start of module execution
func (cli *CLI) startModuleExecution() {
	moduleExecutor.running = true
}

// stopModuleExecution marks the end of module execution
func (cli *CLI) stopModuleExecution() {
	moduleExecutor.running = false
	moduleExecutor.pid = 0
}

// setupSignalHandler sets up Ctrl+C handling
func (cli *CLI) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for range sigChan {
			if moduleExecutor.running {
				// Module is running - just mark it as interrupted
				// The module process will handle its own cleanup
				fmt.Println()
				fmt.Println()
				// Return to prompt without exiting the CLI
				continue
			} else {
				// CLI is idle - ask user if they want to exit
				fmt.Println()
				fmt.Println()
				// Just continue, readline will handle the interrupt
				continue
			}
		}
	}()
}
