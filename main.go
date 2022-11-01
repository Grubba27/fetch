package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"net/http"
	"os"
	"time"
)

var URL = ""
var SHOULD_QUIT = false

type (
	responseText string
	errMsg       struct {
		err error
	}
)

func (e errMsg) Error() string {
	return e.err.Error()
}

func callServer() tea.Msg {
	time.Sleep(2 * time.Second)
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(URL)
	if err != nil {
		return errMsg{err}
	}
	body, err := io.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		return nil
	}
	d := string(body)
	return printMsg(d)

}

type printMsg responseText

type model struct {
	spinner  spinner.Model
	quitting bool
	err      error
	isDone   bool
	result   string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, callServer)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case printMsg:
		m.isDone = true
		m.result = string(msg)
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n  %s  loading... \n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	if m.isDone {
		return fmt.Sprintf("\n\n Here is your data:\n%s\n\n", m.result)
	}
	return str
}

func initalizeLoader() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func main() {
	if len(os.Args) == 2 {
		URL = os.Args[1]
	} else {
		initializeInquerer()
	}
	if SHOULD_QUIT {
		fmt.Printf("Bro what have you done\n\n")
		os.Exit(1)
	}
	err := tea.NewProgram(initalizeLoader()).Start()
	if err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
