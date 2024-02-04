package ui

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// Test calling the Escape key to close the program
func TestProgramQuit(test *testing.T) {
	// Launch the test program
	testModel := teatest.NewTestModel(
		test,
		initialiseModel(),
	)
	// We can only get the final model output when the program has finished.
	// As such, we sent the Escape key to end the program.
	testModel.Send(tea.KeyMsg{
		Type: tea.KeyEscape,
	})
	// Use type assertation to access the state of the model
	modelState, ok := testModel.FinalModel(test).(model)
	if ok != true {
		test.Error("Error fetching program state")
	}
	if modelState.Quitting != true {
		test.Errorf("Quit value is %t when it should be true!", modelState.Quitting)
	}
}

// https://github.com/charmbracelet/x/tree/main
// This test will check that upon the user entering a valid discord guild ID, the view will switch from the TextInput component to the Spinner component
func TestViewSwitching(test *testing.T) {
	// Launch the test program
	testModel := teatest.NewTestModel(
		test,
		initialiseModel(),
	)
	// Send an example discord guild ID (Integer 17-19 digits long)
	testModel.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("123456789123456789"),
	})
	// Send the enter key to switch to the next view
	testModel.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
	teatest.WaitFor(test, testModel.Output(), func(byts []byte) bool {
		return bytes.Contains(byts, []byte("Querying Mee6 data for Discord Guild"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*1))
}
