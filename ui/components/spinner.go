package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// Creates a new Loading Spinner
func Spinner() spinner.Model {
        var spinner spinner.Model
	spinner = spinner.New()
	spinner.Spinner = spinner.Dot
	// ANSI-256 Colour Code for Aqua https://github.com/charmbracelet/lipgloss
	spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	return spinner
}
