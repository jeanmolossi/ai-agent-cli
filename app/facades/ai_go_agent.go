package facades

import "github.com/jeanmolossi/ai-agent-cli/app/contracts/console"

func AiGoAgent() console.AiGoAgent {
	return App().MakeAiGoAgent()
}
