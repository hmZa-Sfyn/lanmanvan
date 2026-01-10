package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"lanmanvan/core"
)

type TEMP_Module struct {
	Name string
	Path string
	Type string // "python", "bash", etc
}

// CLI manages the interactive command-line interface
type CLI struct {
	manager     *core.ModuleManager
	running     bool
	history     []string
	envMgr      *EnvironmentManager
	logger      *Logger
	tempModules map[string]TEMP_Module
}

// NewCLI creates a new CLI instance
func NewCLI(modulesDir string) *CLI {
	return &CLI{
		manager: core.NewModuleManager(modulesDir),
		running: true,
		history: make([]string, 0),
		envMgr:  NewEnvironmentManager(),
		logger:  NewLogger(),
	}
}

// Start begins the CLI loop
func (cli *CLI) Start(banner__ bool) error {
	if err := cli.manager.DiscoverModules(); err != nil {
		return err
	}

	if banner__ { // why is this showing when i told it not to? this one too!
		cli.PrintBanner()
	}
	cli.setupSignalHandler()

	// temp modules
	cli.tempModules = make(map[string]TEMP_Module)

	// Create readline instance with history support
	rl, err := cli.getReadlineInstance()
	if err != nil {
		return err
	}
	defer rl.Close()

	for cli.running {
		rl.SetPrompt(cli.GetPrompt())

		input, err := rl.Readline()
		if err != nil {
			if err.Error() == "Interrupt" {
				fmt.Println()
				continue
			}
			if err.Error() == "EOF" {
				break
			}
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cli.history = append(cli.history, input)
		cli.ExecuteCommand(input)
	}

	return nil
}

// Idle start
func (cli *CLI) IdleStart(banner__ bool, command__ string) error {
	if err := cli.manager.DiscoverModules(); err != nil {
		return err
	}

	if banner__ { // why is this showing when i told it not to?
		cli.PrintBanner()
	}
	cli.setupSignalHandler()

	// Create readline instance with history support
	rl, err := cli.getReadlineInstance()
	if err != nil {
		return err
	}
	defer rl.Close()

	for cli.running {
		//rl.SetPrompt(cli.GetPrompt())

		input := command__

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cli.history = append(cli.history, input)
		cli.ExecuteCommand(input)

		break
	}

	return nil
}

// ExecuteCommand processes user commands
func (cli *CLI) ExecuteCommand(input string) {
	// Handle for loops: for VAR in START..END -> COMMAND
	if strings.HasPrefix(input, "for ") && strings.Contains(input, " in ") && strings.Contains(input, " -> ") {
		cli.executeForLoop(input)
		return
	}

	// Handle pipe syntax: cmd1 |> cmd2 |> cmd3
	if strings.Contains(input, "|>") {
		cli.executePipedCommands(input)
		return
	}

	// Handle global environment variable syntax (key=value or key=?)
	if strings.Contains(input, "=") && !strings.Contains(input, " ") {
		parts := strings.SplitN(input, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Check if it's a view operation (key=?)
			if value == "?" {
				if val, exists := cli.envMgr.Get(key); exists {
					fmt.Println()
					fmt.Printf("   %s = %s\n", core.Color("cyan", key), core.Color("green", val))
					fmt.Println()
				} else {
					core.PrintWarning(fmt.Sprintf("Environment variable '%s' not set, skipping...", key))
					fmt.Println()
				}
				return
			}

			// Set environment variable
			if err := cli.envMgr.Set(key, value); err != nil {
				core.PrintError(fmt.Sprintf("Failed to set environment variable: %v, skipping...", err))
				return
			}
			fmt.Println()
			core.PrintSuccess(fmt.Sprintf("Set %s = %s", key, value))
			fmt.Println()
			return
		}
	}

	// Handle shell commands
	if strings.HasPrefix(input, "$") {
		cli.ExecuteShellCommand(input)
		return
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "help", "h", "?":
		cli.PrintHelp()
	case "list", "ls":
		cli.ListModules()
	case "env", "envs":
		cli.envMgr.Display()
	case "search":
		if len(args) > 0 {
			cli.SearchModules(strings.Join(args, " "))
		} else {
			core.PrintError("Usage: search <keyword> ... example: search network")
		}
	case "info":
		if len(args) > 0 {
			cli.ShowModuleInfo(args[0], 1)
		} else {
			core.PrintError("Usage: info <module_name> ... example: info network")
		}
	case "run":
		if len(args) > 0 {
			cli.RunModule(args[0], args[1:])
		} else {
			core.PrintError("Usage: run <module_name> [args...] ... example: run network target_network=$target_network_suffix port=80")
		}
	case "create", "new":
		if len(args) > 0 {
			cli.CreateModule(args[0], args[1:])
		} else {
			core.PrintError("Usage: create <module_name> [python|bash] ... example: create mymodule python")
		}
	case "edit":
		if len(args) > 0 {
			cli.EditModule(args[0])
		} else {
			core.PrintError("Usage: edit <module_name> ... example: edit mymodule")
		}
	case "delete", "remove", "rm":
		if len(args) > 0 {
			cli.DeleteModule(args[0])
		} else {
			core.PrintError("Usage: delete <module_name> ... example: delete mymodule")
		}
	case "history":
		cli.PrintHistory()
	case "clear", "cls":
		cli.ClearScreen()
	case "refresh", "reload":
		cli.RefreshModules()
	case "import", "include":
		if len(args) == 1 {
			cli.ImportModules(args[0])
		} else {
			core.PrintError("Usage: import /path/to/modules")
		}
	case "exit", "quit", "q":
		cli.running = false
		fmt.Println()
		core.PrintSuccess("Goodbye! See you next time.")
		fmt.Println()
	default:
		// Check if command ends with ! (show module info)
		if strings.HasSuffix(cmd, "!") {
			moduleName := strings.TrimSuffix(cmd, "!")
			cli.ShowModuleInfo(moduleName, 0)
		} else {
			// Try to run as a module if command is not recognized
			cli.RunModule(cmd, args)
		}
	}
}

func (cli *CLI) ImportModules(dir string) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		core.PrintError("Invalid module directory")
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		core.PrintError("Failed to read module directory")
		return
	}

	count := 0

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		modulePath := filepath.Join(dir, e.Name())
		moduleType := detectModuleType(modulePath)
		if moduleType == "" {
			continue
		}

		cli.tempModules[e.Name()] = TEMP_Module{
			Name: e.Name(),
			Path: modulePath,
			Type: moduleType,
		}
		count++
	}

	if count == 0 {
		core.PrintWarning("No valid modules found")
		return
	}

	core.PrintSuccess(fmt.Sprintf("Imported %d modules (temporary)", count))
}

