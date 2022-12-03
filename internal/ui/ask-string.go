package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type StringQuestionModel struct {
	prompt       string
	defaultValue string
	input        textinput.Model
	quitting     bool
}

func AskString(prompt, defaultValue string) string {
	input := textinput.New()
	input.Placeholder = defaultValue
	input.Prompt = ""
	input.Focus()

	model := StringQuestionModel{
		prompt:       prompt,
		defaultValue: defaultValue,
		input:        input,
		quitting:     false,
	}

	p := tea.NewProgram(model)
	if lastModel, err := p.Run(); err != nil {
		panic(err)
	} else {
		model = lastModel.(StringQuestionModel)
	}

	return model.input.Value()
}

func (m StringQuestionModel) Init() tea.Cmd {
	return nil
}

func (m StringQuestionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			os.Exit(1)
		}

		if msg.Type == tea.KeyEnter {
			if m.input.Value() != "" {
				m.quitting = true
				return m, tea.Quit
			}

			if m.input.Value() == "" && m.defaultValue != "" {
				m.input.SetValue(m.defaultValue)
				m.quitting = true
				return m, tea.Quit
			}
		}

		if msg.Type == tea.KeyTab {
			if m.input.Value() == "" {
				m.input.SetValue(m.defaultValue)
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m StringQuestionModel) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf("%s %s\n", m.prompt, m.input.View())
}
