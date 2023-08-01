package interpolate

import (
	"cobi/util"
	"testing"
)

func TestInterpolation(t *testing.T) {
	// 0 . 8
	// . . .
	// . . .
	// 5 . 10
	// TransposeArray needed because the image data is actually stored column-wise, but we create it row-wise here.
	expected := util.TransposeArray([][]byte{
		{0, 4, 8},
		{1, 5, 8},
		{3, 6, 9},
		{5, 7, 10},
	})

	actual := Interpolate(3, 4, [4]byte{0, 8, 5, 10})

	util.AssertArrayEqual(t, expected, actual)
}

func TestInterpolation_noInterpolationNeeded(t *testing.T) {
	// 0 8
	// 5 10
	// TransposeArray needed because the image data is actually stored column-wise, but we create it row-wise here.
	expected := util.TransposeArray([][]byte{
		{0, 8},
		{5, 10},
	})

	actual := Interpolate(2, 2, [4]byte{0, 8, 5, 10})

	util.AssertArrayEqual(t, expected, actual)
}

func TestInterpolation_equalValues(t *testing.T) {
	// 0 . 8
	// . . .
	// . . .
	// 5 . 8
	// TransposeArray needed because the image data is actually stored column-wise, but we create it row-wise here.
	expected := util.TransposeArray([][]byte{
		{0, 4, 8},
		{1, 4, 8},
		{3, 5, 8},
		{5, 6, 8},
	})

	actual := Interpolate(3, 4, [4]byte{0, 8, 5, 8})

	util.AssertArrayEqual(t, expected, actual)
}

func TestInterpolation_singleRow(t *testing.T) {
	// 0 . . 6
	// TransposeArray needed because the image data is actually stored column-wise, but we create it row-wise here.
	expected := util.TransposeArray([][]byte{
		{0, 2, 4, 6},
	})

	actual := Interpolate(4, 1, [4]byte{0, 6, 0, 6})

	util.AssertArrayEqual(t, expected, actual)
}

func TestInterpolation_singleColumn(t *testing.T) {
	// 0
	// .
	// .
	// 6
	// TransposeArray needed because the image data is actually stored column-wise, but we create it row-wise here.
	expected := util.TransposeArray([][]byte{
		{0},
		{2},
		{4},
		{6},
	})

	actual := Interpolate(1, 4, [4]byte{0, 0, 6, 6})

	util.AssertArrayEqual(t, expected, actual)
}

func TestInterpolation_singlePoint(t *testing.T) {
	// 42
	expected := [][]byte{
		{42},
	}

	actual := Interpolate(1, 1, [4]byte{42, 42, 42, 42})

	util.AssertArrayEqual(t, expected, actual)
}