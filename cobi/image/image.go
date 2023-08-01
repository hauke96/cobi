package image

import (
	"fmt"
	"github.com/pkg/errors"
	"image"
)

type Image struct {
	R [][]byte
	G [][]byte
	B [][]byte
	A [][]byte
}

func (img *Image) Print() {
	for y := 0; y < len(img.R[0]); y++ {
		for x := 0; x < len(img.R); x++ {
			fmt.Printf("[%3d %3d %3d], ", img.R[x][y], img.G[x][y], img.B[x][y])
		}
		fmt.Print("\n")
	}
}

func FromGoImage(goImg image.Image) (*Image, error) {
	w, h := goImg.Bounds().Max.X, goImg.Bounds().Max.Y

	img := Image{
		R: make([][]byte, w),
		G: make([][]byte, w),
		B: make([][]byte, w),
		A: make([][]byte, w),
	}

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

	return &img, nil
}
