package cli

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

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

// PrintBanner prints a random application banner
func (cli *CLI) PrintBanner() {
	// Seed random (only needed once)
	rand.Seed(time.Now().UnixNano())

	banners := []string{
		// Banner 1: Original LANMANVAN
		`
  ██▓     ▄▄▄       ███▄    █  ███▄ ▄███▓ ▄▄▄       ███▄    █  ██▒   █▓ ▄▄▄       ███▄    █ 
 ▓██▒    ▒████▄     ██ ▀█   █  ▓██▒▀█▀ ██▒▒████▄     ██ ▀█   █ ▓██░   █▒▒████▄     ██ ▀█   █ 
 ▒██░    ▒██  ▀█▄  ▓██  ▀█ ▄█▒ ▓██    ▓██░▒██  ▀█▄  ▓██  ▀█ ▄█  ▓██  █▒░▒██  ▀█▄  ▓██  ▀█ ▄█ 
 ▒██░    ░██▄▄▄▄██ ▓██▒  ▐▌██▒ ▒██    ▒██ ░██▄▄▄▄██ ▓██▒  ▐▌██▒   ▓██ ░▒░░██▄▄▄▄██ ▓██▒  ▐▌██▒
 ░██████▒ ▓█   ▓██▒▒██░   ▓██░ ▒██▒   ░██▒ ▓█   ▓██▒▒██░   ▓██░   ▒▀█░░░  ▓█   ▓██▒▒██░   ▓██░
 ░ ▒░▓  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒  ░ ▒░   ░  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒    ░ ▐░░░  ▒▒   ▓▒█░░ ▒░   ▒ ▒ 
 ░ ░ ▒  ░  ▒   ▒▒ ░░ ░░   ░ ▒░░  ░      ░  ▒   ▒▒ ░░ ░░   ░ ▒░   ░ ░░░░   ▒   ▒▒ ░░ ░░   ░ ▒░
   ░ ░     ░   ▒      ░   ░ ░ ░      ░     ░   ▒      ░   ░ ░      ░     ░   ▒      ░   ░ ░ 
     ░  ░      ░  ░        ░        ░         ░  ░        ░        ░           ░  ░        ░
`,

		// Banner 2: Apache-inspired
		`
                                                                                             
@@@        @@@@@@   @@@  @@@  @@@@@@@@@@    @@@@@@   @@@  @@@  @@@  @@@   @@@@@@   @@@  @@@  
@@@       @@@@@@@@  @@@@ @@@  @@@@@@@@@@@  @@@@@@@@  @@@@ @@@  @@@  @@@  @@@@@@@@  @@@@ @@@  
@@!       @@!  @@@  @@!@!@@@  @@! @@! @@!  @@!  @@@  @@!@!@@@  @@!  @@@  @@!  @@@  @@!@!@@@  
!@!       !@!  @!@  !@!!@!@!  !@! !@! !@!  !@!  @!@  !@!!@!@!  !@!  @!@  !@!  @!@  !@!!@!@!  
@!!       @!@!@!@!  @!@ !!@!  @!! !!@ @!@  @!@!@!@!  @!@ !!@!  @!@  !@!  @!@!@!@!  @!@ !!@!  
!!!       !!!@!!!!  !@!  !!!  !@!   ! !@!  !!!@!!!!  !@!  !!!  !@!  !!!  !!!@!!!!  !@!  !!!  
!!:       !!:  !!!  !!:  !!!  !!:     !!:  !!:  !!!  !!:  !!!  :!:  !!:  !!:  !!!  !!:  !!!  
 :!:      :!:  !:!  :!:  !:!  :!:     :!:  :!:  !:!  :!:  !:!   ::!!:!   :!:  !:!  :!:  !:!  
 :: ::::  ::   :::   ::   ::  :::     ::   ::   :::   ::   ::    ::::    ::   :::   ::   ::  
: :: : :   :   : :  ::    :    :      :     :   : :  ::    :      :       :   : :  ::    :   
`,

		// Banner 5: Matrix-style / Hacker vibe
		`
██╗      █████╗ ███╗   ██╗███╗   ███╗ █████╗ ███╗   ██╗██╗   ██╗ █████╗ ███╗   ██╗
██║     ██╔══██╗████╗  ██║████╗ ████║██╔══██╗████╗  ██║██║   ██║██╔══██╗████╗  ██║
██║     ███████║██╔██╗ ██║██╔████╔██║███████║██╔██╗ ██║██║   ██║███████║██╔██╗ ██║
██║     ██╔══██║██║╚██╗██║██║╚██╔╝██║██╔══██║██║╚██╗██║╚██╗ ██╔╝██╔══██║██║╚██╗██║
███████╗██║  ██║██║ ╚████║██║ ╚═╝ ██║██║  ██║██║ ╚████║ ╚████╔╝ ██║  ██║██║ ╚████║
╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝  ╚═══╝  ╚═╝  ╚═╝╚═╝  ╚═══╝
`,
	}

	// Pick random banner
	banner := banners[2]

	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println(banner)
	fmt.Println()

	// Common footer for all banners
	color.New(color.FgGreen, color.Bold).Println("╔═════════════════════════════════════════════════════════════════╗")
	color.New(color.FgGreen, color.Bold).Println("║   ✦ LANMANVAN v2.0 - Advanced Modular Tooling Framework ✦       ║")
	color.New(color.FgGreen, color.Bold).Println("║   Go Core | Python3/Bash Modules | Dynamic UI | Security Tools  ║")
	color.New(color.FgGreen, color.Bold).Println("╚═════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("Type %s for available commands, have fun!\n\n", color.CyanString("'help'"))
}

// ClearScreen clears the terminal
func (cli *CLI) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
