package facades

import (
	foundationcontract "github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
	"github.com/jeanmolossi/ai-agent-cli/app/errors"
	"github.com/jeanmolossi/ai-agent-cli/app/foundation"
)

func App() foundationcontract.Application {
	if foundation.App == nil {
		panic(errors.ApplicationNotSet.SetModule(errors.ModuleFacade))
	}

	return foundation.App
}
