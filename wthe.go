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

type HSV struct { // {0..1}
	H, S, V float64
}

type RGB struct { // {0..255}
	R, G, B float64
}

func Wthe(imagename string) *image.RGBA {
	m, err := getImageFromFilePath(imagename)
	if err != nil {
		fmt.Println(err) // debugging
	}

	img_bounds := m.Bounds()
	m_hsv := image.NewRGBA(img_bounds)
	m_v_of_hsv := image.NewGray(img_bounds)
	m_out := image.NewRGBA(img_bounds)

	// convert rgb image to hsv
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			r, g, b, _ := m.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			rgb := &RGB{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
			// convert rgb to hsv
			hsv := rgb.HSV()
			h := float2uint8(hsv.H)
			s := float2uint8(hsv.S)
			v := float2uint8(hsv.V)
			m_hsv.Set(x, y, color.RGBA{h, s, v, 255})
			m_v_of_hsv.Set(x, y, color.RGBA{v, v, v, 255})

			// // check if rgb2hsv and then hsv2rgb of the image is equal to the image
			// rc := float64(r) / 65535
			// gc := float64(g) / 65535
			// bc := float64(b) / 65535
			// rgb_out := hsv.RGB()
			// rn := rgb_out.R
			// gn := rgb_out.G
			// bn := rgb_out.B

			// // fmt.Println(rc - rn)
			// // fmt.Println(gc - gn)
			// // fmt.Println(bc - bn)

			// if (rc - rn) > 1e-15 {
			// 	fmt.Println(rc-rn, r == uint32(math.Round(rn*65535)))
			// } else if (gc - gn) > 1e-15 {
			// 	fmt.Println(gc-gn, g == uint32(math.Round(gn*65535)))
			// } else if (bc - bn) > 1e-15 {
			// 	fmt.Println(bc-bn, b == uint32(math.Round(bn*65535)))
			// }
		}
	}

	v_hist := imhist(m_v_of_hsv)
	h_v := img_bounds.Max.Y
	w_v := img_bounds.Max.X

	v_pmf := [256]float64{}
	for i, _ := range v_hist {
		v_pmf[i] = float64(v_hist[i]) / float64(h_v*w_v)
	}

	r := 0.5
	v := 0.5
	Pl := 1e-5
	Pu := v * max(v_pmf)

	v_pmf_new := v_pmf
	for i, v := range v_pmf {
		if v < Pl {
			v_pmf_new[i] = 0
		} else if v > Pu {
			v_pmf_new[i] = Pu
		} else {
			v_pmf_new[i] = (math.Pow((v-Pl)/(Pu-Pl), r)) * Pu
		}
	}

	Win := float64(len(whereLessThan(v_pmf, 0)))
	Gmax := 1.5 // 1.5 .. 2
	Wout := math.Min(255.0, Gmax*Win)
	// fmt.Println("Wout:", Wout) // debugging

	v_cdf := cumsum(v_pmf)

	// m_v_of_hsv_mean
	m_v_of_hsv_mean := meanOfGray(m_v_of_hsv)

	// make changes to the value layer of the hsv image
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			h, s, v, _ := m_hsv.At(x, y).RGBA()
			v_new := uint8(Wout * v_cdf[v>>8])

			m_hsv.Set(x, y, color.RGBA{uint8(h >> 8), uint8(s >> 8), v_new, 255})
			m_v_of_hsv.Set(x, y, color.RGBA{v_new, v_new, v_new, 255})

		}
	}

	// adjust the average brightness of the new hsv image to match the original image
	m_v_of_new_hsv_mean := meanOfGray(m_v_of_hsv)
	mean_diff := m_v_of_hsv_mean - m_v_of_new_hsv_mean
	// fmt.Println(mean_diff) // debugging
	for x := img_bounds.Min.X; x < img_bounds.Max.X; x++ {
		for y := img_bounds.Min.Y; y < img_bounds.Max.Y; y++ {
			h, s, v, _ := m_hsv.At(x, y).RGBA()
			v_new := int32(v>>8) + int32(mean_diff*255)
			if v_new > 255 {
				v_new = 255
			} else if v_new < 0 {
				v_new = 0
			}
			m_hsv.Set(x, y, color.RGBA{uint8(h >> 8), uint8(s >> 8), uint8(v_new), 255})
		}
	}

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
	return sum / float64(bounds.Max.Y*bounds.Max.X)
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

func whereLessThan(arr [256]float64, threshold float64) []int {
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

func (c *RGB) HSV() *HSV {
	var h, s, v float64

	r := c.R
	g := c.G
	b := c.B

	high := math.Max(r, math.Max(g, b))
	low := math.Min(r, math.Min(g, b))

	v = high
	d := high - low

	if d == 0.0 {
		s = 0.0
		h = 0.0 // undefined, maybe nan?
		hsv := &HSV{h, s, v}
		return hsv
	}

	if high > 0.0 { // NOTE: if Max is == 0, this divide would cause a crash
		s = (d / high) // s
	} else {
		// if max is 0, then r = g = b = 0
		// s = 0, h is undefined
		s = 0.0
		h = 0.5 // its now undefined
		hsv := &HSV{h, s, v}
		return hsv
	}

	if r == high {
		offset := 0.0
		if g < b {
			offset = 6.0
		}
		h = (g-b)/d + offset

	} else if g == high {
		h = (b-r)/d + 2.0

	} else if b == high {
		h = (r-g)/d + 4.0
	}
	h = h * 60.0 / 360.0 //  normalize

	hsv := &HSV{h, s, v}
	return hsv
}

func (c *HSV) RGB() *RGB {
	var r, g, b float64

	h := c.H
	s := c.S
	v := c.V
	if h >= 1.0 {
		h = math.Mod(h, 1.0)
	}

	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)

	switch i {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	case 5:
		r = v
		g = p
		b = q
	}
	rgb := &RGB{r, g, b}
	return rgb
}
