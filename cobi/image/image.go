package image

import (
	"fmt"
	"github.com/pkg/errors"
	"image"
	"image/color"
)

type Image struct {
	R      [][]uint8
	G      [][]uint8
	B      [][]uint8
	A      [][]uint8
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

func (img *Image) ColorModel() color.Model {
	return color.AlphaModel
}

func (img *Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img *Image) At(x, y int) color.Color {
	return color.RGBA{
		R: img.R[x][y],
		G: img.G[x][y],
		B: img.B[x][y],
		A: img.A[x][y],
	}
}

// New create an Image instance in which the RGBA arrays are initialized but empty. The first dimension of the array
// represents the horizontal x-direction in the image. The second array dimension the vertical y-direction.
func New(width, height int) *Image {
	return &Image{
		R:      make([][]uint8, width),
		G:      make([][]uint8, width),
		B:      make([][]uint8, width),
		A:      make([][]uint8, width),
		Width:  width,
		Height: height,
	}
}

func FromGoImage(goImg image.Image) (*Image, error) {
	w, h := goImg.Bounds().Max.X, goImg.Bounds().Max.Y

	img := New(w, h)

	// convert data from uint32 into uint8s for smaller images
	for x := 0; x < w; x++ {
		img.R[x] = make([]uint8, h)
		img.G[x] = make([]uint8, h)
		img.B[x] = make([]uint8, h)
		img.A[x] = make([]uint8, h)

		for y := 0; y < h; y++ {
			r, g, b, a := goImg.At(x, y).RGBA()
			r /= 256
			g /= 256
			b /= 256
			a /= 256

			if r > 255 || g > 255 || b > 255 || a > 255 {
				return nil, errors.New(fmt.Sprintf("Only 8-bit per channel are allowed but found (R,G,B,A): %d,%d,%d,%d", r, g, b, a))
			}

			img.R[x][y] = uint8(r)
			img.G[x][y] = uint8(g)
			img.B[x][y] = uint8(b)
			img.A[x][y] = uint8(a)
		}
	}

	return img, nil
}
