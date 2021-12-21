package stego

import (
	"encoding/binary"
	"errors"
	. "github.com/krezefal/cox-stegosystem/bmp"
	. "github.com/krezefal/cox-stegosystem/dct"
	"image"
	"math"
)

// EmbedMessage embeds given message into image-container so this message becomes hidden inside the image. It also returns
// the number of bits which fit inside the container.
func EmbedMessage(img image.Image, message []byte, alpha float64) (image.Image, int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	singleComponentArray := GetComponent(img, BLUE)
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
	newImg := SetComponent(img, idct, BLUE)

	return newImg, cnt
}

// ExtractMessage extracts hidden message from image-container and returns it. It returns an error, when the sizes of
// given source and target images are not the same.
func ExtractMessage(srcImg image.Image, tgImg image.Image, messageLen int) ([]byte, error) {
	boundsSrc := srcImg.Bounds()
	boundsTg := tgImg.Bounds()

	widthSrc, heightSrc := boundsSrc.Max.X, boundsSrc.Max.Y
	widthTg, heightTg := boundsTg.Max.X, boundsTg.Max.Y

	if widthSrc != widthTg || heightSrc != heightTg {
		return nil, errors.New("unable to extract message: image sizes do not match")
	}

	singleComponentArraySrc := GetComponent(srcImg, BLUE)
	dctSrc := DCT(singleComponentArraySrc)

	singleComponentArrayTg := GetComponent(tgImg, BLUE)
	dctTg := DCT(singleComponentArrayTg)

	var message []byte
	for y := 0; y < heightSrc; y += DataUnitSize {
		for x := 0; x < widthSrc; x += DataUnitSize {
			if len(message) < messageLen {
				detectedMessageBit := compareUnitsMaxAC(dctSrc, dctTg, x, y)
				if detectedMessageBit != -1 {
					message = append(message, byte(detectedMessageBit))
				}
			} else {
				break
			}
		}
		if len(message) == messageLen {
			break
		}
	}

	return message, nil
}

// embedMesBitInUnitMaxAC finds max AC coefficient in current image data unit and embeds there current message bit by
// adding or subtracting specified alpha.
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

// compareUnitsMaxAC compares two current data units of source and target images by calculating the difference between
// max AC coefficient in source data unit and the AC coefficient in target data unit taken on the same position. If
// calculated difference is positive, it decodes "1", otherwise "0". In case of equality, it returns -1 code, that means
// that the current data unit has been marked as exactly the same for source and target images.
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

// EmbedMessageLen embeds message length into the reserved section of the file header (BITMAPINFOHEADER).
func EmbedMessageLen(file []byte, messageLen int) []byte {
	ml := make([]byte, 4)
	binary.LittleEndian.PutUint32(ml, uint32(messageLen))

	offset := 6
	for i := range ml {
		file[offset+i] = ml[i]
	}

	return file
}

// DetectEmbedding checks the reserved section of BITMAPINFOHEADER on presence of message length embedded before.
// Returns the detection flag.
func DetectEmbedding(file []byte) bool {
	offset := 6
	if file[offset] == 0 && file[offset+1] == 0 && file[offset+2] == 0 && file[offset+3] == 0 {
		return false
	}
	return true
}

// ExtractMessageLen extracts hidden message length from the reserved section of BITMAPINFOHEADER.
func ExtractMessageLen(file []byte) int {
	offset := 6
	ml := file[offset : offset+4]
	messageLen := binary.LittleEndian.Uint32(ml)

	return int(messageLen)
}
