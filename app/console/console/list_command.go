package console

import (
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/console"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/console/command"
)

type ListCommand struct {
	aigoagent console.AiGoAgent
}

func NewListCommand(aigoagent console.AiGoAgent) *ListCommand {
	return &ListCommand{
		aigoagent: aigoagent,
	}
}

// Signature The name and signature of the console command.
func (r *ListCommand) Signature() string {
	return "list"
}

// Description The console command description.
func (r *ListCommand) Description() string {
	return "List commands"
}

// Extend The console command extend.
func (r *ListCommand) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (r *ListCommand) Handle(ctx console.Context) error {
	return r.aigoagent.Call("--help")
}
