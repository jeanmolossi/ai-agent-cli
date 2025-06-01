package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jeanmolossi/ai-agent-cli/infra/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	if err := facades.Route().Run(); err != nil {
	// 		facades.Log().Errorf("Route run error: %v", err)
	// 	}
	// }()

	// Listen for the OS signal
	go func() {
		<-quit

		// if err := facades.Route().Shutdown(); err != nil {
		//     facades.Log().Errorf("Route Shutdown error: %v", err)
		// }

		os.Exit(0)
	}()

	select {}
}
