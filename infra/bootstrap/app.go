package bootstrap

import (
	"github.com/jeanmolossi/ai-agent-cli/app/foundation"
	"github.com/jeanmolossi/ai-agent-cli/infra/config"
)

func Boot() {
	app := foundation.NewApplication()

	app.Boot()

	config.Boot()
}
