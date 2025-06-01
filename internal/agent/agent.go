package agent

import (
	"context"
	"fmt"
	"strings"

	contractsagent "github.com/jeanmolossi/ai-agent-cli/internal/contracts/agent"
	"github.com/jeanmolossi/ai-agent-cli/internal/prompt"
	"github.com/spf13/viper"
)

// The Agent represents CLI's core,
// will be responsible for load LLMs, configure context, etc.
type Agent struct {
	llm       contractsagent.LLMProvider
	templates []string
}

func New() (*Agent, error) {
	provider, err := NewLLMProvider()
	if err != nil {
		return nil, err
	}

	tplPath := viper.GetString("prompt.templates_path")
	tplContents, err := prompt.LoadTemplates(tplPath)
	if err != nil {
		return nil, err
	}

	return &Agent{
		llm:       provider,
		templates: tplContents,
	}, nil
}

// Run is how the agent perform an action. It load templates and generate an
// LLM request with loaded templates and user prompt.
func (a *Agent) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please pass an action")
	}

	userPrompt := strings.Join(args, " ")
	fullPrompt := a.LoadTemplates(userPrompt)

	resp, err := a.llm.Generate(context.Background(), fullPrompt)
	if err != nil {
		return err
	}

	fmt.Println("---------------------------------------------------------------")
	fmt.Println()
	fmt.Println(resp)
	fmt.Println()
	fmt.Println("---------------------------------------------------------------")
	return nil
}
