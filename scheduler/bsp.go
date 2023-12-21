package scheduler

import (
	"encoding/json"
	"image"
	"io"
	"os"
	"proj2/png"
	"strings"
)

type bspWorkerContext struct {
	// Define the necessary fields for your implementation
	workers   int
	dirs      string
	imageTask *png.ImageTask
	effect    string
	barrierS  Barrier
	barrierE  Barrier
}

func NewBSPContext(config Config) *bspWorkerContext {
	//Initialize the context
	return &bspWorkerContext{
		config.ThreadCount,
		config.DataDirs,
		nil,
		"",
		*NewBarrier(config.ThreadCount),
		*NewBarrier(config.ThreadCount),
	}
}

func generateTasks(ctx *bspWorkerContext, imageDataArr []ImageData) {

	images := make([]*png.ImageTask, 0)
	for _, dir := range strings.Split(ctx.dirs, "+") {
		for _, imageData := range imageDataArr {
			newImageTask, err := png.Load(
				"../../data/in/"+dir+"/"+imageData.InPath,
				imageData.OutPath,
				dir,
				imageData.Effects,
			)
			if err != nil {
				panic(err)
			}
			for i, e := range newImageTask.Effects {
				ctx.imageTask = newImageTask
				ctx.effect = e

				ctx.barrierS.Arrive()

				// main go routine works on the last chunk of the image as well
				{
					minX, minY, maxX, maxY := ctx.imageTask.Bounds.Min.X, ctx.imageTask.Bounds.Min.Y,
						ctx.imageTask.Bounds.Max.X, ctx.imageTask.Bounds.Max.Y
					bounds := image.Rect(
						minX,
						minY+(ctx.workers-1)*maxY/(ctx.workers),
						maxX,
						maxY,
					)
					if ctx.effect == "G" {
						ctx.imageTask.Grayscale(&bounds)
					} else if ctx.effect == "E" {
						ctx.imageTask.EdgeDetection(&bounds)
					} else if ctx.effect == "S" {
						ctx.imageTask.Sharpen(&bounds)
					} else if ctx.effect == "B" {
						ctx.imageTask.Blur(&bounds)
					} else {
						panic("No matching effect")
					}
				}

				ctx.barrierE.Arrive()

				if i != len(newImageTask.Effects)-1 {
					ctx.imageTask.SwapInAndOut()
				}
			}
			newImageTask.OutPath = "../../data/out/" +
				newImageTask.Size +
				"_" +
				newImageTask.OutPath
			images = append(images, newImageTask)
		}
	}

	barrier := NewBarrier(len(images) + 1)
	for i := 0; i < len(images); i++ {
		go func(i int) {
			images[i].Save(images[i].OutPath)
			barrier.Arrive()
		}(i)
	}
	barrier.Arrive()
}

func RunBSPWorker(id int, ctx *bspWorkerContext) {
	for {
		// Implement the BSP model here.
		// No additional loops can be used in this implementation. T
		// This goes to calling other functions. No other called
		// function you define or are using can have looping being done for you.
		if id == ctx.workers-1 {
			effectsPathFile := "../../data/effects.txt"
			effectsFile, err := os.Open(effectsPathFile)
			if err != nil {
				panic(err)
			}
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
			generateTasks(ctx, imageDataArr)
			return
		} else {
			ctx.barrierS.Arrive()
			minX, minY, maxX, maxY := ctx.imageTask.Bounds.Min.X, ctx.imageTask.Bounds.Min.Y,
				ctx.imageTask.Bounds.Max.X, ctx.imageTask.Bounds.Max.Y
			bounds := image.Rect(
				minX,
				minY+id*maxY/(ctx.workers),
				maxX,
				minY+(id+1)*maxY/(ctx.workers),
			)

			if ctx.effect == "G" {
				ctx.imageTask.Grayscale(&bounds)
			} else if ctx.effect == "E" {
				ctx.imageTask.EdgeDetection(&bounds)
			} else if ctx.effect == "S" {
				ctx.imageTask.Sharpen(&bounds)
			} else if ctx.effect == "B" {
				ctx.imageTask.Blur(&bounds)
			} else {
				panic("No matching effect")
			}
			ctx.barrierE.Arrive()

		}
	}
}
