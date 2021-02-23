package run

import (
	"github.com/faiface/pixel"
	"math"
)

// Adapted code from https://stackoverflow.com/questions/3018313/algorithm-to-convert-rgb-to-hsv-and-hsv-to-rgb-in-range-0-255-for-both

func HsvToRgb(h, s, v float64) pixel.RGBA { // input from 0.0 to 1.0
	if s <= 0.0 {
		return pixel.RGB(v, v, v)
	}
	hh := h * 6
	for hh < 0 {
		hh += 6
	}
	for hh >= 6 {
		hh -= 6.
	}
	i := math.Floor(hh)
	ff := hh - i
	p := v * (1 - s)
	q := v * (1 - s*ff)
	t := v * (1 - s*(1-ff))
	switch i {
	case 0.:
		return pixel.RGB(v, t, p)
	case 1.:
		return pixel.RGB(q, v, p)
	case 2.:
		return pixel.RGB(p, v, t)
	case 3.:
		return pixel.RGB(p, q, v)
	case 4.:
		return pixel.RGB(t, p, v)
	default:
		return pixel.RGB(v, p, q)
	}
}
