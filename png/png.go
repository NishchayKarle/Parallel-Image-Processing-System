// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// The ImageTask represents a structure for working with PNG images.
// From Professor Samuels: You are allowed to update this and change it as you wish!
type ImageTask struct {
	in      *image.RGBA64   //The original pixels before applying the effect
	out     *image.RGBA64   //The updated pixels after applying teh effect
	Bounds  image.Rectangle //The size of the image
	Size    string          // Type of image - small, mixture, big
	OutPath string          // Path to save the image
	Effects []string        // Effects to be applied on the image
}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
// From Professor Samuels:  You are allowed to modify and update this as you wish
func Load(fileInputPath string, fileOutputPath string, size string, effects []string) (*ImageTask, error) {

	inReader, err := os.Open(fileInputPath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inOrig, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	bounds := inOrig.Bounds()

	outImg := image.NewRGBA64(bounds)
	inImg := image.NewRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := inOrig.At(x, y).RGBA()
			inImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	task := &ImageTask{}
	task.in = inImg
	task.out = outImg
	task.Bounds = bounds
	task.OutPath = fileOutputPath
	task.Size = size
	task.Effects = make([]string, 0)
	task.Effects = append(task.Effects, effects...)
	return task, nil
}

// Save saves the image to the given file
// From Professor Samuels:  You are allowed to modify and update this as you wish
func (img *ImageTask) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.out)
	if err != nil {
		return err
	}
	return nil
}

// clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}

// SwapInAndOut will swap the in and out image pointers
func (img *ImageTask) SwapInAndOut() {
	img.in, img.out = img.out, img.in
}
