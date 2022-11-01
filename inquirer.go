package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func initializeInquerer() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type (
	inqErr error
)

type input struct {
	textInput textinput.Model
	err       error
}

func initialModel() input {
	ti := textinput.New()
	ti.Placeholder = "someurl.com"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return input{
		textInput: ti,
		err:       nil,
	}
}

func (m input) Init() tea.Cmd {
	return textinput.Blink
}

func (m input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			SHOULD_QUIT = true
			return m, tea.Quit
		case tea.KeyEnter:
			URL = m.textInput.Value()
			return m, tea.Quit

		}

	// We handle errors just like any other message
	case inqErr:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m input) View() string {
	return fmt.Sprintf(
		"Your new postman is here.\n Type the api you want to fetch\n\n%s\n\n%s",
		m.textInput.View(),
		"(enter to fetch)",
	) + "\n"
}
