package scheduler

import (
	"encoding/json"
	"image"
	"io"
	"os"
	"proj2/png"
	"strings"
)

func RunSequential(config Config) {
	images := make([]ImageData, 0)
	effectsPathFile := "../../data/effects.txt"
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		panic(err)
	}

	reader := json.NewDecoder(effectsFile)
	for {
		var newImage ImageData
		err := reader.Decode(&newImage)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		images = append(images, newImage)
	}

	for _, dir := range strings.Split(config.DataDirs, "+") {
		for _, currImage := range images {
			pngImg, err := png.Load("../../data/in/"+dir+"/"+currImage.InPath, currImage.OutPath, dir, currImage.Effects)
			if err != nil {
				panic(err)
			}

			minX, minY, maxX, maxY := pngImg.Bounds.Min.X, pngImg.Bounds.Min.Y, pngImg.Bounds.Max.X, pngImg.Bounds.Max.Y
			bounds := image.Rect(minX, minY, maxX, maxY)
			for i, effect := range pngImg.Effects {
				if effect == "G" {
					pngImg.Grayscale(&bounds)
				} else if effect == "E" {
					pngImg.EdgeDetection(&bounds)
				} else if effect == "S" {
					pngImg.Sharpen(&bounds)
				} else if effect == "B" {
					pngImg.Blur(&bounds)
				} else {
					panic("No matching effect")
				}

				if i != len(pngImg.Effects)-1 {
					pngImg.SwapInAndOut()
				}
			}

			err = pngImg.Save("../../data/out/" + pngImg.Size + "_" + pngImg.OutPath)
			if err != nil {
				panic(err)
			}
		}
	}
}
