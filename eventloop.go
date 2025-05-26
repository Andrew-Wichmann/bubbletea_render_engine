package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qeesung/image2ascii/convert"
)

const targetFps = 60
const linesGoRoutine = 10

type GraphicsEngine struct {
	frameNum       int64
	converter      convert.ImageConverter
	convertOptions convert.Options
	Frame          string
	FPS            int64
}

func NewGraphicEngine(width, height int) GraphicsEngine {
	converter := convert.NewImageConverter()
	convertOptions := convert.Options{
		FitScreen:       false,
		Colored:         true,
		Reversed:        false,
		StretchedScreen: false,
		FixedWidth:      width,
		FixedHeight:     height,
	}
	return GraphicsEngine{converter: *converter, convertOptions: convertOptions}
}

func (e *GraphicsEngine) Run() tea.Msg {
	start := time.Now().Unix() - 1
	for {
		e.frameNum += 1
		e.Frame = e.converter.ImageFile2ASCIIString(fmt.Sprintf("./frames/frame_%d.png", e.frameNum%9+1), &e.convertOptions)
		time.Sleep(time.Second / targetFps)
		e.FPS = e.frameNum / (time.Now().Unix() - start)
	}
}
