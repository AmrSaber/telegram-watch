package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type BoolQuestionModel struct {
	prompt       string
	defaultValue bool
	response     bool
	quitting     bool
}

func AskBool(prompt string, defaultValue bool) bool {
	model := BoolQuestionModel{
		prompt:       prompt,
		defaultValue: defaultValue,
		quitting:     false,
	}

	p := tea.NewProgram(model)
	if lastModel, err := p.Run(); err != nil {
		panic(err)
	} else {
		model = lastModel.(BoolQuestionModel)
	}

	return model.response
}

func (m BoolQuestionModel) Init() tea.Cmd {
	return nil
}

func (m BoolQuestionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			os.Exit(1)
		}

		if msg.Type == tea.KeyEnter {
			m.response = m.defaultValue
			m.quitting = true
			return m, tea.Quit
		}

		if strings.ToLower(msg.String()) == "y" {
			m.response = true
			m.quitting = true
			return m, tea.Quit
		}

		if strings.ToLower(msg.String()) == "n" {
			m.response = false
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m BoolQuestionModel) View() string {
	if m.quitting {
		return ""
	}

	y, n := "y", "n"
	if m.defaultValue {
		y = "Y"
	} else {
		n = "N"
	}

	return fmt.Sprintf("%s [%s/%s]\n", m.prompt, y, n)
}
