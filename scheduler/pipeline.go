package scheduler

import (
	"encoding/json"
	"image"
	"io"
	"os"
	"proj2/png"
	"strings"
)

type ImageData struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
}

func ImageTaskGenerator(dirs string) chan *png.ImageTask {
	imageTasks := make(chan *png.ImageTask)

	effectsPathFile := "../../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		panic(err)
	}

	go func() {
		defer close(imageTasks)

		var imageDataArr []ImageData
		reader := json.NewDecoder(effectsFile)

		for {
			var imageData ImageData
			err := reader.Decode(&imageData)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			imageDataArr = append(imageDataArr, imageData)
		}

		for _, dir := range strings.Split(dirs, "+") {
			for _, imageData := range imageDataArr {
				imageTask, err := png.Load("../../data/in/"+dir+"/"+imageData.InPath, imageData.OutPath, dir, imageData.Effects)
				if err != nil {
					panic(err)
				}
				imageTasks <- imageTask
			}
		}
	}()

	return imageTasks
}

func Worker(imageTasks chan *png.ImageTask, numMiniWorkers int) chan *png.ImageTask {
	c := make(chan *png.ImageTask)

	go func() {
		defer close(c)
		for img := range imageTasks {
			minX, minY, maxX, maxY := img.Bounds.Min.X, img.Bounds.Min.Y,
				img.Bounds.Max.X, img.Bounds.Max.Y

			for i, effect := range img.Effects {
				miniWg := make(chan bool)

				for w := 0; w < numMiniWorkers; w++ {

					go func(w int, miniWg chan bool) {
						var bounds image.Rectangle

						if w != numMiniWorkers-1 {
							bounds = image.Rect(
								minX,
								minY+w*maxY/numMiniWorkers,
								maxX,
								minY+(w+1)*maxY/numMiniWorkers,
							)
						} else {
							bounds = image.Rect(
								minX,
								minY+(numMiniWorkers-1)*maxY/numMiniWorkers,
								maxX,
								maxY,
							)
						}

						if effect == "G" {
							img.Grayscale(&bounds)
						} else if effect == "E" {
							img.EdgeDetection(&bounds)
						} else if effect == "S" {
							img.Sharpen(&bounds)
						} else if effect == "B" {
							img.Blur(&bounds)
						} else {
							panic("No matching effect")
						}

						miniWg <- true
					}(w, miniWg)
				}

				for j := 0; j < numMiniWorkers; j++ {
					<-miniWg
				}

				if i != len(img.Effects)-1 {
					img.SwapInAndOut()
				}

			}
			c <- img
		}
	}()

	return c
}

func ResultsAggregator(channels ...chan *png.ImageTask) {
	wg := make(chan bool)

	multiplex := func(c <-chan *png.ImageTask) {
		for img := range c {
			err := img.Save("../../data/out/" + img.Size + "_" + img.OutPath)
			if err != nil {
				panic(err)
			}
		}
		wg <- true
	}

	for _, c := range channels {
		go multiplex(c)
	}

	for i := 0; i < len(channels); i++ {
		<-wg
	}
}

func RunPipeline(config Config) {

	imageTasks := ImageTaskGenerator(config.DataDirs)

	channels := make([]chan *png.ImageTask, config.ThreadCount)
	for i := 0; i < config.ThreadCount; i++ {
		channels[i] = Worker(imageTasks, config.ThreadCount)
	}

	ResultsAggregator(channels...)
}
