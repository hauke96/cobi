package encoding

import (
	"cobi/image"
	"cobi/interpolate"
	"github.com/hauke96/sigolo"
)

// EncodedArea represents
type EncodedArea struct {
	X, Y, W, H int
	Values     [4]uint8
}

func (e *EncodedArea) Contains(x, y int) bool {
	return e.X <= x && x <= e.X+e.W-1 &&
		e.Y <= y && y <= e.Y+e.H-1
}

func (e *EncodedArea) GetInterpolatedArea() [][]uint8 {
	return interpolate.Interpolate(e.W, e.H, e.Values)
}

type ChannelEncoder struct {
	coveredPixel       [][]bool
	minUncoveredPixelX int
	minUncoveredPixelY int
	imageWidth         int
	imageHeight        int
}

func newChannelEncoder(width, height int) *ChannelEncoder {
	coveredPixel := make([][]bool, width)
	for x := 0; x < width; x++ {
		coveredPixel[x] = make([]bool, height)
	}

	return &ChannelEncoder{
		minUncoveredPixelX: 0,
		minUncoveredPixelY: 0,
		imageWidth:         width,
		imageHeight:        height,
		coveredPixel:       coveredPixel,
	}
}

// Encode determines the encoded areas per color channel R (0), G (1), B (2) and A (3).
func Encode(img image.Image) ([4][]EncodedArea, error) {
	sigolo.Debug("Encode channel R")
	channelR := newChannelEncoder(img.Width, img.Height).encodeChannel(img.R)

	sigolo.Debug("Encode channel G")
	channelG := newChannelEncoder(img.Width, img.Height).encodeChannel(img.G)

	sigolo.Debug("Encode channel B")
	channelB := newChannelEncoder(img.Width, img.Height).encodeChannel(img.B)
	
	sigolo.Debug("Encode channel A")
	channelA := newChannelEncoder(img.Width, img.Height).encodeChannel(img.A)

	return [4][]EncodedArea{
		channelR,
		channelG,
		channelB,
		channelA,
	}, nil
}

func (e *ChannelEncoder) encodeChannel(values [][]uint8) []EncodedArea {
	var result []EncodedArea

	for {
		area := e.findLargestNonEncodedArea(values, result)
		if area == nil {
			break
		}
		result = append(result, *area)
	}

	return result
}

// findLargestNonEncodedArea finds the next encoded area following the strategy to find areas from the upper-left to the
// bottom-right of the image.
func (e *ChannelEncoder) findLargestNonEncodedArea(values [][]uint8, areas []EncodedArea) *EncodedArea {
	imgWidth := len(values)
	imgHeight := len(values[0])

	//areaX, areaY := e.FindMinUncoveredPixel(imgWidth, imgHeight)
	areaX, areaY := e.minUncoveredPixelX, e.minUncoveredPixelY
	if areaX == -1 || areaY == -1 {
		return nil
	}

	// TODO use real interpolation methods to find encoded areas
	areaWidth := 3
	if areaX+areaWidth >= imgWidth {
		areaWidth = imgWidth - areaX
	}
	areaHeight := 3
	if areaY+areaHeight >= imgHeight {
		areaHeight = imgHeight - areaY
	}

	encodedArea := &EncodedArea{
		X: areaX,
		Y: areaY,
		W: areaWidth,
		H: areaHeight,
		Values: [4]uint8{
			values[areaX][areaY],
			values[areaX+areaWidth-1][areaY],
			values[areaX][areaY+areaHeight-1],
			values[areaX+areaWidth-1][areaY+areaHeight-1],
		},
	}

	e.addToCoverageMap(*encodedArea)

	return encodedArea
}

func (e *ChannelEncoder) addToCoverageMap(encodedArea EncodedArea) {
	for y := encodedArea.Y; y < encodedArea.Y+encodedArea.H; y++ {
		for x := encodedArea.X; x < encodedArea.X+encodedArea.W; x++ {
			e.coveredPixel[x][y] = true
		}
	}
	e.minUncoveredPixelX, e.minUncoveredPixelY = e.FindMinUncoveredPixel()
}

// FindMinUncoveredPixel determines the smallest pixel that is not covered by any area. It is assumed that the encoded
// areas grow from the upper-left to the bottom-right. This means for example, when (3, 5) is the first non-covered
// pixel, all pixels in rows 0, 1 or 2 are covered and all pixels of row 3 in columns 0-4 are covered.
func (e *ChannelEncoder) FindMinUncoveredPixel() (int, int) {
	for y := e.minUncoveredPixelY; y < e.imageHeight; y++ {
		for x := 0; x < e.imageWidth; x++ {
			if !e.coveredPixel[x][y] {
				return x, y
			}
		}
	}

	// No pixel has been found that's not covered
	return -1, -1
}
