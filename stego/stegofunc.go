package stego

import (
	"errors"
	. "github.com/krezefal/cox-stegosystem/dct"
	"image"
	"math"
)

func EmbedMessage(img image.Image, message []byte, alpha float64) (image.Image, int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	singleComponentArray := getComponent(img, BLUE)
	dct := DCT(singleComponentArray)

	var cnt int
	for y := 0; y < height; y += DataUnitSize {
		for x := 0; x < width; x += DataUnitSize {
			if cnt < len(message) {
				embedMesBitInUnitMaxAC(dct, x, y, message[cnt], alpha)
				cnt++
			} else {
				break
			}
		}
		if cnt == len(message) {
			break
		}
	}

	idct := IDCT(dct)
	newImg := setComponent(img, idct, BLUE)

	return newImg, cnt
}

func ExtractMessage(srcImg image.Image, tgImg image.Image) ([]byte, error) {
	boundsSrc := srcImg.Bounds()
	boundsTg := tgImg.Bounds()

	widthSrc, heightSrc := boundsSrc.Max.X, boundsSrc.Max.Y
	widthTg, heightTg := boundsTg.Max.X, boundsTg.Max.Y

	if widthSrc != widthTg || heightSrc != heightTg {
		return nil, errors.New("image sizes do not match")
	}

	singleComponentArraySrc := getComponent(srcImg, BLUE)
	dctSrc := DCT(singleComponentArraySrc)

	singleComponentArrayTg := getComponent(tgImg, BLUE)
	dctTg := DCT(singleComponentArrayTg)

	var message []byte
	for y := 0; y < heightSrc; y += DataUnitSize {
		for x := 0; x < widthSrc; x += DataUnitSize {
			detectedMessageBit := compareUnitsMaxAC(dctSrc, dctTg, x, y)
			if detectedMessageBit != -1 {
				message = append(message, byte(detectedMessageBit))
			}
		}
	}

	return message, nil
}

func embedMesBitInUnitMaxAC(dct [][]float64, startX, startY int, messageBit byte, alpha float64) {

	endX := startX + DataUnitSize
	endY := startY + DataUnitSize

	skipDC := 1
	var maxX, maxY int
	maxAC := dct[startY][startX+skipDC]

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if x%DataUnitSize == 0 && y%DataUnitSize == 0 {
				// This is DC coefficient
				continue
			}
			if math.Abs(maxAC) <= math.Abs(dct[y][x]) {
				maxX = x
				maxY = y
				maxAC = dct[maxY][maxX]
			}
		}
	}

	if messageBit == 0 {
		dct[maxY][maxX] += alpha
	} else {
		dct[maxY][maxX] -= alpha
	}
}

func compareUnitsMaxAC(dctSrc [][]float64, dctTg [][]float64, startX, startY int) int {
	endX := startX + DataUnitSize
	endY := startY + DataUnitSize

	skipDC := 1
	var maxX, maxY int
	maxAC := dctSrc[startY][startX+skipDC]

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if x%DataUnitSize == 0 && y%DataUnitSize == 0 {
				// This is DC coefficient
				continue
			}
			if math.Abs(maxAC) <= math.Abs(dctSrc[y][x]) {
				maxX = x
				maxY = y
				maxAC = dctSrc[maxY][maxX]
			}
		}
	}

	switch {
	case dctSrc[maxY][maxX] > dctTg[maxY][maxX]:
		return 1
	case dctSrc[maxY][maxX] < dctTg[maxY][maxX]:
		return 0
	}

	return -1
}
