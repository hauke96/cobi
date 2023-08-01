package decoding

import (
	"cobi/encoding"
	"cobi/image"
	"fmt"
	"github.com/pkg/errors"
)

func Decode(areas [4][]encoding.EncodedArea) (*image.Image, error) {
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

func getAndEnsureWidthHeight(areas [4][]encoding.EncodedArea) (int, int, error) {
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

func getSizeOfChannel(areas []encoding.EncodedArea) (int, int, error) {
	width := 0
	height := 0

	for _, area := range areas {
		if width < area.X+area.W {
			width = area.X + area.W
		}
		if height < area.Y+area.H {
			height = area.Y + area.H
		}
	}

	x, y := encoding.FindMinUncoveredPixel(areas, width, height)
	if x != -1 || y != -1 {
		return -1, -1, errors.New("Encoded areas do not cover the whole image")
	}

	return width, height, nil
}

func interpolateChannel(areas []encoding.EncodedArea, width, height int) [][]byte {
	result := make([][]byte, width)

	for x := 0; x < width; x++ {
		result[x] = make([]byte, height)
	}

	for _, area := range areas {
		interpolatedValues := area.GetInterpolatedArea()

		for x := 0; x < area.W; x++ {
			for y := 0; y < area.H; y++ {
				result[x+area.X][y+area.Y] = interpolatedValues[x][y]
			}
		}
	}

	return result
}
