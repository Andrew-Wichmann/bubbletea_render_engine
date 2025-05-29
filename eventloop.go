package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qeesung/image2ascii/convert"
)

const targetFps = 60
const goRoutines = 27

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
		FixedHeight:     height / goRoutines,
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
	var wg sync.WaitGroup
	for {
		e.frameNum += 1
		rect := e.frame_images[1].Bounds()
		results := make(chan workerResult, rect.Max.Y/goRoutines)
		for i := range goRoutines {
			wg.Add(1)
			go e.worker(i, results, &wg)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		e.Frame = collectResults(results)
		time.Sleep(time.Second / targetFps)
		e.FPS = e.frameNum / (time.Now().Unix() - start)
	}
}

func collectResults(results chan workerResult) string {
	collectedResults := make([]string, goRoutines)
	var builder strings.Builder
	builder.Grow(30000 * goRoutines)
	for results := range results {
		collectedResults[results.lineNumber] = results.renderedString
	}
	for _, line := range collectedResults {
		builder.WriteString(line)
	}
	return builder.String()
}

type workerResult struct {
	lineNumber     int
	renderedString string
}

func (e GraphicsEngine) worker(id int, results chan workerResult, wg *sync.WaitGroup) {
	img := e.frame_images[e.frameNum%8+1]
	defer wg.Done()
	rect := img.Bounds()
	cropArea := image.Rectangle{image.Pt(0, id*(rect.Max.Y/goRoutines)), image.Pt(rect.Max.X, min(((id+1)*(rect.Max.Y/goRoutines)), rect.Max.Y))}
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropArea)
	result := e.converter.Image2ASCIIString(croppedImg, &e.convertOptions)
	results <- workerResult{lineNumber: id, renderedString: result}
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
