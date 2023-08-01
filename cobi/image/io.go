package image

import "image"

type Reader interface {
	Read(filePath string) (*Image, error)
}

type Writer interface {
	Write(filepath string, img image.Image) error
}
