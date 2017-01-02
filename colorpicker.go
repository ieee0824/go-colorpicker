package colorpicker

import (
	"image"
	"image/color/palette"

	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mdesenfants/gokmeans"
)

func thinOutColor(img image.Image) image.Image {
	paletted := image.NewPaletted(img.Bounds(), palette.WebSafe)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			paletted.Set(x, y, img.At(x, y))
		}
	}
	return paletted
}

func ExtractTypicalColors(img image.Image, k int) []colorful.Color {
	img = imaging.AdjustContrast(img, 95)

	indexes := make([]int, k)
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	selectColor := make([]colorful.Color, 0, w*h)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := colorful.Color{
				R: float64(r>>8) / 255,
				G: float64(g>>8) / 255,
				B: float64(b>>8) / 255,
			}
			h, s, v := c.Hsv()

			if s > 0.5 && v > 0.7 {
				selectColor = append(selectColor, colorful.Hsv(h, s, 1))
			}
		}
	}

	observations := func() []gokmeans.Node {
		ret := make([]gokmeans.Node, len(selectColor))
		for i := 0; i < len(selectColor); i++ {
			ret[i] = func() []float64 {
				h, s, v := selectColor[i].Hsv()
				return []float64{h, s, v}
			}()
		}
		return ret
	}()
	success, centroids := gokmeans.Train(observations, k, len(selectColor))
	if success {
		for _, observation := range observations {
			index := gokmeans.Nearest(observation, centroids)
			indexes[index]++
		}
	}

	ret := make([]colorful.Color, 0, len(centroids))
	for i, centroid := range centroids {
		if float64(indexes[i])/float64(len(observations)) < 0.01 {
			continue
		}
		hsv := centroid
		ret = append(ret, colorful.Hsv(hsv[0], hsv[1], hsv[2]))
	}

	return ret
}
