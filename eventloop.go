package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qeesung/image2ascii/convert"
)

const targetFps = 60

type EventLoop struct {
	frameNum  int64
	converter convert.ImageConverter
	Frame     string
	FPS       int64
}

func NewEventLoop() EventLoop {
	converter := convert.NewImageConverter()
	return EventLoop{converter: *converter}
}

func (e *EventLoop) Run() tea.Msg {
	start := time.Now().Unix() - 1
	for {
		e.frameNum += 1
		e.Frame = e.converter.ImageFile2ASCIIString(fmt.Sprintf("/home/twisted/Code/playing_with_goroutines/frames/frame_%d.png", e.frameNum%9+1), &convert.DefaultOptions)
		time.Sleep(time.Second / targetFps)
		e.FPS = e.frameNum / (time.Now().Unix() - start)
	}
}
