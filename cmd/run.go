package cmd

import (
	"fmt"

	"github.com/jeanmolossi/ai-agent-cli/internal/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executa o agente com um prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		ag, err := agent.New()
		if err != nil {
			return fmt.Errorf("fail on create agent: %w", err)
		}

		// args[0] is an prompt; if want support multiple parts, concat here
		return ag.Run(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// model flag
	runCmd.
		Flags().
		StringP("model", "m", "", "Model to use (ex: openai, anthropic, ollama). It overrides llm.provider")

	// temperature flag
	runCmd.
		Flags().
		Float64("temperature", 1.0, "LLM temperature")

	_ = viper.BindPFlag("llm.provider", runCmd.Flags().Lookup("model"))
	_ = viper.BindPFlag("llm.temperature", runCmd.Flags().Lookup("temperature"))
}
