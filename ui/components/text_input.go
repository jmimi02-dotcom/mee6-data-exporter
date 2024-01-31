package components

import "github.com/charmbracelet/bubbles/textinput"

// Creates a new Text input
func TextInput() textinput.Model {
        var input textinput.Model
	input = textinput.New()
	input.Focus()
	input.CharLimit = 19
	input.Width = 19
	return input
}
