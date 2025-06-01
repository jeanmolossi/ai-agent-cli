package foundation

import (
	"log/slog"

	"github.com/jeanmolossi/ai-agent-cli/app/config"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
	"github.com/jeanmolossi/ai-agent-cli/app/foundation/json"
	"github.com/jeanmolossi/ai-agent-cli/app/support"
	"github.com/jeanmolossi/ai-agent-cli/app/support/env"
)

var App foundation.Application

func init() {
	setEnv()
	setRootPath()

	app := &Application{
		Container:     NewContainer(),
		publishes:     make(map[string]map[string]string),
		publishGroups: make(map[string]map[string]string),
	}

	app.registerBaseServiceProviders()
	app.bootBaseServiceProviders()
	app.SetJson(json.New())
	App = app
}

type Application struct {
	*Container
	publishes     map[string]map[string]string
	publishGroups map[string]map[string]string
	json          foundation.Json
}

func NewApplication() foundation.Application {
	return App
}

func (a *Application) Boot() {
	a.setTimezone()
	a.registerConfiguredServiceProviders()
	a.bootConfiguredServiceProviders()
	// a.registerCommands([]contractsconsole.Command{
	// 	console.NewAboutCommand(a),
	// })
}

func (a *Application) SetJson(j foundation.Json) {
	if j != nil {
		a.json = j
	}
}

func (a *Application) GetJson() foundation.Json {
	return a.json
}

func (a *Application) getBaseServiceProviders() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		&config.ServiceProvider{},
	}
}

func (a *Application) getConfiguredServiceProviders() []foundation.ServiceProvider {
	configFacade := a.MakeConfig()
	if configFacade == nil {
		slog.Warn("config facade is not initialized. Skipping registering service providers.")
		return []foundation.ServiceProvider{}
	}

	providers, ok := configFacade.Get("app.providers").([]foundation.ServiceProvider)
	if !ok {
		slog.Warn("providers config is not of type []foundation.ServiceProvider. Skipping registering service providers.")
		return []foundation.ServiceProvider{}
	}

	return providers
}

func (a *Application) registerServiceProviders(serviceProviders []foundation.ServiceProvider) {
	for _, sp := range serviceProviders {
		sp.Register(a)
	}
}

func (a *Application) bootServiceProviders(serviceProviders []foundation.ServiceProvider) {
	for _, sp := range serviceProviders {
		sp.Boot(a)
	}
}

func (a *Application) registerConfiguredServiceProviders() {
	a.registerServiceProviders(a.getConfiguredServiceProviders())
}

func (a *Application) bootConfiguredServiceProviders() {
	a.bootServiceProviders(a.getConfiguredServiceProviders())
}

func (a *Application) registerBaseServiceProviders() {
	a.registerServiceProviders(a.getBaseServiceProviders())
}

func (a *Application) bootBaseServiceProviders() {
	a.bootServiceProviders(a.getBaseServiceProviders())
}

func (a *Application) setTimezone() {
	configFacade := a.MakeConfig()
	if configFacade == nil {
		slog.Warn("config facade is not initialized. Using default timezone UTC.")
		return
	}
}

func setEnv() {}

func setRootPath() {
	support.RootPath = env.CurrentAbsolutePath()
}