func detectModuleType(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, "run.py")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(dir, "run.sh")); err == nil {
		return "bash"
	}
	return ""
}

// GetModuleManager returns the module manager instance
func (cli *CLI) GetModuleManager() *core.ModuleManager {
	return cli.manager
}

// IsRunning returns the running state
func (cli *CLI) IsRunning() bool {
	return cli.running
}

// AddHistory adds a command to history
func (cli *CLI) AddHistory(cmd string) {
	cli.history = append(cli.history, cmd)
}

// GetHistory returns the command history
func (cli *CLI) GetHistory() []string {
	return cli.history
}

// Stop stops the CLI loop
func (cli *CLI) Stop() {
	cli.running = false
}

// RefreshModules refreshes and reloads all modules from the modules directory
func (cli *CLI) RefreshModules() {
	fmt.Println()
	core.PrintInfo("Refreshing modules...")
	fmt.Println()

	// Clear and reinitialize the module manager with the same directory
	modulesDirPath := cli.manager.ModulesDir
	cli.manager = core.NewModuleManager(modulesDirPath)

	// Discover modules again
	if err := cli.manager.DiscoverModules(); err != nil {
		core.PrintError(fmt.Sprintf("Failed to refresh modules: %v", err))
		fmt.Println()
		return
	}

	// Count loaded modules
	modules := cli.manager.ListModules()
	moduleCount := len(modules)

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("✓ Modules refreshed successfully! Loaded %d module(s)", moduleCount))
	fmt.Println()

	// Display summary of loaded modules
	if moduleCount > 0 {
		fmt.Println(core.NmapBox("Loaded Modules"))
		for i, module := range modules {
			status := "✓"
			fmt.Printf("   [%d] %s %s\n", i+1, status, core.Color("cyan", module.Name))
		}
		fmt.Println()
	}
}

