package stego

import (
	"errors"
	. "github.com/krezefal/cox-stegosystem/bmp"
	"image"
	"math"
	"math/rand"
	"time"
)

const roundingAccuracy = 1000

// GenerateMessage generates specified binary sequence.
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

// GenerateRandomMessage generates random binary sequence with specified length.
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

// PSNR calculates peak signal-to-noise ratio between one size arrays of data. It returns an error, when given arrays are
// completely the same.
func PSNR(data1, data2 [][]byte) (float64, error) {
	height := len(data1)
	width := len(data1[0])

	var sum float64
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			sum += math.Pow(float64(data1[i][j])-float64(data2[i][j]), 2)
		}
	}

	if sum == 0 {
		return 0, errors.New("unable to calculate PSNR: arrays of bytes are the same")
	}

	psnr := 10 * math.Log10((float64(height)*float64(width)*math.Pow(255, 2))/sum)
	return math.Floor(psnr*roundingAccuracy) / roundingAccuracy, nil
}

// PSNRimg calculates peak signal-to-noise ratio between one size images. It returns an error, when all pixels in given
// images are completely the same.
func PSNRimg(img1, img2 image.Image) (float64, error) {
	data1 := GetArrayRGB(img1)
	data2 := GetArrayRGB(img2)

	return PSNR(data1, data2)
}
