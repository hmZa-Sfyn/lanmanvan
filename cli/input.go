package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
)

// getHistoryPath returns the persistent history file path
func getHistoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/lanmanvan_history"
	}
	configDir := filepath.Join(homeDir, ".lanmanvan")
	os.MkdirAll(configDir, 0700)
	return filepath.Join(configDir, "history")
}

// getReadlineInstance creates a readline instance with history support and copy-paste enabled
func (cli *CLI) getReadlineInstance() (*readline.Instance, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:         "",
		HistoryFile:    getHistoryPath(),
		FuncIsTerminal: func() bool { return true }, // Treat as terminal for paste support
	})

	if err != nil {
		return nil, err
	}

	return rl, nil
}

// readWithHistory reads input with history navigation support
func (cli *CLI) readWithHistory() (string, error) {
	rl, err := cli.getReadlineInstance()
	if err != nil {
		return "", err
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err.Error() == "Interrupt" {
				fmt.Println()
				continue
			}
			return "", err
		}

		return line, nil
	}
}