// executeForLoop handles for loop syntax: for VAR in START..END -> COMMAND
func (cli *CLI) executeForLoop(input string) {
	// Parse: for VAR in START..END -> COMMAND
	forIdx := strings.Index(input, "for ")
	inIdx := strings.Index(input, " in ")
	arrowIdx := strings.Index(input, " -> ")

	if forIdx == -1 || inIdx == -1 || arrowIdx == -1 {
		core.PrintError("Invalid for loop syntax. Use: for VAR in 0..256 -> COMMAND")
		return
	}

	varName := strings.TrimSpace(input[forIdx+4 : inIdx])
	rangeStr := strings.TrimSpace(input[inIdx+4 : arrowIdx])
	command := strings.TrimSpace(input[arrowIdx+4:])

	// Parse range: START..END
	rangeParts := strings.Split(rangeStr, "..")
	if len(rangeParts) != 2 {
		core.PrintError("Invalid range syntax. Use: 0..256")
		return
	}

	startStr := strings.TrimSpace(rangeParts[0])
	endStr := strings.TrimSpace(rangeParts[1])

	start, errStart := strconv.Atoi(startStr)
	end, errEnd := strconv.Atoi(endStr)

	if errStart != nil || errEnd != nil {
		core.PrintError("Range must contain valid integers")
		return
	}

	fmt.Println()
	core.PrintInfo(fmt.Sprintf("Executing loop: for %s in %d..%d", varName, start, end))
	fmt.Println()

	results := []string{}
	for i := start; i <= end; i++ {
		// Substitute variable in command
		expandedCmd := strings.ReplaceAll(command, "$"+varName, fmt.Sprintf("%d", i))

		// Show what we're executing
		fmt.Printf("  [%d/%d] Executing: %s\n", i-start+1, end-start+1, expandedCmd)

		// Execute the command
		if strings.Contains(expandedCmd, "|>") {
			// For pipes, capture output
			result := cli.executePipedCommandsForLoop(expandedCmd)
			results = append(results, result)
		} else {
			// For modules, execute normally
			parts := strings.Fields(expandedCmd)
			if len(parts) > 0 {
				cli.ExecuteCommand(expandedCmd)
			}
		}
	}

	// Display results if any were captured
	if len(results) > 0 {
		fmt.Println()
		core.PrintSuccess("Loop Results:")
		for i, result := range results {
			fmt.Printf("   [%d] %s\n", i, result)
		}
		fmt.Println()
	}
}

// executePipedCommandsForLoop handles pipes and returns output instead of printing
func (cli *CLI) executePipedCommandsForLoop(input string) string {
	parts := strings.Split(input, "|>")
	if len(parts) < 2 {
		return ""
	}

	var result string
	var err error

	// Execute first command
	firstCmd := strings.TrimSpace(parts[0])
	result, err = cli.executePipedCommand(firstCmd, "")
	if err != nil {
		return ""
	}

	// Execute remaining commands, passing output as input
	for i := 1; i < len(parts); i++ {
		nextCmd := strings.TrimSpace(parts[i])
		result, err = cli.executePipedCommand(nextCmd, result)
		if err != nil {
			return ""
		}
	}

	return result
}

// executePipedCommands handles piped commands with |> syntax
// Example: whoami() |> sha256() or cat(file.txt) |> base64()
func (cli *CLI) executePipedCommands(input string) {
	parts := strings.Split(input, "|>")
	if len(parts) < 2 {
		return
	}

	var result string
	var err error

	// Execute first command
	firstCmd := strings.TrimSpace(parts[0])
	result, err = cli.executePipedCommand(firstCmd, "")
	if err != nil {
		core.PrintError(fmt.Sprintf("Pipe error in first command: %v", err))
		return
	}

	// Execute remaining commands, passing output as input
	for i := 1; i < len(parts); i++ {
		nextCmd := strings.TrimSpace(parts[i])
		result, err = cli.executePipedCommand(nextCmd, result)
		if err != nil {
			core.PrintError(fmt.Sprintf("Pipe error at step %d: %v", i+1, err))
			return
		}
	}

	// Only print result if the last command is not file() - file() handles its own output
	lastCmd := strings.TrimSpace(parts[len(parts)-1])
	if !strings.HasPrefix(lastCmd, "file(") {
		fmt.Println()
		fmt.Println(result)
		fmt.Println()
	}
}

