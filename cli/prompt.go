package cli

import (
	"fmt"
	"os"
	"os/user"

	"github.com/fatih/color"
)

// GetPrompt returns the CLI prompt
func (cli *CLI) GetPrompt() string {
	user, _ := user.Current()
	hostname, _ := os.Hostname()

	// Simple, colorful prompt
	return fmt.Sprintf("%s%s%s%s ",
		color.CyanString(user.Username),
		color.WhiteString("@"),
		color.MagentaString(hostname),
		color.GreenString(" ❯"),
	)
}

// PrintBanner prints the application banner
func (cli *CLI) PrintBanner() {
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println(`
  ██▓     ▄▄▄       ███▄    █  ███▄ ▄███▓ ▄▄▄       ███▄    █  ██▒   █▓ ▄▄▄       ███▄    █ 
 ▓██▒    ▒████▄     ██ ▀█   █  ▓██▒▀█▀ ██▒▒████▄     ██ ▀█   █ ▓██░   █▒▒████▄     ██ ▀█   █ 
 ▒██░    ▒██  ▀█▄  ▓██  ▀█ ▄█▒ ▓██    ▓██░▒██  ▀█▄  ▓██  ▀█ ▄█  ▓██  █▒░▒██  ▀█▄  ▓██  ▀█ ▄█ 
 ▒██░    ░██▄▄▄▄██ ▓██▒  ▐▌██▒ ▒██    ▒██ ░██▄▄▄▄██ ▓██▒  ▐▌██▒   ▓██ ░▒░░██▄▄▄▄██ ▓██▒  ▐▌██▒
 ░██████▒ ▓█   ▓██▒▒██░   ▓██░ ▒██▒   ░██▒ ▓█   ▓██▒▒██░   ▓██░   ▒▀█░░░  ▓█   ▓██▒▒██░   ▓██░
 ░ ▒░▓  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒  ░ ▒░   ░  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒    ░ ▐░░░  ▒▒   ▓▒█░░ ▒░   ▒ ▒ 
 ░ ░ ▒  ░  ▒   ▒▒ ░░ ░░   ░ ▒░░  ░      ░  ▒   ▒▒ ░░ ░░   ░ ▒░   ░ ░░░░   ▒   ▒▒ ░░ ░░   ░ ▒░
   ░ ░     ░   ▒      ░   ░ ░ ░      ░     ░   ▒      ░   ░ ░      ░     ░   ▒      ░   ░ ░ 
     ░  ░      ░  ░        ░        ░         ░  ░        ░        ░           ░  ░        ░
`)

	fmt.Println()
	color.New(color.FgGreen, color.Bold).Println("╔════════════════════════════════════════════════════════════════════╗")
	color.New(color.FgGreen, color.Bold).Println("║   ✦ LANMANVAN v2.0 - Advanced Modular Exploitation Framework ✦      ║")
	color.New(color.FgGreen, color.Bold).Println("║   Go Core | Python3/Bash Modules | Dynamic UI | Security Tools       ║")
	color.New(color.FgGreen, color.Bold).Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Type %s for available commands\n\n", color.CyanString("'help'"))
}

// ClearScreen clears the terminal
func (cli *CLI) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
