package dct

import (
	"math"
)

const DataUnitSize = 8

// DCT performs discrete-cosine-transform on array of bytes; returns calculated values in frequency range.
func DCT(data [][]byte) [][]float64 {

	height := len(data)
	width := len(data[0])
	result := make([][]float64, height)
	for i := range result {
		result[i] = make([]float64, width)
	}

	for h := 0; h < height; h += DataUnitSize {
		for w := 0; w < width; w += DataUnitSize {
			for k := 0; k < DataUnitSize; k++ {
				var Ck float64
				if k == 0 {
					Ck = 1 / float64(DataUnitSize)
				} else {
					Ck = 2 / float64(DataUnitSize)
				}

				for l := 0; l < DataUnitSize; l++ {
					var Cl float64
					if l == 0 {
						Cl = 1 / float64(DataUnitSize)
					} else {
						Cl = 2 / float64(DataUnitSize)
					}

					var sum float64
					for i := 0; i < DataUnitSize; i++ {
						cos1 := math.Cos(((float64(2*i + 1)) * math.Pi * float64(k)) / float64(2*DataUnitSize))
						for j := 0; j < DataUnitSize; j++ {
							cos2 := math.Cos((float64(2*j+1) * math.Pi * float64(l)) / float64(2*DataUnitSize))
							sum += float64(data[h+i][w+j]) * cos1 * cos2
						}
					}
					val := math.Sqrt(Ck) * math.Sqrt(Cl) * sum
					result[h+k][w+l] = val
				}
			}
		}
	}

	return result
}

// IDCT performs inverse discrete-cosine-transform on array of values in frequency range.
func IDCT(data [][]float64) [][]byte {

	height := len(data)
	width := len(data[0])
	result := make([][]byte, height)
	for i := range result {
		result[i] = make([]byte, width)
	}

	for h := 0; h < height; h += DataUnitSize {
		for w := 0; w < height; w += DataUnitSize {
			for i := 0; i < DataUnitSize; i++ {
				for j := 0; j < DataUnitSize; j++ {
					var sum float64
					for k := 0; k < DataUnitSize; k++ {
						var Ck float64
						if k == 0 {
							Ck = 1 / float64(DataUnitSize)
						} else {
							Ck = 2 / float64(DataUnitSize)
						}

						cos1 := math.Cos((float64(2*i+1) * math.Pi * float64(k)) / float64(2*DataUnitSize))
						for l := 0; l < DataUnitSize; l++ {
							var Cl float64
							if l == 0 {
								Cl = 1 / float64(DataUnitSize)
							} else {
								Cl = 2 / float64(DataUnitSize)
							}

							cos2 := math.Cos((float64(2*j+1) * math.Pi * float64(l)) / float64(2*DataUnitSize))
							sum += math.Sqrt(Ck) * math.Sqrt(Cl) * data[h+k][w+l] * cos1 * cos2
						}
					}
					result[h+i][w+j] = byte(sum)
				}
			}
		}
	}

	return result
}
