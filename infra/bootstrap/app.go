package bootstrap

import (
	"github.com/jeanmolossi/ai-agent-cli/app/foundation"
	"github.com/jeanmolossi/ai-agent-cli/infra/config"
)

func Boot() {
	app := foundation.NewApplication()

	// Bootstrap the application.
	app.Boot()

	// Bootstrap the config.
	config.Boot()
}
