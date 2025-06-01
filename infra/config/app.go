package config

import (
	"github.com/dromara/carbon/v2"
	"github.com/jeanmolossi/ai-agent-cli/app/contracts/foundation"
	"github.com/jeanmolossi/ai-agent-cli/app/facades"
)

func Boot() {}

func init() {
	config := facades.Config()
	config.Add("app", map[string]any{
		// Application Name
		//
		// This value is the name of your application. This is used when the
		// framework needs to place the application's name in a notification or
		// any other location as required by the application or its packages.
		"name": config.Env("APP_NAME", "AiGoAgent"),

		// Application Environment
		//
		"env": config.Env("APP_ENV", "production"),

		// Application Debug Mode
		"debug": config.Env("APP_DEBUG", false),

		// Application timezone
		//
		// Here you may specify the default timezone for your application.
		// Example: UTC, America/Sao_Paulo
		// More: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
		"timezone": carbon.UTC,

		// Application Locale configuration
		//
		// The application locale determines the default locale that will be used
		// by the translation service provider. You are free to set this value
		// to any of the locales which will be supported by the application.
		"locale": "en",

		// Application Fallback locale
		//
		// The fallback locale determines the locale to use when the current one
		// is not available. You may change the value to correspond to any of
		// the language folders that are provided through your application.
		"fallback_locale": "en",

		// Application Lang path
		//
		// The path to the language files for the application. You may change
		// the path to a different directory if you would like to customize it.
		"lang_path": "lang",

		// Encryption Key
		//
		// 32 character string, otherwise these encrypted strings
		// will not be safe. Please do this before deploying an application!
		"key": config.Env("APP_KEY", ""),

		// Autoload service providers
		//
		// The service providers listed here will be automatically loaded on the
		// request to your application. Feel free to add your own services to
		// this array to grant expanded functionality to your applications.
		"providers": []foundation.ServiceProvider{},
	})
}
