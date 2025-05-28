package agent

type (
	Config struct {
		Provider string
	}

	LLMOption func(cfg *Config)
)

func WithProvider(provider string) LLMOption {
	return func(cfg *Config) { cfg.Provider = provider }
}
