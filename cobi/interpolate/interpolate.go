package interpolate

// Interpolate returns the interpolated area between the four given values. The values represent the following corners:
// [0] - upper left
// [1] - upper right
// [2] - bottom left
// [3] - bottom right
func Interpolate(w, h uint8, v [4]uint8) [][]uint8 {
	result := make([][]uint8, w)
	for x := uint8(0); x < w; x++ {
		result[x] = make([]uint8, h)
	}

	// Fill corner, which are already known by the four given values
	result[0][0] = v[0]
	result[w-1][0] = v[1]
	result[0][h-1] = v[2]
	result[w-1][h-1] = v[3]

	// First, interpolate first and last row
	interpolateRow(result, 0)
	interpolateRow(result, h-1)

	// Then interpolate all columns between first and last row
	for x := uint8(0); x < w; x++ {
		interpolateColumn(result, x)
	}

	return result
}

func interpolateRow(v [][]uint8, y uint8) {
	w := len(v)
	leftXValue := float32(v[0][y])
	rightXValue := float32(v[w-1][y])
	increasePerColumn := (rightXValue - leftXValue) / float32(w-1)
	for x := 1; x < w-1; x++ {
		v[x][y] = uint8(leftXValue + float32(x)*increasePerColumn)
	}
}

func interpolateColumn(v [][]uint8, x uint8) {
	h := len(v[x])
	upperYValue := float32(v[x][0])
	lowerYValue := float32(v[x][h-1])
	increasePerRow := (lowerYValue - upperYValue) / float32(h-1)
	for y := 1; y < h-1; y++ {
		v[x][y] = uint8(upperYValue + float32(y)*increasePerRow)
	}
}
