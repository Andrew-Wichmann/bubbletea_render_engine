package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type State int

const INIT, RUNNING, PAUSED State = 0, 1, 2

type app struct {
	e     *GraphicsEngine
	state State
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (a app) Init() tea.Cmd {
	return doTick()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
	}
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		engine := NewGraphicEngine(msg.Width, msg.Height)
		a.e = &engine
		a.state = RUNNING
		return a, engine.Run
	}
	if _, ok := msg.(TickMsg); ok {
		return a, doTick()
	}
	return a, nil
}

func (a app) View() string {
	if a.state == INIT {
		return ""
	} else {
		return fmt.Sprintf("%s\nFPS:%d", a.e.Frame(), a.e.FPS)
	}
}

func main() {
	a := app{state: INIT}
	f, err := tea.LogToFile("./app.log", "")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	prog := tea.NewProgram(a)
	prog.Run()
}
