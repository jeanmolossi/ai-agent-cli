package cmd

import (
	"fmt"
	"strings"

	"github.com/jeanmolossi/ai-agent-cli/internal/agent"
	"github.com/jeanmolossi/ai-agent-cli/internal/rag"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query [pergunta]",
	Short: "Faz uma pergunta acerca do repositorio",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := rag.NewVectorStore()
		if err != nil {
			return err
		}

		//nolint:errcheck
		defer store.Close()

		err = store.Load()
		if err != nil {
			return err
		}

		results, err := store.Search(strings.Join(args, " "), 100)
		if err != nil {
			return err
		}

		var ctxBuilder strings.Builder
		for _, r := range results {
			ctxBuilder.WriteString(r.Content() + "\n---\n")
		}

		prompt := fmt.Sprintf("Contexto:\n%s\n\nPergunta: %s", ctxBuilder.String(), args[0])
		ag, err := agent.New()
		if err != nil {
			return err
		}

		err = ag.Run([]string{prompt})
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
