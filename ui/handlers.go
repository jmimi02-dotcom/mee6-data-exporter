package ui

// Launch the program
func CreateProgram() {
	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		fmt.Println("Failed to launch the program:", err)
	}
}

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