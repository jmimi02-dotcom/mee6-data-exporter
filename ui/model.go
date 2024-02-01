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

// Called as an event when an update is processed to the main application
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Get the keypress event
	if msg, ok := msg.(tea.KeyMsg); ok {
		key := msg.String()
		// If the key pressed is escape or control c, quit the application
		if key == "ctrl+c" || key == "esc" {
			m.Quitting = true
			return m, tea.Quit
		}
		// Ensure that the TextInput field is updated with the key that has been entered
		m.TextInput, _ = m.TextInput.Update(msg)
	}

	// If an input hasn't been entered, watch for the enter key being pressed
	if !m.InputEntered {
		return setEntered(msg, m)
	}
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  Written by Luis / github.com/luisjones\n\n"
	}
	if !m.InputEntered {
		s = inputView(m)
	} else {
		s = spinnerView(m)
	}
	return indent.String(fmt.Sprintf("\n%s"\n\n", s), 2)
}