// executePipedCommand executes a single command in a pipe chain
// Supports: builtin(args), module, module arg=value
func (cli *CLI) executePipedCommand(cmd string, input string) (string, error) {
	cmd = strings.TrimSpace(cmd)

	// Handle string literals in pipes: "\n", "\t", "text", etc.
	if (strings.HasPrefix(cmd, "\"") && strings.HasSuffix(cmd, "\"")) ||
		(strings.HasPrefix(cmd, "'") && strings.HasSuffix(cmd, "'")) {
		// Remove quotes and process escape sequences
		literal := cmd[1 : len(cmd)-1]

		// Process escape sequences
		literal = strings.ReplaceAll(literal, "\\n", "\n")
		literal = strings.ReplaceAll(literal, "\\t", "\t")
		literal = strings.ReplaceAll(literal, "\\r", "\r")
		literal = strings.ReplaceAll(literal, "\\\\", "\\")

		// String literals just pass through, potentially appending to input
		return input + literal, nil
	}

	// If input from previous command, inject it appropriately
	if input != "" {
		// If command is a builtin function call
		if strings.Contains(cmd, "(") && strings.Contains(cmd, ")") {
			openParen := strings.Index(cmd, "(")
			closeParen := strings.LastIndex(cmd, ")")
			if openParen > 0 && closeParen > openParen {
				funcName := cmd[:openParen]
				args := cmd[openParen+1 : closeParen]
				if args != "" {
					args += ", \"" + input + "\""
				} else {
					args = "\"" + input + "\""
				}
				cmd = funcName + "(" + args + ")"
			}
		} else {
			// It's a module call with potential arguments
			// Check if there's an argument pattern like: modulename ip=$somevar
			if strings.Contains(cmd, "=") {
				// Module with specific arguments - find what argument to inject into
				// If pattern is "module arg=$var", replace $var with input
				if strings.Contains(cmd, "$") {
					// Find the variable and replace it
					parts := strings.Split(cmd, "=")
					if len(parts) >= 2 {
						// Replace the variable value with piped input
						lastPart := parts[len(parts)-1]
						if strings.HasPrefix(lastPart, "$") {
							// Replace the variable
							varName := strings.TrimSpace(lastPart)
							cmd = strings.Replace(cmd, varName, "\""+input+"\"", 1)
						} else {
							// Append input as new argument
							cmd = cmd + " input=\"" + input + "\""
						}
					}
				} else {
					// Append input as new argument
					cmd = cmd + " input=\"" + input + "\""
				}
			} else {
				// No arguments, just module name
				cmd = cmd + " input=\"" + input + "\""
			}
		}
	}

	// Try to execute as module
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		moduleName := parts[0]
		args := parts[1:]

		// Check if module exists
		if _, err := cli.manager.GetModule(moduleName); err == nil {
			// Execute module and capture output
			return cli.executeModuleForPipe(moduleName, args)
		}
	}

	return "", fmt.Errorf("invalid pipe command: %s", cmd)
}

// executeModuleForPipe executes a module and returns its output
func (cli *CLI) executeModuleForPipe(moduleName string, args []string) (string, error) {
	_, err := cli.manager.GetModule(moduleName)
	if err != nil {
		return "", err
	}

	// Parse arguments with support for variable expansion
	moduleArgs := make(map[string]string)
	parsedArgs := cli.parseArguments(args)

	for key, value := range parsedArgs {
		switch key {
		case "threads", "save":
			// Skip these
		default:
			moduleArgs[key] = value
		}
	}

	// Merge global environment variables
	for key, value := range cli.envMgr.GetAll() {
		if _, exists := moduleArgs[key]; !exists {
			moduleArgs[key] = value
		}
	}

	// Save original stdout to restore later
	saveOut := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		// Fallback: execute without capturing stdout
		result, execErr := cli.manager.ExecuteModule(moduleName, moduleArgs)
		if execErr != nil {
			return "", execErr
		}
		return strings.TrimSpace(result.Output), nil
	}

	// Redirect stdout to our pipe
	os.Stdout = writer

	// Execute module
	result, err := cli.manager.ExecuteModule(moduleName, moduleArgs)

	// Restore stdout
	writer.Close()
	os.Stdout = saveOut

	// Read captured stdout (not used but needed to drain the pipe)
	_ = reader
	reader.Close()

	if err != nil {
		return "", err
	}

	// Return the module output directly, which contains the structured output
	return strings.TrimSpace(result.Output), nil
}

