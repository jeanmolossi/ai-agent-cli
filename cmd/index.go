package cmd

import (
	"log/slog"

	"github.com/jeanmolossi/ai-agent-cli/internal/rag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var docs []string

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Indexa o repositorio atual para o agente",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("Indexing files...")

		currentDir := []string{"."}
		cfgDocs := viper.GetStringSlice("rag.docs_paths")
		allDocs := append(cfgDocs, docs...)

		return rag.ScanAndIndex(append(currentDir, allDocs...)...)
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	// docs folder flag
	indexCmd.
		Flags().
		StringSliceVarP(&docs, "docs", "d", nil,
			"diretório(s) extra de documentação para indexar (sobrepõe config)")

		// bind
	_ = viper.BindPFlag("rag.docs_paths", indexCmd.Flags().Lookup("docs"))
}
