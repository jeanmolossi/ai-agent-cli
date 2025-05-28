package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jeanmolossi/ai-agent-cli/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.1.0"

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "ai-agent-cli",
		Short: "Agente de IA para auxiliar no desenvolvimento de software",
		Long: `CLI para orquestrar e configurar um agente de IA 
que pode carregar modelos locais ou remotos.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// global flags, ex: --config
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de configuração (default ./.ai-agent-cli.yaml ou ~/.ai-agent-cli.yaml)")
	rootCmd.PersistentFlags().String("log-level", "info", "nível do logger (debug, info, warn error)")
	rootCmd.PersistentFlags().String("log-format", "text", "formato de saida do logger (text, json)")

	// binding of flags
	_ = viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("log-format"))

	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.Flags().BoolP("version", "v", false, "exibe a versão")

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if ver, _ := cmd.Flags().GetBool("version"); ver {
			fmt.Println("ai-agent-cli version", version)
			os.Exit(0)
		}
	}
}

func initConfig() {
	cfgFilename := ".ai-agent-cli"

	if cfgFile != "" {
		// usuario definiu um config file
		viper.SetConfigFile(cfgFile)
	} else {
		// por default procura o config file no diretorio onde ta rodando
		viper.AddConfigPath(".")

		// como fallback, procura na home do usuario
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "não foi possível determinar a home dir: %v\n", err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)

		viper.SetConfigName(cfgFilename)
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Não foi possível ler a config", cfgFilename)
	}

	slog.Info("Usando config file", slog.String("file", viper.ConfigFileUsed()))
}

func initLogger() {
	lvl := viper.GetString("log.level")
	if err := log.Init(lvl); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}
}
