package main

import (
	"github.com/vcoder4c/jirakits/cli"
	"os"
)

const (
	Name        = "jirakits"
	Version     = "1.0.0"
	SettingDir  = ".jirakits"
	MappingFile = "mapping.json"
)

func main() {
	_ = os.Mkdir(projectDir(), 0700)
	handlers := []*cli.Handler{
		&cli.Handler{
			Pattern:     "version",
			Description: "Print version",
			Callback:    printVersion,
		},
		&cli.Handler{
			Pattern:     "help",
			Description: "Print help",
			Callback:    printHelp,
		},
		&cli.Handler{
			Pattern:     "help <command>",
			Description: "Print command help",
			Callback:    printCommandHelp,
		},
		&cli.Handler{
			Pattern:     "help <command> <subcommand>",
			Description: "Print subcommand help",
			Callback:    printSubCommandHelp,
		},
		&cli.Handler{
			Pattern:     "jira [options]",
			Description: "Jira Configuration",
			Callback:    configJira,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "check",
						Patterns:    []string{"--check"},
						Description: "Check Jira configuration",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "url",
						Patterns:    []string{"--url"},
						Description: "Jira URL",
					},
					cli.StringFlag{
						Name:        "token",
						Patterns:    []string{"--token"},
						Description: "Jira Token",
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "drive [options]",
			Description: "Trigger to authorize with Google",
			Callback:    configDrive,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "check",
						Patterns:    []string{"--check"},
						Description: "Check Drive configuration",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "setup",
						Patterns:    []string{"--setup"},
						Description: "Setup Drive",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "fetch",
			Description: "Fetch Jira issues with JQL",
			Callback:    fetchIssues,
		},
		&cli.Handler{
			Pattern:     "check",
			Description: "Check configuration",
			Callback:    configCheck,
		},
	}

	cli.SetHandlers(handlers)

	if ok := cli.Handle(os.Args[1:]); !ok {
		exitF("No valid arguments given, use '%s help' to see available commands", Name)
	}
}
