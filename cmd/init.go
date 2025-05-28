package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	initcmd "github.com/jeanmolossi/ai-agent-cli/internal/init"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa o agente (gera config padrão)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := filepath.Join(".", ".ai-agent-cli.yaml")

		_, err := os.ReadFile(cfgPath)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("unknown err file: %w", err)
		}

		if err != nil && errors.Is(err, fs.ErrNotExist) {
			m, err := tea.NewProgram(initcmd.InitialModel(), tea.WithAltScreen()).Run()
			if err != nil {
				return err
			}

			slog.Debug("final model", slog.Any("model", m))

			data, err := m.(initcmd.Model).Yaml()
			if err != nil {
				return err
			}

			if err := os.WriteFile(cfgPath, data, 0644); err != nil {
				return err
			}

			fmt.Printf("Config criado em %s\n", cfgPath)
		} else {
			fmt.Printf("O arquivo de configuração já existe.\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	// flags locais de init, se precisar:
	// initCmd.Flags().StringP("template", "t", "", "template a usar")
}
