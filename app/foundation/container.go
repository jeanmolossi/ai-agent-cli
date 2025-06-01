package foundation

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/jeanmolossi/ai-agent-cli/app/contracts"
	contractsconfig "github.com/jeanmolossi/ai-agent-cli/app/contracts/config"
	contractsconsole "github.com/jeanmolossi/ai-agent-cli/app/contracts/console"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
)

type instance struct {
	concrete any
	shared   bool
}

type Container struct {
	bindings  sync.Map
	instances sync.Map
}

func NewContainer() *Container {
	return &Container{}
}

func (r *Container) Bind(key any, callback func(app foundation.Application) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: false})
}

func (r *Container) BindWith(key any, callback func(app foundation.Application, parameters map[string]any) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: false})
}

func (r *Container) Fresh(bindings ...any) {
	if len(bindings) == 0 {
		r.instances.Range(func(key, value any) bool {
			if key != contracts.BindingConfig {
				r.instances.Delete(key)
			}

			return true
		})
	} else {
		for _, binding := range bindings {
			r.instances.Delete(binding)
		}
	}
}

func (r *Container) Instance(key any, inst any) {
	r.bindings.Store(key, instance{concrete: inst, shared: true})
}

func (r *Container) Make(key any) (any, error) {
	return r.make(key, nil)
}

func (r *Container) MakeAiGoAgent() contractsconsole.AiGoAgent {
	instance, err := r.Make(contracts.BindingConsole)
	if err != nil {
		return nil
	}

	return instance.(contractsconsole.AiGoAgent)
}

func (r *Container) MakeConfig() contractsconfig.Config {
	instance, err := r.Make(contracts.BindingConfig)
	if err != nil {
		slog.Error(
			"error making config",
			slog.String("err", err.Error()),
		)
		return nil
	}

	return instance.(contractsconfig.Config)
}

func (r *Container) MakeWith(key any, parameters map[string]any) (any, error) {
	return r.make(key, parameters)
}

func (r *Container) Singleton(key any, callback func(app foundation.Application) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: true})
}

func (r *Container) make(key any, parameters map[string]any) (any, error) {
	binding, ok := r.bindings.Load(key)
	if !ok {
		return nil, fmt.Errorf("binding not found: %+v", key)
	}

	if parameters == nil {
		instance, ok := r.instances.Load(key)
		if ok {
			return instance, nil
		}
	}

	bindingImpl := binding.(instance)
	switch concrete := bindingImpl.concrete.(type) {
	case func(app foundation.Application) (any, error):
		concreteImpl, err := concrete(App)
		if err != nil {
			return nil, err
		}

		if bindingImpl.shared {
			r.instances.Store(key, concreteImpl)
		}

		return concreteImpl, nil

	case func(app foundation.Application, parameters map[string]any) (any, error):
		concreteImpl, err := concrete(App, parameters)
		if err != nil {
			return nil, err
		}

		return concreteImpl, nil
	default:
		r.instances.Store(key, concrete)
		return concrete, nil
	}
}
