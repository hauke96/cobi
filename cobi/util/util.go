package util

import "fmt"

func PrintArray[T any](v [][]T) {
	for y := 0; y < len(v[0]); y++ {
		fmt.Print("[")
		for x := 0; x < len(v); x++ {
			fmt.Printf("%3d ", v[x][y])
		}
		fmt.Print("]\n")
	}
}

func TransposeArray[T any](slice [][]T) [][]T {
	width := len(slice)
	height := len(slice[0])
	result := make([][]T, height)
	for i := range result {
		result[i] = make([]T, width)
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
