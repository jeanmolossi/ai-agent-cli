package main

import (
	"log"

	"github.com/jeanmolossi/ai-agent-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
