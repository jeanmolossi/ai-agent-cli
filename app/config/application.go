package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/jeanmolossi/ai-agent-cli/app/contracts/config"
	"github.com/jeanmolossi/ai-agent-cli/app/support"
	"github.com/jeanmolossi/ai-agent-cli/app/support/convert"
	"github.com/jeanmolossi/ai-agent-cli/app/support/file"
	"github.com/spf13/cast"
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
			slog.Error("invalid config error", slog.String("err", err.Error()))
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
			slog.Warn("Example command: \n./aigoagent key:generate")
			os.Exit(0)
		}
	}

	return app
}

// Env implements config.Config.
func (app *Application) Env(envName string, defaultValue ...any) any {
	value := app.Get(envName, defaultValue...)
	if cast.ToString(value) == "" {
		return convert.Default(defaultValue...)
	}

	return value
}

// Add implements config.Config.
func (app *Application) Add(name string, configuration any) {
	app.vip.Set(name, configuration)
}

// Get implements config.Config.
func (app *Application) Get(path string, defaultValue ...any) any {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}

	return app.vip.Get(path)
}

// GetBool implements config.Config.
func (app *Application) GetBool(path string, defaultValue ...bool) bool {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}

	return app.vip.GetBool(path)
}

// GetDuration implements config.Config.
func (app *Application) GetDuration(path string, defaultValue ...time.Duration) time.Duration {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}

	return app.vip.GetDuration(path)
}

// GetInt implements config.Config.
func (app *Application) GetInt(path string, defaultValue ...int) int {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}

	return app.vip.GetInt(path)
}

// GetString implements config.Config.
func (app *Application) GetString(path string, defaultValue ...string) string {
	if !app.vip.IsSet(path) {
		return convert.Default(defaultValue...)
	}
	return app.vip.GetString(path)
}
