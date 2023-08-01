package encoding

import (
	"cobi/image"
)

// EncodedArea represents
type EncodedArea struct {
	X, Y, W, H int
	Values     [4]byte
}

func (e *EncodedArea) Contains(x, y int) bool {
	return e.X <= x && x <= e.X+e.W-1 &&
		e.Y <= y && y <= e.Y+e.H-1
}

// Encode determines the encoded areas per color channel R (0), G (1), B (2) and A (3).
func Encode(img image.Image) ([4][]EncodedArea, error) {
	return [4][]EncodedArea{
		EncodeChannel(img.R),
		EncodeChannel(img.G),
		EncodeChannel(img.B),
		EncodeChannel(img.A),
	}, nil
}

func EncodeChannel(values [][]byte) []EncodedArea {
	var result []EncodedArea

	for {
		area := findLargestNonEncodedArea(values, result)
		if area == nil {
			break
		}
		result = append(result, *area)
	}

	return result
}

// findLargestNonEncodedArea finds the next encoded area following the strategy to find areas from the upper-left to the
// bottom-right of the image.
func findLargestNonEncodedArea(values [][]byte, areas []EncodedArea) *EncodedArea {
	width := len(values)
	height := len(values[0])

	minX, minY := findMinUncoveredPixel(areas, width, height)
	if minX == -1 || minY == -1 {
		return nil
	}

	// TODO use real interpolation methods to find encoded areas
	w := 8
	if minX+w >= width {
		w = width - minX
	}
	h := 8
	if minY+h >= height {
		h = height - minY
	}

	return &EncodedArea{
		X: minX,
		Y: minY,
		W: w,
		H: h,
		Values: [4]byte{
			values[minX][minY],
			values[minX+w-1][minY],
			values[minX][minY+h-1],
			values[minX+w-1][minY+h-1],
		},
	}
}

// findMinUncoveredPixel determines the smallest pixel that is not covered by any area. It is assumed that the encoded
// areas grow from the upper-left to the bottom-right. This means for example, when (3, 5) is the first non-covered
// pixel, all pixels in rows 0, 1 or 2 are covered and all pixels of row 3 in columns 0-4 are covered.
func findMinUncoveredPixel(areas []EncodedArea, width, height int) (int, int) {
	var x, y int
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			// Assume this pixel (x, y) is not covered and return (x, y) if it indeed isn't covered.
			isCovered := false
			for _, area := range areas {
				isCovered = isCovered || area.Contains(x, y)
				if isCovered {
					// Cell is covered -> Abort and process with next cell
					break
				}
			}

			if !isCovered {
				// Minimum non-covered cell found
				return x, y
			}
		}
	}

	// No pixel has been found that's not covered
	return -1, -1
}
