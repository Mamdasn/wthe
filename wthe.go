package wthe

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"math"
	"os"
)

func Wthe(imagename string) *image.RGBA {
	input_image, err := getImageFromFilePath(imagename)
	if err != nil {
		fmt.Println(err) // debugging
	}

	img_bounds := input_image.Bounds()
	m_out := image.NewRGBA(img_bounds)

	m_hsv := image.NewRGBA(img_bounds)
	v_of_m_hsv := image.NewGray(img_bounds)

	// convert rgb image to hsv
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			r, g, b, _ := input_image.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			rgb := &RGB{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
			// convert rgb to hsv
			hsv := rgb.HSV()
			h := float2uint8(hsv.H)
			s := float2uint8(hsv.S)
			v := float2uint8(hsv.V)

			m_hsv.Set(x, y, color.RGBA{h, s, v, 255})
			v_of_m_hsv.SetGray(x, y, color.Gray{Y: v})
		}
	}

	v_hist := imhist(v_of_m_hsv)
	h_v := img_bounds.Max.Y - img_bounds.Min.Y + 1
	w_v := img_bounds.Max.X - img_bounds.Min.X + 1

	v_pmf := [256]float64{}
	for i, _ := range v_hist {
		v_pmf[i] = float64(v_hist[i]) / float64(h_v*w_v)
	}

	r := 0.5
	v := 0.5
	Pl := 1e-5
	Pu := v * max(v_pmf)

	v_pmf_modified := v_pmf
	for i, v := range v_pmf {
		if v <= Pl {
			v_pmf_modified[i] = 0
		} else if v > Pu {
			v_pmf_modified[i] = Pu
		} else {
			v_pmf_modified[i] = (math.Pow((v-Pl)/(Pu-Pl), r)) * Pu
		}
	}

	Win := float64(len(whereBiggerThan(v_pmf, 0)))
	Gmax := 1.5 // 1.5 .. 2
	Wout := math.Min(255.0, Gmax*Win)

	v_cdf := cumsum(v_pmf_modified)
	for i, _ := range v_cdf {
		v_cdf[i] /= v_cdf[len(v_cdf)-1]
	}
	// make changes to the value layer of the hsv image
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			h, s, v, _ := m_hsv.At(x, y).RGBA()
			v_new := uint8(Wout * v_cdf[v>>8])

			m_hsv.Set(x, y, color.RGBA{uint8(h >> 8), uint8(s >> 8), v_new, 255})
			// // adjust the average brightness of the new hsv image to match the original image
			// v_of_m_hsv.Set(x, y, color.RGBA{v_new, v_new, v_new, 255})

		}
	}

	// // adjust the average brightness of the new hsv image to match the original image
	// m_v_of_new_hsv_mean := meanOfGray(v_of_m_hsv)
	// mean_diff := m_v_of_hsv_mean - m_v_of_new_hsv_mean
	// // fmt.Println(mean_diff) // debugging
	// for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
	// 	for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
	// 		h, s, v, _ := m_hsv.At(x, y).RGBA()
	// 		v_new := int32(v>>8) + int32(mean_diff*255)
	// 		if v_new > 255 {
	// 			v_new = 255
	// 		} else if v_new < 0 {
	// 			v_new = 0
	// 		}
	// 		m_hsv.Set(x, y, color.RGBA{uint8(h >> 8), uint8(s >> 8), uint8(v_new), 255})
	// 	}
	// }

	// convert hsv to rgb, img_out
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			h, s, v, _ := m_hsv.At(x, y).RGBA()

			hsv_out := &HSV{float64(h) / 65535, float64(s) / 65535, float64(v) / 65535}
			hsv2rgb_out := hsv_out.RGB()
			r := float2uint8(hsv2rgb_out.R)
			g := float2uint8(hsv2rgb_out.G)
			b := float2uint8(hsv2rgb_out.B)
			m_out.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return m_out
}
func float2uint8(c float64) uint8 {
	return uint8(math.Round(c * 255))
}

func meanOfGray(img *image.Gray) float64 {
	var sum float64 = 0
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gg, _, _, _ := img.At(x, y).RGBA()
			g := float64(gg) / 65535
			sum += g
		}
	}
	return sum / float64((bounds.Max.Y-bounds.Min.Y+1)*(bounds.Max.X-bounds.Min.X+1))
}

func imhist(img *image.Gray) [256]int {
	var histogram [256]int
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			g, _, _, _ := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			histogram[g>>8]++
		}
	}
	return histogram
}
func cumsum(pmf [256]float64) [256]float64 {
	var cum [256]float64
	cum[0] = pmf[0]
	for i := 1; i < 256; i++ {
		cum[i] = cum[i-1] + pmf[i]
	}
	return cum
}
func max(arr [256]float64) float64 {
	maxm := arr[0]
	for i, _ := range arr {
		if maxm < arr[i] {
			maxm = arr[i]
		}
	}
	return maxm
}

func whereBiggerThan(arr [256]float64, threshold float64) []int {
	where := []int{}
	for i := 0; i < 256; i++ {
		if arr[i] > threshold {
			where = append(where, i)
		}
	}
	return where
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func SaveImageToFilePath(filePath string, img *image.RGBA) error {
	// Somewhere in the same package
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Specify the quality, between 0-100
	// Higher is better
	opt := jpeg.Options{
		Quality: 100,
	}
	err = jpeg.Encode(f, img, &opt)
	if err != nil {
		return err
	}
	return nil
}

type HSV struct { // {0..1}
	H, S, V float64
}

func (c *RGB) HSV() *HSV {
	r, g, b := c.R, c.G, c.B
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	v := max
	delta := max - min

	var h, s float64

	if max != 0 {
		s = delta / max
	} else {
		// r = g = b = 0
		s = 0
		h = -1 // Undefined
		return &HSV{H: h, S: s, V: v}
	}

	if r == max {
		h = (g - b) / delta // Between yellow & magenta
	} else if g == max {
		h = 2 + (b-r)/delta // Between cyan & yellow
	} else {
		h = 4 + (r-g)/delta // Between magenta & cyan
	}

	h *= 60 // degrees
	if h < 0 {
		h += 360
	}

	// Normalize H to [0, 1)
	h = h / 360

	return &HSV{H: h, S: s, V: v}
}

type RGB struct { // {0..255}
	R, G, B float64
}

func (c *HSV) RGB() *RGB {
	h, s, v := c.H, c.S, c.V

	h = math.Mod(h, 1.0) // Ensures h is within [0, 1)

	region := int(math.Floor(h * 6))
	fraction := h*6 - float64(region)
	p := v * (1.0 - s)
	q := v * (1 - fraction*s)
	t := v * (1 - (1-fraction)*s)

	var r, g, b float64

	switch region {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return &RGB{R: r, G: g, B: b}
}
