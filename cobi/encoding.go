package main

import (
	"cobi/image"
)

// EncodedArea represents
type EncodedArea struct {
	X, Y, W, H int
	Values     [4]byte
}

func (e *EncodedArea) Contains(x, y int) bool {
	return e.X <= x && x <= e.X+e.W &&
		e.Y <= y && y <= e.Y+e.H
}

// Encode determines the encoded areas per color channel R (0), G (1), B (2) and A (3).
func Encode(img image.Image) ([4][]EncodedArea, error) {
	for y := 0; y < img.Width; y++ {
		for x := 0; x < img.Height; x++ {
		}
	}

	return [4][]EncodedArea{}, nil
}

func EncodeChannel(values [][]byte) []EncodedArea {
	var result []EncodedArea

	return result
}

// getMinUncoveredPixel determines the smallest pixel that is not covered by any area. It is assumed that the encoded
// areas grow from the upper-left to the bottom-right. This means for example, when (3, 5) is the first non-covered
// pixel, there is no other pixel p with p.X <= 3 && p.Y <= 5.
func getMinUncoveredPixel(areas []EncodedArea, width, height int) (int, int) {
	x := 0
	y := 0
	for ; y < height; y++ {
		for ; x < height; x++ {
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
