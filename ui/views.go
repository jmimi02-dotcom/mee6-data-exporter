package ui

import "fmt"

func inputView(m model) string {
	return fmt.Sprintf(
		"Discord Server ID\n\n%s",
		m.TextInput.View(),
	) + "\n"
}

func spinnerView(m model) string {
	label := m.Spinner.View() + "Querying Mee6 data for Discord Guild " + m.TextInput.Value()
	return "\n\n" + label
}
