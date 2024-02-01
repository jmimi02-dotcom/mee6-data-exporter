package ui

// Launch the program
func CreateProgram() {
	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		fmt.Println("Failed to launch the program:", err)
	}
}
