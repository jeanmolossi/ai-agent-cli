package initcmd

type stage int

func (s stage) isLastStage() bool {
	return s == quit
}

const (
	llmProvider stage = iota
	ragProvider
	logLevel
	quit

	ollamaHost stage = iota
	ollamaPort
	ollamaModel
	ollamaTemp
)
