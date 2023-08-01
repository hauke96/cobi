package encoding

import (
	"cobi/util"
	"testing"
)

func Test_encodedArea_contains(t *testing.T) {
	area := EncodedArea{X: 3, Y: 2, W: 4, H: 2}

	util.AssertFalse(t, area.Contains(2, 1))
	util.AssertFalse(t, area.Contains(3, 1))
	util.AssertFalse(t, area.Contains(4, 1))
	util.AssertFalse(t, area.Contains(5, 1))
	util.AssertFalse(t, area.Contains(6, 1))
	util.AssertFalse(t, area.Contains(7, 1))

	util.AssertFalse(t, area.Contains(2, 2))
	util.AssertTrue(t, area.Contains(3, 2))
	util.AssertTrue(t, area.Contains(4, 2))
	util.AssertTrue(t, area.Contains(5, 2))
	util.AssertTrue(t, area.Contains(6, 2))
	util.AssertFalse(t, area.Contains(7, 2))

	util.AssertFalse(t, area.Contains(2, 3))
	util.AssertTrue(t, area.Contains(3, 3))
	util.AssertTrue(t, area.Contains(4, 3))
	util.AssertTrue(t, area.Contains(5, 3))
	util.AssertTrue(t, area.Contains(6, 3))
	util.AssertFalse(t, area.Contains(7, 3))

	util.AssertFalse(t, area.Contains(2, 4))
	util.AssertFalse(t, area.Contains(3, 4))
	util.AssertFalse(t, area.Contains(4, 4))
	util.AssertFalse(t, area.Contains(5, 4))
	util.AssertFalse(t, area.Contains(6, 4))
	util.AssertFalse(t, area.Contains(7, 4))
}

func Test_findMinUncoveredPixel_noAreas(t *testing.T) {
	x, y := FindMinUncoveredPixel([]EncodedArea{}, 10, 10)
	util.AssertEqual(t, 0, x)
	util.AssertEqual(t, 0, y)
}

func Test_findMinUncoveredPixel_oneArea(t *testing.T) {
	// 111.....
	// 111.....
	// 111.....
	// 111.....
	// 111.....
	// ........
	// ........
	// ........
	areas := []EncodedArea{
		{X: 0, Y: 0, W: 3, H: 5}, // 1
	}
	x, y := FindMinUncoveredPixel(areas, 8, 8)
	util.AssertEqual(t, 3, x)
	util.AssertEqual(t, 0, y)
}

func Test_findMinUncoveredPixel_completelyCoveredRows(t *testing.T) {
	// 11112223
	// 11112223
	// 1111...3
	// 1111....
	// ........
	// ........
	// ........
	// ........
	areas := []EncodedArea{
		{X: 0, Y: 0, W: 4, H: 5}, // 1
		{X: 4, Y: 0, W: 3, H: 2}, // 3
		{X: 7, Y: 0, W: 1, H: 3}, // 2
	}
	x, y := FindMinUncoveredPixel(areas, 8, 8)
	util.AssertEqual(t, 4, x)
	util.AssertEqual(t, 2, y)
}

func Test_findLargestNonEncodedArea(t *testing.T) {
	// 11112223
	// 11112223
	// 1111...3
	// 1111....
	// ........
	// ........
	// ........
	// ........
	areas := []EncodedArea{
		{X: 0, Y: 0, W: 4, H: 5}, // 1
		{X: 4, Y: 0, W: 3, H: 2}, // 3
		{X: 7, Y: 0, W: 1, H: 3}, // 2
	}
	values := util.TransposeArray([][]uint8{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 3, 4, 5, 6, 7},
	})

	newArea := *(findLargestNonEncodedArea(values, areas))

	util.AssertEqual(t, 4, newArea.X)
	util.AssertEqual(t, 2, newArea.Y)
	util.AssertEqual(t, 3, newArea.W)
	util.AssertEqual(t, 3, newArea.H)
	util.AssertEqual(t, [4]uint8{4, 6, 4, 6}, newArea.Values)
}
