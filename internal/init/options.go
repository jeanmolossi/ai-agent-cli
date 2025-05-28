package initcmd

import "github.com/charmbracelet/bubbles/list"

var (
	defaultWidth = 80
	listHeight   = 21
)

var llmOptions = []list.Item{
	item{"openai", "API OpenAI (ex: GPT-4)"},
	item{"anthropic", "API Anthropic"},
	item{"ollama", "Servidor Ollama"},
}

var llmOllamaOptions = []list.Item{
	item{"host", "O host do servidor do ollama (default http://localhost"},
	item{"port", "A porta do servidor do ollama (default 11434)"},
	item{"model", "O modelo que o ollama vai requisitar (default gemma)"},
	item{"temperature", "A temperatura de criatividade do modelo (default 0)"},
}

var ragOptions = []list.Item{
	item{"local", "Vector store SQLite local"},
	item{"qdrant", "Qdrant (via HTTP API)"},
	item{"pgvector", "PostgreSQL + PGVector"},
}

var logLevelOptions = []list.Item{
	item{"debug", "Printa todos os logs, incluindo os de debug"},
	item{"info", "Performa logs com level iguais ou superior a info (info, warn, error)"},
	item{"warn", "Performa logs com level iguais ou superior a warn (warn, error)"},
	item{"error", "Performa apenas logs de erro"},
}

func defaultUI(title string, options []list.Item) list.Model {
	ll := list.New(options, list.NewDefaultDelegate(), defaultWidth, listHeight)
	ll.Title = title
	ll.SetShowStatusBar(false)
	ll.SetFilteringEnabled(false)

	return ll
}
