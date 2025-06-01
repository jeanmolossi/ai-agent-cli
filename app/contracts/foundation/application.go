package foundation

import "github.com/jeanmolossi/ai-agent-cli/app/contracts/config"

type Application interface {
	// Boot register and bootstrap configured service providers.
	Boot()

	// Path
	SetJson(j Json)
	GetJson() Json

	// Container
	// Bind registers a binding with the container.
	Bind(key any, callback func(app Application) (any, error))
	// BindWith registers a binding with the container.
	BindWith(key any, callback func(app Application, parameters map[string]any) (any, error))
	// Fresh modules after changing config, will fresh all bindings except for config if no bindings provided.
	// Notice, this method only freshs the facade, if another facade injects the facade previously, the another
	// facades should be fresh simulaneously.
	Fresh(bindings ...any)
	// Instance registers an existing instance as shared in the container.
	Instance(key, instance any)
	// MakeConfig resolves the config instance.
	MakeConfig() config.Config
	// Make resolves the given type from the container.
	Make(key any) (any, error)
	// MakeWith resolves the given type with the given parameters from the container.
	MakeWith(key any, parameters map[string]any) (any, error)
	// Singleton registers a shared binding in the container.
	Singleton(key any, callback func(app Application) (any, error))
}
