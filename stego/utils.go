package stego

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"
)

const RED = "RED"
const GREEN = "GREEN"
const BLUE = "BLUE"

func ReadImage(filePath string) (image.Image, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return bmp.Decode(imgFile)
}

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

func getComponent(img image.Image, component string) [][]byte {

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

func setComponent(img image.Image, singleComponentArray [][]byte, component string) *image.RGBA {

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

// GenerateMessage func generates specified binary sequence
func GenerateMessage(sequence string) []byte {
	message := make([]byte, len(sequence))

	for i, bit := range sequence {
		if bit == '1' {
			message[i] = 0x01
		} else {
			message[i] = 0x00
		}
	}

	return message
}

// GenerateRandomMessage func generates random binary sequence with specified length
func GenerateRandomMessage(length int) []byte {
	rMessage := make([]byte, length)

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	for i := 0; i < cap(rMessage); i++ {
		if r.Int()%2 == 1 {
			rMessage[i] = 0x00
		} else {
			rMessage[i] = 0x01
		}
	}

	return rMessage
}
