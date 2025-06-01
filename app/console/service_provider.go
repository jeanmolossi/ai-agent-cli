package console

import (
	"log/slog"

	"github.com/jeanmolossi/ai-agent-cli/app/console/console"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts"
	consolecontract "github.com/jeanmolossi/ai-agent-cli/app/contracts/console"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
)

type ServiceProvider struct{}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(contracts.BindingConsole, func(app foundation.Application) (any, error) {
		name := "aigoagent"
		usage := "AiGoAgent Framework"
		usageText := "aigoagent [global options] command [options] [arguments...]"
		return NewApplication(name, usage, usageText, app.Version(), true), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	r.registerCommands(app)
}

func (r *ServiceProvider) registerCommands(app foundation.Application) {
	aigoagentFacade := app.MakeAiGoAgent()
	if aigoagentFacade == nil {
		slog.Warn("AiGoAgent Facade is not initialized. Skipping command registration.")
		return
	}

	configFacade := app.MakeConfig()
	if configFacade == nil {
		slog.Warn("Config Facade is not initialized. Skipping certain command registrations.")
		return
	}

	aigoagentFacade.Register([]consolecontract.Command{
		console.NewListCommand(aigoagentFacade),
		console.NewKeyGenerateCommand(configFacade),
		// console.NewMakeCommand(),
		// console.NewBuildCommand(configFacade),
	})
}
