// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image"
	"image/color"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *ImageTask) Grayscale(bounds *image.Rectangle) {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	// bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

// applyEffect will apply the convolution on the image with the given kernel
func applyEffect(img *ImageTask, kernel [3][3]float64, bounds *image.Rectangle) {
	outBounds := img.out.Bounds()
	var r, g, b, a uint32

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			accR, accG, accB, accA := 0.0, 0.0, 0.0, 0.0
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					X, Y := x+i-1, y+j-1

					// make sure X and Y are within the image boundaries
					if X >= outBounds.Min.X && X < outBounds.Max.X && Y >= outBounds.Min.Y && Y < outBounds.Max.Y {
						r, g, b, a = img.in.At(X, Y).RGBA()
						accR += float64(r) * kernel[i][j]
						accG += float64(g) * kernel[i][j]
						accB += float64(b) * kernel[i][j]

						if i == 1 && j == 1 {
							accA = float64(a)
						}
					}
				}
			}
			img.out.Set(x, y, color.RGBA64{clamp(accR), clamp(accG), clamp(accB), uint16(accA)})
		}
	}
}

// EdgeDetection applies the Edge Detection filter on the image
func (img *ImageTask) EdgeDetection(bounds *image.Rectangle) {
	applyEffect(img, [3][3]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}}, bounds)
}

// Sharpen applies the Sharpen filter on the image
func (img *ImageTask) Sharpen(bounds *image.Rectangle) {
	applyEffect(img, [3][3]float64{{0, -1, 0}, {-1, 5, -1}, {0, -1, 0}}, bounds)
}

// Blur applies the Blur filter on the image
func (img *ImageTask) Blur(bounds *image.Rectangle) {
	applyEffect(img, [3][3]float64{{1 / 9.0, 1 / 9.0, 1 / 9.0}, {1 / 9.0, 1 / 9.0, 1 / 9.0}, {1 / 9.0, 1 / 9.0, 1 / 9.0}}, bounds)
}
