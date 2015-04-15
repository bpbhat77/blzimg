package main

import "image"
import "image/draw"
import "image/color"
import "errors"
import "os"
import "fmt"

type LightestOperation struct{}

func (c LightestOperation) lightest(colors []color.RGBA) color.Color {
	lightest := color.RGBA{0, 0, 0, 0}

	for _, color := range colors {
		if luminance(color) > luminance(lightest) {
			lightest = color
		}
	}

	return lightest
}

func (c LightestOperation) ResultFiles(files []string) (image.Image, error) {
	firstFile, _ := os.Open(files[0])
	defer firstFile.Close()
	firstImage, _, _ := image.Decode(firstFile)
	bounds := firstImage.Bounds()
	lightest := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(lightest, bounds, firstImage, bounds.Min, draw.Src)

	for _, currentImageFile := range files {
		fmt.Printf("Processing %s...\n", currentImageFile)

		currentFile, _ := os.Open(currentImageFile)
		defer currentFile.Close()
		currentImage, _, _ := image.Decode(currentFile)

		if currentImage.Bounds() != bounds {
			return nil, errors.New("The images have different size!")
		}

		imageToCompare := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(imageToCompare, bounds, currentImage, bounds.Min, draw.Src)
		c.getlightestImageBetweenTwo(lightest, imageToCompare)
	}

	return lightest, nil
}

func (c LightestOperation) Result(images []image.Image) (image.Image, error) {
	firstImage := images[0]
	bounds := firstImage.Bounds()
	lightest := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(lightest, bounds, firstImage, bounds.Min, draw.Src)

	for _, currentImage := range images {
		if currentImage.Bounds() != bounds {
			return nil, errors.New("The images have different size!")
		}

		imageToCompare := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(imageToCompare, bounds, currentImage, bounds.Min, draw.Src)
		c.getLightestImageBetweenTwo(lightest, imageToCompare)
	}

	return lightest, nil
}

func (c LightestOperation) getLightestImageBetweenTwo(current, other *image.RGBA) {
	for i := current.Bounds().Min.X; i < current.Bounds().Max.X; i++ {
		for j := current.Bounds().Min.Y; j < current.Bounds().Max.Y; j++ {
			currentLightestImagePixel := current.At(i, j).(color.RGBA)
			otherImagePixel := other.At(i, j).(color.RGBA)

			lightestColor := c.lightest([]color.RGBA{currentLightestImagePixel, otherImagePixel})
			current.Set(i, j, lightestColor)
		}
	}
}

// http://stackoverflow.com/questions/596216/formula-to-determine-brightness-of-rgb-color
func luminance(someColor color.Color) uint32 {
	r, g, b, _ := someColor.RGBA()
	return uint32(0.2126*float32(r) + 0.7152*float32(g) + 0.0722*float32(b))
}

func average(someColor color.Color) uint32 {
	r, g, b, _ := someColor.RGBA()
	return (r + g + b) / 3
}