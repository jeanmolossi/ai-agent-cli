package providers

import (
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
	"github.com/jeanmolossi/ai-agent-cli/app/facades"
	"github.com/jeanmolossi/ai-agent-cli/infra/app/console"
)

type ConsoleServiceProvider struct{}

func (r *ConsoleServiceProvider) Register(app foundation.Application) {
	kernel := console.Kernel{}
	facades.AiGoAgent().Register(kernel.Commands())
}

func (r *ConsoleServiceProvider) Boot(app foundation.Application) {}
