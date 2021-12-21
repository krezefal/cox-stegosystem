package bmp

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
)

const RED = "RED"
const GREEN = "GREEN"
const BLUE = "BLUE"

// ReadImage reads image at specified path. For bmp only.
func ReadImage(filePath string) (image.Image, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return bmp.Decode(imgFile)
}

// WriteImage writes image to specified path. For bmp only.
func WriteImage(filePath string, img image.Image) error {
	if file, err := os.Create(filePath); err != nil {
		return err
	} else {
		if imgErr := bmp.Encode(file, img); imgErr != nil {
			return imgErr
		}
	}

	return nil
}

// GetComponent gets specified component (RED, GREEN or BLUE) from image.
func GetComponent(img image.Image, component string) [][]byte {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	singleComponentArray := make([][]byte, height)
	for i := range singleComponentArray {
		singleComponentArray[i] = make([]byte, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			switch component {
			case RED:
				singleComponentArray[y][x] = uint8(r)
			case GREEN:
				singleComponentArray[y][x] = uint8(g)
			case BLUE:
				singleComponentArray[y][x] = uint8(b)
			}
		}
	}

	return singleComponentArray
}

// GetArrayRGB gets all component (RED, GREEN and BLUE) from image and returns it as single 2d array where the pixel
// components alternately follow each other in the following order: r, g, b.
func GetArrayRGB(img image.Image) [][]byte {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	arrayRGB := make([][]byte, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			pix := img.At(x, y)
			r, g, b, _ := pix.RGBA()

			arrayRGB[y] = append(arrayRGB[y], byte(r), byte(g), byte(b))
		}
	}

	return arrayRGB
}

// SetComponent sets array of components into image (replace values of specified component (RED, GREEN or BLUE) in source
// image with provided values).
func SetComponent(img image.Image, singleComponentArray [][]byte, component string) *image.RGBA {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	imgRGBA := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			pix := img.At(x, y)
			r, g, b, a := pix.RGBA()

			newPix := color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			}

			switch component {
			case RED:
				newPix.R = singleComponentArray[y][x]
			case GREEN:
				newPix.G = singleComponentArray[y][x]
			case BLUE:
				newPix.B = singleComponentArray[y][x]
			}

			imgRGBA.Set(x, y, newPix)
		}
	}

	return imgRGBA
}
