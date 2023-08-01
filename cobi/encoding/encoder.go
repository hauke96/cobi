package encoding

import (
	"cobi/image"
	"cobi/interpolate"
	"github.com/hauke96/sigolo"
	"math"
)

// EncodedArea represents
type EncodedArea struct {
	X, Y   int
	W, H   uint8
	Values [4]uint8
}

func (e *EncodedArea) Contains(x, y int) bool {
	return e.X <= x && x <= e.X+int(e.W)-1 &&
		e.Y <= y && y <= e.Y+int(e.H)-1
}

func (e *EncodedArea) GetInterpolatedArea() [][]uint8 {
	return interpolate.Interpolate(e.W, e.H, e.Values)
}

func GetDebugImage(width, height int, areas [4][]EncodedArea) *image.Image {
	img := image.New(width, height)

	for x := 0; x < width; x++ {
		img.R[x] = make([]uint8, height)
		img.G[x] = make([]uint8, height)
		img.B[x] = make([]uint8, height)
		img.A[x] = make([]uint8, height)
	}

	for i, _ := range areas {
		var channel [][]uint8
		switch i {
		case 0:
			channel = img.R
		case 1:
			channel = img.G
		case 2:
			channel = img.B
		case 3:
			channel = img.A
		}

		for _, area := range areas[i] {
			channel[area.X][area.Y] = 255
			channel[area.X+int(area.W)-1][area.Y] = 255
			channel[area.X][area.Y+int(area.H)-1] = 255
			channel[area.X+int(area.W)-1][area.Y+int(area.H)-1] = 255
		}
	}

	for x, _ := range img.A {
		for y, _ := range img.A[x] {
			img.A[x][y] = 255
		}
	}

	return img
}

type ChannelEncoder struct {
	coveredPixel       [][]bool
	minUncoveredPixelX int
	minUncoveredPixelY int
	imageWidth         int
	imageHeight        int
	channel            [][]uint8
}

func newChannelEncoder(width, height int, channel [][]uint8) *ChannelEncoder {
	coveredPixel := make([][]bool, width)
	for x := 0; x < width; x++ {
		coveredPixel[x] = make([]bool, height)
	}

	return &ChannelEncoder{
		coveredPixel:       coveredPixel,
		minUncoveredPixelX: 0,
		minUncoveredPixelY: 0,
		imageWidth:         width,
		imageHeight:        height,
		channel:            channel,
	}
}

// Encode determines the encoded areas per color channel R (0), G (1), B (2) and A (3).
func Encode(img image.Image) ([4][]EncodedArea, error) {
	sigolo.Debug("Encode channel R")
	channelR := newChannelEncoder(img.Width, img.Height, img.R).encodeChannel(img.R)

	sigolo.Debug("Encode channel G")
	channelG := newChannelEncoder(img.Width, img.Height, img.G).encodeChannel(img.G)

	sigolo.Debug("Encode channel B")
	channelB := newChannelEncoder(img.Width, img.Height, img.B).encodeChannel(img.B)

	sigolo.Debug("Encode channel A")
	channelA := newChannelEncoder(img.Width, img.Height, img.A).encodeChannel(img.A)

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
	areaX, areaY := e.minUncoveredPixelX, e.minUncoveredPixelY
	if areaX == -1 || areaY == -1 {
		return nil
	}

	areaWidth, areaHeight := e.getAreaSize(areaX, areaY)

	encodedArea := &EncodedArea{
		X: areaX,
		Y: areaY,
		W: areaWidth,
		H: areaHeight,
		Values: [4]uint8{
			values[areaX][areaY],
			values[areaX+int(areaWidth)-1][areaY],
			values[areaX][areaY+int(areaHeight)-1],
			values[areaX+int(areaWidth)-1][areaY+int(areaHeight)-1],
		},
	}

	e.addToCoverageMap(*encodedArea)

	return encodedArea
}

func (e *ChannelEncoder) addToCoverageMap(encodedArea EncodedArea) {
	for y := encodedArea.Y; y < encodedArea.Y+int(encodedArea.H); y++ {
		for x := encodedArea.X; x < encodedArea.X+int(encodedArea.W); x++ {
			e.coveredPixel[x][y] = true
		}
	}
	e.minUncoveredPixelX, e.minUncoveredPixelY = e.findMinUncoveredPixel()
}

// findMinUncoveredPixel determines the smallest pixel that is not covered by any area. It is assumed that the encoded
// areas grow from the upper-left to the bottom-right. This means for example, when (3, 5) is the first non-covered
// pixel, all pixels in rows 0, 1 or 2 are covered and all pixels of row 3 in columns 0-4 are covered.
func (e *ChannelEncoder) findMinUncoveredPixel() (int, int) {
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

func (e *ChannelEncoder) getAreaSize(x, y int) (uint8, uint8) {
	maxWidth := uint8(0)
	for _x := x; _x < e.imageWidth && y+int(maxWidth) < e.imageHeight && maxWidth < math.MaxUint8; _x++ {
		if e.coveredPixel[_x][y] {
			break
		}
		maxWidth++
	}

	width := uint8(2)
	if maxWidth < width {
		width = maxWidth
	}

	// TODO make this configurable
	differenceThreshold := 0.0005

	for ; width < maxWidth && y+int(width) < e.imageHeight && width < math.MaxUint8; width++ {
		// Build a square -> width = height
		difference := e.calculateDifference(x, y, x+int(width), y+int(width))
		if difference > differenceThreshold {
			// We passed the point of tolerable quality -> Previous iteration was the last one with an okay difference.
			width--
			break
		}
	}

	return width, width
}

func (e *ChannelEncoder) calculateDifference(x1, y1, x2, y2 int) float64 {
	interpolatedData := interpolate.Interpolate(uint8(x2-x1), uint8(y2-y1), [4]uint8{
		e.channel[x1][y1],
		e.channel[x2][y1],
		e.channel[x1][y2],
		e.channel[x2][y2],
	})

	// Calculating the sum square difference as one distance measurement between the two image sections
	// See https://datascience.stackexchange.com/questions/48642/how-to-measure-the-similarity-between-two-images

	squaredSum := 0.0
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			squaredSum += math.Pow(float64(e.channel[x][y])-float64(interpolatedData[x-x1][y-y1]), 2)
		}
	}

	squaredSumOriginal := 0.0
	squaredSumInterpolated := 0.0
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			squaredSumOriginal += math.Pow(float64(e.channel[x][y]), 2)
			squaredSumInterpolated += math.Pow(float64(interpolatedData[x-x1][y-y1]), 2)
		}
	}

	normalizedDifference := squaredSum / math.Sqrt(squaredSumOriginal*squaredSumInterpolated)
	return normalizedDifference
}