// parseAdvancedArguments parses function arguments with support for:
// - Quoted strings (both "..." and '...')
// - Nested builtins $(builtin args) and builtin() function call syntax
// - Variable expansion $var
// - Space-separated arguments
func (cli *CLI) parseAdvancedArguments(argsStr string) []string {
	var args []string
	var currentArg strings.Builder
	i := 0

	for i < len(argsStr) {
		ch := argsStr[i]

		// Handle quoted strings
		if ch == '"' || ch == '\'' {
			quote := ch
			i++ // skip opening quote
			for i < len(argsStr) && argsStr[i] != quote {
				if argsStr[i] == '\\' && i+1 < len(argsStr) {
					// Handle escape sequences
					i++
					currentArg.WriteByte(argsStr[i])
				} else {
					currentArg.WriteByte(argsStr[i])
				}
				i++
			}
			i++ // skip closing quote
			continue
		}

		// Handle variable expansion: $varname
		if ch == '$' && i+1 < len(argsStr) && isValidVarChar(rune(argsStr[i+1])) {
			i++ // skip $
			var varName strings.Builder
			for i < len(argsStr) && isValidVarChar(rune(argsStr[i])) {
				varName.WriteByte(argsStr[i])
				i++
			}
			varVal := cli.expandVariable(varName.String())
			currentArg.WriteString(varVal)
			continue
		}

		// Handle comma-separated arguments
		if ch == ',' {
			arg := strings.TrimSpace(currentArg.String())
			if arg != "" {
				args = append(args, arg)
			}
			currentArg.Reset()
			i++
			continue
		}

		// Handle spaces (space-separated arguments)
		if ch == ' ' {
			arg := strings.TrimSpace(currentArg.String())
			if arg != "" {
				args = append(args, arg)
			}
			currentArg.Reset()
			i++
			continue
		}

		currentArg.WriteByte(ch)
		i++
	}

	// Add final argument
	arg := strings.TrimSpace(currentArg.String())
	if arg != "" {
		args = append(args, arg)
	}

	return args
}

// isValidVarChar checks if a rune is valid in a variable name
func isValidVarChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

// collectIdentifier extracts an identifier starting at position i
func (cli *CLI) collectIdentifier(s string, i *int) string {
	var ident strings.Builder
	for *i < len(s) && isValidVarChar(rune(s[*i])) {
		ident.WriteByte(s[*i])
		*i++
	}
	return ident.String()
}

// findMatchingParen finds the index of the closing parenthesis that matches
// the opening parenthesis at startIdx
func (cli *CLI) findMatchingParen(s string, startIdx int) int {
	depth := 1
	i := startIdx
	inQuote := false
	quoteChar := byte(0)

	for i < len(s) && depth > 0 {
		ch := s[i]

		// Handle quotes
		if (ch == '"' || ch == '\'') && (i == 0 || s[i-1] != '\\') {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
			}
		}

		// Handle parentheses (only outside quotes)
		if !inQuote {
			if ch == '(' {
				depth++
			} else if ch == ')' {
				depth--
				if depth == 0 {
					return i
				}
			}
		}

		i++
	}

	return -1 // Not found
}

// expandVariable expands a variable reference
func (cli *CLI) expandVariable(varName string) string {
	if val, exists := cli.envMgr.Get(varName); exists {
		return val
	}
	if val, exists := os.LookupEnv(varName); exists {
		return val
	}
	return "$" + varName
}
