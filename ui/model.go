package ui

import (
	"fmt"
	"mee6xport/mee6"
	"mee6xport/ui/components"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

// This holds the application state
type model struct {
	TextInput     textinput.Model
	Spinner       spinner.Model
	InputEntered  bool
	Finished      bool
	Quitting      bool
	CurrentStatus string

	StartGenerating  bool
	Channel          chan mee6.Response
	CurrentPage      int
	ContinueCrawling bool
}

func (m model) Listen() tea.Cmd {
	if m.ContinueCrawling {
		return func() tea.Msg {
			for i := 0; i < 2; i++ {
				// At the time of writing I've been rate limited by the API, oops!
				// As such, we're having to reply on cached API responses I have stored in the /mock folder
				x, err := mee6.MockGetInfo(1234, i)
				//fmt.Println(i)
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(time.Second)
				m.Channel <- x
			}
			m.ContinueCrawling = false
			m.Finished = true
			return nil
		}
	} else {
		return nil
	}
}

type responseMsg mee6.Response

func waitForActivity(ch chan mee6.Response) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-ch)
	}
}

// Creates a new model{} structure, using default config
func initialiseModel() model {
	return model{
		TextInput:        components.TextInput(),
		Spinner:          components.Spinner(),
		InputEntered:     false,
		Finished:         false,
		Quitting:         false,
		CurrentStatus:    "",
		StartGenerating:  false,
		Channel:          make(chan mee6.Response),
		CurrentPage:      0,
		ContinueCrawling: true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, waitForActivity(m.Channel))
}

// Called as an event when an update is processed to the main application
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "esc" {
			m.Quitting = true
			return m, tea.Quit
		}
		m.TextInput, _ = m.TextInput.Update(msg)
	case responseMsg:
		if m.ContinueCrawling {
			m.CurrentStatus = fmt.Sprintf("Crawling %d", m.CurrentPage)
			m.CurrentPage++
			return m, waitForActivity(m.Channel)
		}
		m.Finished = true
		return m, nil
	}

	// If an input hasn't been entered, watch for the enter key being pressed
	if !m.InputEntered && !m.Finished {
		return setEntered(msg, m)
	}

	// If an input has been entered, start calling the mee6 api
	if m.InputEntered && !m.Finished {
		if m.StartGenerating {
			return m, nil
		} else {
			m.StartGenerating = false
			return m, m.Listen()
		}
	}

	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)

	return m, tea.Batch(cmd, m.Spinner.Tick)

}

func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  Written by Luis / github.com/luisjones\n\n"
	}

	s = inputView(m)

	if m.InputEntered {
		s = spinnerView(m)
	}

	return indent.String(fmt.Sprintf("\n%s\n\n", s), 2)
}

func (m model) isValidDiscordGuildID() bool {
	// Regular Expression returns true for digits with a length of 17-19 characters.
	/*
		TODO: Double check discord snowflake length
		Recently increased to 19 but check that this length is consistent across guilds and not variable
	*/
	regex, _ := regexp.Compile("\\d{17,19}")
	return regex.MatchString(m.TextInput.Value())
}
