package main

import (
	"fmt"
	"image"
	"os"
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

	// temp
	frame_images []image.Image
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
	// temp
	frames, err := os.ReadDir("./frames")
	if err != nil {
		panic(err)
	}
	for _, frame := range frames {
		image, err := OpenImageFile(fmt.Sprintf("%s/%s", "./frames/", frame.Name()))
		if err != nil {
			panic(err)
		}
		e.frame_images = append(e.frame_images, image)
	}

	start := time.Now().Unix() - 1
	for {
		e.frameNum += 1
		e.Frame = e.converter.Image2ASCIIString(e.frame_images[e.frameNum%8+1], &e.convertOptions)
		time.Sleep(time.Second / targetFps)
		e.FPS = e.frameNum / (time.Now().Unix() - start)
	}
}

func OpenImageFile(imageFilename string) (image.Image, error) {
	f, err := os.Open(imageFilename)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return img, nil
}
