package facades

import "github.com/jeanmolossi/ai-agent-cli/app/contracts/config"

func Config() config.Config {
	return App().MakeConfig()
}
