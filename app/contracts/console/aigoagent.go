package console

type AiGoAgent interface {
	// Register commands.
	Register(commands []Command)

	// Call run an AiGoAgent console command by name.
	Call(command string) error

	// CallAndExit run an AiGoAgent console command by name and exit.
	CallAndExit(command string)

	// Run a command. args include ["./bin", "aigoagent", "command"]
	Run(args []string, exitIfGoAgent bool) error
}
