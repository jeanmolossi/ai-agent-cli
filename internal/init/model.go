package initcmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

type Model struct {
	stage        stage
	llmList      list.Model
	ragList      list.Model
	loglvlList   list.Model
	llmChoice    string
	ragChoice    string
	loglvlChoice string
}

func InitialModel() Model {
	return Model{
		stage:      llmProvider,
		llmList:    defaultUI("Selecione o provedor de LLM", llmOptions),
		ragList:    defaultUI("Selecione o provedor de Vector Store", ragOptions),
		loglvlList: defaultUI("Selecione o log level", logLevelOptions),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.stage {
	case llmProvider:
		m.llmList, cmd = m.llmList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.llmChoice = m.llmList.SelectedItem().(item).title
				m.stage = ragProvider

				if m.stage.isLastStage() {
					return m, tea.Quit
				}
			}
		}
		return m, cmd

	case ragProvider:
		m.ragList, cmd = m.ragList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.ragChoice = m.ragList.SelectedItem().(item).title
				m.stage = logLevel

				if m.stage.isLastStage() {
					return m, tea.Quit
				}
			}
		}
		return m, cmd

	case logLevel:
		m.loglvlList, cmd = m.loglvlList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.loglvlChoice = m.loglvlList.SelectedItem().(item).title
				m.stage = quit

				if m.stage.isLastStage() {
					return m, tea.Quit
				}
			}
		}
		return m, cmd

	default:
		return m, tea.Quit
	}
}

func (m Model) View() string {
	switch m.stage {
	case llmProvider:
		return m.llmList.View() + "\n(Use arrow keys e Enter)"
	case ragProvider:
		return m.ragList.View() + "\n(Use arrow keys e Enter)"
	case logLevel:
		return m.loglvlList.View() + "\n(Use arrow keys e Enter)"
	case quit:
		return fmt.Sprintf("Gerando configuração...\nLLM: %s\nRAG: %s\nLog Level: %s\n", m.llmChoice, m.ragChoice, m.loglvlChoice)
	}

	return ""
}

func (m Model) ollamaTpl() map[string]any {
	return map[string]any{
		"provider": m.llmChoice,
		"ollama": map[string]any{
			"host":        "http://localhost",
			"port":        11434,
			"model":       "gemma",
			"temperature": 0.0,
		},
		"openai":    m.genericTpl()["openai"],
		"anthropic": m.genericTpl()["openai"],
	}
}

func (m Model) genericTpl() map[string]any {
	return map[string]any{
		"provider": m.llmChoice,
		"openai": map[string]any{
			"api_key": "sk-XXXXXXXXXXXXXXXXXXXXXXXX",
		},
	}
}

func (m Model) ragTpl() map[string]any {
	return map[string]any{
		"provider":       m.ragChoice,
		"embed_provider": m.llmChoice,
		"ignore": []string{
			".git",
			".vscode",
			".docker",
			".idea",
			"node_modules",
			"vendor",
		},
		"local": map[string]any{
			"chunk_size": 512,
		},
		"docs_paths": []string{
			"docs",
		},
	}
}

func (m Model) Yaml() ([]byte, error) {
	llm := m.genericTpl()

	if m.llmChoice == "ollama" {
		llm = m.ollamaTpl()
	}

	cfg := map[string]any{
		"llm": llm,
		"rag": m.ragTpl(),
		"prompt": map[string]any{
			"templates_path": "./templates",
		},
		"log": map[string]any{
			"level":  m.loglvlChoice,
			"format": "text",
		},
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Println("--------------------------------")
	fmt.Println()
	fmt.Println(string(data))
	fmt.Println()
	fmt.Println("--------------------------------")

	return data, nil
}
