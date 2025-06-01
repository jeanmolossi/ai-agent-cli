package config

import (
	"github.com/jeanmolossi/ai-agent-cli/app/contracts"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
	"github.com/jeanmolossi/ai-agent-cli/app/support"
)

type ServiceProvider struct{}

func (sp *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(contracts.BindingConfig, func(app foundation.Application) (any, error) {
		return NewApplication(support.EnvFilePath), nil
	})
}

func (sp *ServiceProvider) Boot(app foundation.Application) {}
