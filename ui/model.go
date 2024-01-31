package ui

import (
	"mee6xport/ui/components"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	TextInput    textinput.Model
	Spinner      spinner.Model
	InputEntered bool
	Quitting     bool
}

// Creates a new model{} structure, using default config
func initialiseModel() model {
	return model{
		TextInput:    components.TextInput(),
		Spinner:      components.Spinner(),
		InputEntered: false,
		Quitting:     false,
	}
}

func (m model) Init() tea.Cmd {
	return m.Spinner.Tick
}
