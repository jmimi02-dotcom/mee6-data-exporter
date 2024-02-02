package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// Creates a new Loading Spinner
func Spinner() spinner.Model {
        var spinnerComponent spinner.Model
	spinnerComponent = spinner.New()
	spinnerComponent .Spinner = spinner.Dot
	// ANSI-256 Colour Code for Aqua https://github.com/charmbracelet/lipgloss
	spinnerComponent .Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	return spinnerComponent 
}
