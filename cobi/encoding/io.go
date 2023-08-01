package encoding

import (
	"os"
)

func Write(filePath string, areas [4][]EncodedArea) error {
	var data []uint8

	for _, channel := range areas {
		for _, area := range channel {
			data = append(data, serialize(area)...)
		}
	}

	return os.WriteFile(filePath, data, 0644)
}

func serialize(area EncodedArea) []uint8 {
	return []byte{
		area.Values[0],
		area.Values[1],
		area.Values[2],
		area.Values[3],
		area.W,
	}
}
