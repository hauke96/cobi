package image

import (
	"fmt"
	"github.com/pkg/errors"
	"image"
)

type Image struct {
	R      [][]byte
	G      [][]byte
	B      [][]byte
	A      [][]byte
	Width  int
	Height int
}

func (img *Image) Print() {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			fmt.Printf("[%3d %3d %3d], ", img.R[x][y], img.G[x][y], img.B[x][y])
		}
		fmt.Print("\n")
	}
}

// New create an Image instance in which the RGBA arrays are initialized but empty. The first dimension of the array
// represents the horizontal x-direction in the image. The second array dimension the vertical y-direction.
func New(width, height int) *Image {
	return &Image{
		R:      make([][]byte, width),
		G:      make([][]byte, width),
		B:      make([][]byte, width),
		A:      make([][]byte, width),
		Width:  width,
		Height: height,
	}
}

func FromGoImage(goImg image.Image) (*Image, error) {
	w, h := goImg.Bounds().Max.X, goImg.Bounds().Max.Y

	img := New(w, h)

	// convert data from uint32 into bytes for smaller images
	for x := 0; x < w; x++ {
		img.R[x] = make([]byte, h)
		img.G[x] = make([]byte, h)
		img.B[x] = make([]byte, h)
		img.A[x] = make([]byte, h)

		for y := 0; y < h; y++ {
			r, g, b, a := goImg.At(x, y).RGBA()
			r /= 256
			g /= 256
			b /= 256
			a /= 256

			if r > 255 || g > 255 || b > 255 || a > 255 {
				return nil, errors.New(fmt.Sprintf("Only 8-bit per channel are allowed but found (R,G,B,A): %d,%d,%d,%d", r, g, b, a))
			}

			img.R[x][y] = byte(r)
			img.G[x][y] = byte(g)
			img.B[x][y] = byte(b)
			img.A[x][y] = byte(a)
		}
	}

	return img, nil
}
