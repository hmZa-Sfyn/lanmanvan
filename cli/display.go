package cli

import (
	"fmt"
	"sort"
	"strings"

	"lanmanvan/core"

	"github.com/fatih/color"
)

// PrintHelp prints available commands
func (cli *CLI) PrintHelp() {
	fmt.Println()
	color.New(color.FgCyan, color.Bold).Println("╔════════════════════════════════════════════════════════════════╗")
	color.New(color.FgCyan, color.Bold).Println("║                    📚 AVAILABLE COMMANDS                        ║")
	color.New(color.FgCyan, color.Bold).Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	commands := []struct {
		name string
		desc string
	}{
		{"help, h, ?", "Show this help message"},
		{"list, ls", "List all modules"},
		{"search <keyword>", "Search modules by name/tag"},
		{"info <module>", "Show detailed module information"},
		{"<module>!", "Quick show module options and usage"},
		{"run <module> [args]", "Execute a module with arguments"},
		{"<module> [args]", "Shorthand: <module> arg_key=value"},
		{"<module> arg_key = value", "Format with spaces (alternative)"},
		{"env, envs", "Show all global environment variables"},
		{"key=value", "Set global environment variable (persistent)"},
		{"key=?", "View global environment variable value"},
		{"create <name> [type]", "Create a new module (python/bash)"},
		{"edit <module>", "Edit module files"},
		{"delete <module>", "Delete a module"},
		{"history", "Show command history"},
		{"clear", "Clear screen"},
		{"exit, quit, q", "Exit framework"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %s%-32s %s\n",
			color.GreenString("❯"),
			color.CyanString(cmd.name),
			cmd.desc,
		)
	}
	fmt.Println()
}

// ListModules displays all available modules
func (cli *CLI) ListModules() {
	modules := cli.manager.ListModules()
	if len(modules) == 0 {
		core.PrintWarning("No modules loaded.")
		fmt.Println()
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("AVAILABLE MODULES (%d)", len(modules))))

	// Sort modules by name
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})

	for i, module := range modules {
		typeBadge := cli.getTypeBadge(module.Type)
		desc := ""
		tags := ""

		if module.Metadata != nil {
			if module.Metadata.Description != "" {
				desc = module.Metadata.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}
			}
			if len(module.Metadata.Tags) > 0 {
				tags = strings.Join(module.Metadata.Tags[:1], "")
			}
		}

		prefix := "   \\_ "
		if i == len(modules)-1 {
			prefix = "   \\_ "
		}

		line := fmt.Sprintf("%s%s %s  %s %s",
			prefix,
			color.CyanString(module.Name),
			typeBadge,
			color.WhiteString(desc),
			color.MagentaString(tags),
		)
		fmt.Println(line)
	}

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("Total: %d modules loaded", len(modules)))
	fmt.Println()
}

// SearchModules searches modules by keyword
func (cli *CLI) SearchModules(keyword string) {
	modules := cli.manager.ListModules()
	keyword = strings.ToLower(keyword)

	var results []*core.ModuleConfig

	for _, module := range modules {
		name := strings.ToLower(module.Name)
		if strings.Contains(name, keyword) {
			results = append(results, module)
			continue
		}

		if module.Metadata != nil {
			desc := strings.ToLower(module.Metadata.Description)
			if strings.Contains(desc, keyword) {
				results = append(results, module)
				continue
			}

			for _, tag := range module.Metadata.Tags {
				if strings.Contains(strings.ToLower(tag), keyword) {
					results = append(results, module)
					break
				}
			}
		}
	}

	if len(results) == 0 {
		core.PrintWarning(fmt.Sprintf("No modules found for '%s'", keyword))
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("SEARCH: %s (%d results)", keyword, len(results))))

	for i, module := range results {
		typeBadge := cli.getTypeBadge(module.Type)
		desc := ""
		if module.Metadata != nil && module.Metadata.Description != "" {
			desc = module.Metadata.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
		}

		prefix := "   \\_ "
		if i == len(results)-1 {
			prefix = "   \\_ "
		}

		line := fmt.Sprintf("%s%s %s  %s",
			prefix,
			color.CyanString(module.Name),
			typeBadge,
			color.WhiteString(desc),
		)
		fmt.Println(line)
	}

	fmt.Println()
	core.PrintSuccess(fmt.Sprintf("Found %d module(s)", len(results)))
	fmt.Println()
}

// ShowModuleInfo displays detailed module information
func (cli *CLI) ShowModuleInfo(moduleName string) {
	module, err := cli.manager.GetModule(moduleName)
	if err != nil {
		core.PrintError(fmt.Sprintf("Error: %v", err))
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("MODULE: %s", moduleName)))

	if module.Metadata != nil {
		meta := module.Metadata
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Description:"), color.WhiteString(meta.Description))
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Type:"), cli.getTypeBadge(meta.Type))
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Author:"), color.YellowString(meta.Author))
		fmt.Printf("   ├─ %s %s\n", color.WhiteString("Version:"), color.MagentaString(meta.Version))

		if len(meta.Tags) > 0 {
			fmt.Printf("   ├─ %s %s\n", color.WhiteString("Tags:"), color.CyanString(strings.Join(meta.Tags, ", ")))
		}

		if len(meta.Options) > 0 {
			fmt.Printf("   └─ %s\n", color.WhiteString("Options:"))

			for optName, opt := range meta.Options {
				required := ""
				if opt.Required {
					required = color.RedString(" [REQUIRED]")
				}

				fmt.Printf("       ├─ %s %s%s\n",
					color.GreenString(optName),
					color.WhiteString(fmt.Sprintf("(%s)", opt.Type)),
					required,
				)
				fmt.Printf("       │  └─ %s\n", color.WhiteString(opt.Description))
			}
		}
	} else {
		fmt.Printf("   └─ %s\n", color.WhiteString("Type: "+cli.getTypeBadge(module.Type)))
		fmt.Println("   (No metadata available)")
	}

	fmt.Println()
}

// PrintHistory shows command history
func (cli *CLI) PrintHistory() {
	if len(cli.history) == 0 {
		core.PrintWarning("No command history")
		return
	}

	fmt.Println()
	fmt.Println(core.NmapBox(fmt.Sprintf("COMMAND HISTORY (%d)", len(cli.history))))

	for i, cmd := range cli.history {
		fmt.Printf("   \\_ %s %s\n",
			color.GreenString(fmt.Sprintf("[%d]", i+1)),
			color.WhiteString(cmd),
		)
	}

	fmt.Println()
}

// getTypeBadge returns a colored badge for module type
func (cli *CLI) getTypeBadge(moduleType string) string {
	switch moduleType {
	case "python":
		return color.BlueString("[PY]")
	case "bash":
		return color.CyanString("[SH]")
	case "go":
		return color.MagentaString("[GO]")
	default:
		return color.WhiteString("[??]")
	}
}
