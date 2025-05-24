package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	e *EventLoop
}

func (a app) Init() tea.Cmd {
	return a.e.Run
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a app) View() string {
	return fmt.Sprintf("%s\nFPS:%d", a.e.Frame, a.e.FPS)
}

func main() {
	eventLoop := NewEventLoop()
	a := app{e: &eventLoop}
	prog := tea.NewProgram(a)
	prog.Run()
}
