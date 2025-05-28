package rag

type (
	Config struct {
		Provider string
		Model    string
	}

	RagOption func(*Config)
)

func WithProvider(p string) RagOption {
	return func(c *Config) { c.Provider = p }
}

func WithModel(m string) RagOption {
	return func(c *Config) { c.Model = m }
}
