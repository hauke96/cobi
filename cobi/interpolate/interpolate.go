package interpolate

// Interpolate returns the interpolated area between the four given values. The values represent the following corners:
// [0] - upper left
// [1] - upper right
// [2] - bottom left
// [3] - bottom right
func Interpolate(w, h int, v [4]byte) [][]byte {
	result := make([][]byte, w)
	for x := 0; x < w; x++ {
		result[x] = make([]byte, h)
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
	for x := 0; x < w; x++ {
		interpolateColumn(result, x)
	}

	return result
}

func interpolateRow(v [][]byte, y int) {
	w := len(v)
	minXValue := v[0][y]
	maxXValue := v[w-1][y]
	increasePerColumn := float32(maxXValue-minXValue) / float32(w-1)
	for x := 1; x < w-1; x++ {
		v[x][y] = byte(float32(minXValue) + float32(x)*increasePerColumn)
	}
}

func interpolateColumn(v [][]byte, x int) {
	h := len(v[x])
	minYValue := v[x][0]
	maxYValue := v[x][h-1]
	increasePerRow := float32(maxYValue-minYValue) / float32(h-1)
	for y := 1; y < h-1; y++ {
		v[x][y] = byte(float32(minYValue) + float32(y)*increasePerRow)
	}
}
