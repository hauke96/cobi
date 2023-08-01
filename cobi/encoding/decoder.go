package encoding

import (
	"cobi/image"
	"fmt"
	"github.com/pkg/errors"
)

func Decode(areas [4][]EncodedArea) (*image.Image, error) {
	width, height, err := getAndEnsureWidthHeight(areas)
	if err != nil {
		return nil, err
	}

	img := image.New(width, height)
	img.R = interpolateChannel(areas[0], width, height)
	img.G = interpolateChannel(areas[1], width, height)
	img.B = interpolateChannel(areas[2], width, height)
	img.A = interpolateChannel(areas[3], width, height)
	return img, nil
}

func getAndEnsureWidthHeight(areas [4][]EncodedArea) (int, int, error) {
	widthR, heightR, err := getSizeOfChannel(areas[0])
	if err != nil {
		return -1, -1, err
	}
	widthG, heightG, err := getSizeOfChannel(areas[1])
	if err != nil {
		return -1, -1, err
	}
	widthB, heightB, err := getSizeOfChannel(areas[2])
	if err != nil {
		return -1, -1, err
	}
	widthA, heightA, err := getSizeOfChannel(areas[3])
	if err != nil {
		return -1, -1, err
	}

	if widthR != widthG || widthG != widthB || widthB != widthA {
		return -1, -1, errors.New(fmt.Sprintf("Width not equal in all channels (R, G, B, A): %d, %d, %d, %d", widthR, widthG, widthB, widthA))
	}
	if heightR != heightG || heightG != heightB || heightB != heightA {
		return -1, -1, errors.New(fmt.Sprintf("Height not equal in all channels (R, G, B, A): %d, %d, %d, %d", heightR, heightG, heightB, heightA))
	}

	return widthR, heightR, nil
}

func getSizeOfChannel(areas []EncodedArea) (int, int, error) {
	width := 0
	height := 0

	for _, area := range areas {
		areaWidth := int(area.W)
		areaHeight := int(area.H)

		if width < area.X+areaWidth {
			width = area.X + areaWidth
		}
		if height < area.Y+areaHeight {
			height = area.Y + areaHeight
		}
	}

	//x, y := encoding.findMinUncoveredPixel(areas, width, height)
	//if x != -1 || y != -1 {
	//	return -1, -1, errors.New("Encoded areas do not cover the whole image")
	//}

	return width, height, nil
}

func interpolateChannel(areas []EncodedArea, width, height int) [][]uint8 {
	result := make([][]uint8, width)

	for x := 0; x < width; x++ {
		result[x] = make([]uint8, height)
	}

	for _, area := range areas {
		interpolatedValues := area.GetInterpolatedArea()

		for x := 0; x < int(area.W); x++ {
			for y := 0; y < int(area.H); y++ {
				result[x+area.X][y+area.Y] = interpolatedValues[x][y]
			}
		}
	}

	return result
}
