package png

import (
	"cobi/image"
	"fmt"
	"github.com/pkg/errors"
	"image/png"
	"os"
)

type Reader struct{}

func (r *Reader) Read(filePath string) (*image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not open input image %s", filePath))
	}
	defer file.Close()

	pngImage, err := png.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not decode input image %s", filePath))
	}

	return image.FromGoImage(pngImage)
}
