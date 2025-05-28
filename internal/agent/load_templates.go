package agent

import (
	"log/slog"
	"strings"
)

func (a *Agent) LoadTemplates(userPrompt string) string {
	var fullPrompt strings.Builder
	for _, tpl := range a.templates {
		fullPrompt.WriteString(tpl)
		fullPrompt.WriteString("\n\n")
	}
	fullPrompt.WriteString("Usuario: ")
	fullPrompt.WriteString(userPrompt)

	slog.Debug("final prompt", slog.String("prompt", fullPrompt.String()))

	return fullPrompt.String()
}
