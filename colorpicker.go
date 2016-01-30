package colorpicker

import (
	"image"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mdesenfants/gokmeans"
)

func ExtractTypicalColors(img image.Image, k int) []colorful.Color {
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
			_, s, v := c.Hsv()

			if s > 0.5 && v > 0.5 {
				selectColor = append(selectColor, c)
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
