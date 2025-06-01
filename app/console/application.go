package console

import (
	"context"
	"os"
	"slices"
	"strings"

	"github.com/jeanmolossi/ai-agent-cli/app/contracts/console"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/console/command"
	"github.com/jeanmolossi/ai-agent-cli/app/support/env"
	"github.com/urfave/cli/v3"
)

var (
	noANSI     bool
	noANSIFlag = &cli.BoolFlag{
		Name:        "no-ansi",
		Destination: &noANSI,
		HideDefault: true,
		Usage:       "Force disable ANSI output",
	}
)

type Application struct {
	instance     *cli.Command
	useAiGoAgent bool
}

// NewApplication create a new AiGoAgent application.
// Will add AiGoAgent flag to the command if useAiGoAgent is true.
func NewApplication(name, usage, usageText, version string, useAiGoAgent bool) console.AiGoAgent {
	instance := &cli.Command{}
	instance.Name = name
	instance.Usage = usage
	instance.UsageText = usageText
	instance.Version = version
	instance.CommandNotFound = commandNotFound
	instance.OnUsageError = onUsageError
	instance.Flags = []cli.Flag{noANSIFlag}

	return &Application{
		instance:     instance,
		useAiGoAgent: useAiGoAgent,
	}
}

// Register implements console.AiGoAgent.
func (a *Application) Register(commands []console.Command) {
	for _, item := range commands {
		cliCommand := cli.Command{
			Name:  item.Signature(),
			Usage: item.Description(),
			Action: func(_ context.Context, cmd *cli.Command) error {
				return item.Handle(NewCliContext(cmd))
			},
			Category:     item.Extend().Category,
			ArgsUsage:    item.Extend().ArgsUsage,
			Flags:        flagsToCliFlags(item.Extend().Flags),
			OnUsageError: onUsageError,
		}

		a.instance.Commands = append(a.instance.Commands, &cliCommand)
	}
}

// Call implements console.AiGoAgent.
func (a *Application) Call(command string) error {
	if len(os.Args) == 0 {
		return nil
	}

	commands := []string{os.Args[0]}

	if a.useAiGoAgent {
		commands = append(commands, "aigoagent")
	}

	return a.Run(append(commands, strings.Split(command, " ")...), false)
}

// CallAndExit implements console.AiGoAgent.
func (a *Application) CallAndExit(command string) {
	if len(os.Args) == 0 {
		return
	}

	commands := []string{os.Args[0]}

	if a.useAiGoAgent {
		commands = append(commands, "aigoagent")
	}

	_ = a.Run(append(commands, strings.Split(command, " ")...), true)
}

// Run implements console.AiGoAgent.
func (a *Application) Run(args []string, exitIfGoAgent bool) error {
	//nolint:staticcheck // branch is empty but it WIP
	if noANSI || env.IsNoANSI() || slices.Contains(args, "--no-ansi") {
		// TODO: color disable
	}

	aigoagentIndex := -1
	if a.useAiGoAgent {
		for i, arg := range args {
			if arg == "aigoagent" {
				aigoagentIndex = i
				break
			}
		}
	} else {
		aigoagentIndex = 0
	}

	if aigoagentIndex != -1 {
		if aigoagentIndex+1 == len(args) {
			args = append(args, "--help")
		}

		cliArgs := append([]string{args[0]}, args[aigoagentIndex+1:]...)
		if err := a.instance.Run(context.Background(), cliArgs); err != nil {
			if exitIfGoAgent {
				panic(err.Error())
			}

			return err
		}

		if exitIfGoAgent {
			os.Exit(0)
		}
	}

	return nil
}

func flagsToCliFlags(flags []command.Flag) []cli.Flag {
	var cliFlags []cli.Flag
	for _, flag := range flags {
		switch flag.Type() {
		case command.FlagTypeBool:
			flag := flag.(*command.BoolFlag)
			cliFlags = append(cliFlags, &cli.BoolFlag{
				Name:        flag.Name,
				Aliases:     flag.Aliases,
				HideDefault: flag.DisableDefaultText,
				Usage:       flag.Usage,
				Required:    flag.Required,
				Value:       flag.Value,
			})
		case command.FlagTypeFloat64:
			flag := flag.(*command.Float64Flag)
			cliFlags = append(cliFlags, &cli.FloatFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeFloat64Slice:
			flag := flag.(*command.Float64SliceFlag)
			cliFlags = append(cliFlags, &cli.FloatSliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    cli.NewFloatSlice(flag.Value...).Value(),
			})
		case command.FlagTypeInt:
			flag := flag.(*command.IntFlag)
			cliFlags = append(cliFlags, &cli.IntFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeIntSlice:
			flag := flag.(*command.IntSliceFlag)
			cliFlags = append(cliFlags, &cli.IntSliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeInt64:
			flag := flag.(*command.Int64Flag)
			cliFlags = append(cliFlags, &cli.Int64Flag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeInt64Slice:
			flag := flag.(*command.Int64SliceFlag)
			cliFlags = append(cliFlags, &cli.Int64SliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeString:
			flag := flag.(*command.StringFlag)
			cliFlags = append(cliFlags, &cli.StringFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeStringSlice:
			flag := flag.(*command.StringSliceFlag)
			cliFlags = append(cliFlags, &cli.StringSliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		}
	}

	return cliFlags
}
