package png

import (
	"cobi/image"
	"fmt"
	"github.com/pkg/errors"
	"image/png"
	"os"
)

type Writer struct{}

func (w *Writer) Write(filePath string, img image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not open output image %s", filePath))
	}

	err = png.Encode(file, &img)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could encode or write PNG file %s", filePath))
	}

	return nil
}
