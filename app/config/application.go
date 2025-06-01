package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/jeanmolossi/ai-agent-cli/app/contracts/config"
	"github.com/jeanmolossi/ai-agent-cli/app/support"
	"github.com/jeanmolossi/ai-agent-cli/app/support/file"
	"github.com/spf13/viper"
)

var _ config.Config = &Application{}

type Application struct {
	vip *viper.Viper
}

func NewApplication(envFilePath string) *Application {
	app := &Application{}
	app.vip = viper.New()
	app.vip.AutomaticEnv()

	if file.Exists(envFilePath) {
		app.vip.SetConfigType("env")
		app.vip.SetConfigFile(envFilePath)

		if err := app.vip.ReadInConfig(); err != nil {
			slog.Error("invalid config error: %v", err)
			os.Exit(0)
		}
	}

	appKey := app.Env("APP_KEY")
	if !support.DontVerifyEnvFileExists {
		if appKey == nil {
			slog.Error("Please initialize APP_KEY first.")
			slog.Error("Create a .env file and run key:generate command.")
			slog.Error("Or set a system variable: APP_KEY={32-bit number}")

			os.Exit(0)
		}

		if len(appKey.(string)) != 32 {
			slog.Error("Invalid APP_KEY, the length must be 32, please reset it.")
			os.Exit(0)
		}
	}

	return app
}

// Add implements config.Config.
func (a *Application) Add(name string, configuration any) {
	panic("unimplemented")
}

// Env implements config.Config.
func (a *Application) Env(envName string, defaultValue ...any) any {
	panic("unimplemented")
}

// Get implements config.Config.
func (a *Application) Get(path string, defaultValue ...any) any {
	panic("unimplemented")
}

// GetBool implements config.Config.
func (a *Application) GetBool(path string, defaultValue ...bool) bool {
	panic("unimplemented")
}

// GetDuration implements config.Config.
func (a *Application) GetDuration(path string, defaultValue ...time.Duration) time.Duration {
	panic("unimplemented")
}

// GetInt implements config.Config.
func (a *Application) GetInt(path string, defaultValue ...int) int {
	panic("unimplemented")
}

// GetString implements config.Config.
func (a *Application) GetString(path string, defaultValue ...string) string {
	panic("unimplemented")
}
