package encoding

import (
	"cobi/util"
	"testing"
)

func Test_decode(t *testing.T) {
	// 11111222
	// 11111222
	// 11111222
	// 11111333
	// 11111333
	areas := []EncodedArea{
		{X: 0, Y: 0, W: 5, H: 5, Values: [4]uint8{0, 10, 5, 20}},   // 1
		{X: 5, Y: 0, W: 3, H: 3, Values: [4]uint8{12, 14, 18, 10}}, // 2
		{X: 5, Y: 3, W: 3, H: 2, Values: [4]uint8{19, 19, 20, 20}}, // 3
	}

	// TransposeArray needed because the image data is actually stored column-wise, but we create it row-wise here.
	expected := util.TransposeArray([][]uint8{
		{0, 2, 5, 7, 10, 12, 13, 14},
		{1, 3, 6, 9, 12, 15, 13, 12},
		{2, 5, 8, 11, 15, 18, 14, 10},
		{3, 7, 10, 13, 17, 19, 19, 19},
		{5, 9, 12, 16, 20, 20, 20, 20},
	})

	img, err := Decode([4][]EncodedArea{
		areas,
		areas,
		areas,
		areas,
	})

	util.AssertNil(t, err)
	util.AssertEqual(t, 8, img.Width)
	util.AssertEqual(t, 5, img.Height)
	util.AssertArrayEqual(t, expected, img.R)
	util.AssertArrayEqual(t, expected, img.G)
	util.AssertArrayEqual(t, expected, img.B)
	util.AssertArrayEqual(t, expected, img.A)
}
