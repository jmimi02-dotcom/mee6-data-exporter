package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func setEntered(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	// If the enter key has been presssed while an input hasn't been entered, return the spinner view
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.InputEntered = true
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}
