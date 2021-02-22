package run

import (
	"github.com/faiface/pixel"
	"math"
)

// Ray and Rect structures; their collisions

type Rect struct {
	X, Y, W, H float64
	C          pixel.RGBA
}

func (r Rect) IsCollideRect(r2 Rect) bool {
	if r.X > r2.X+r2.W || r.X+r.W < r2.X || r.Y > r2.Y+r2.H || r.Y+r.H < r2.Y {
		return false
	}
	return true
}

func (r Rect) Draw() {
	Draw.Rect(r.X, r.Y, r.W, r.H, r.C)
}

type Ray struct {
	X, Y, K float64
	P       bool // направлен ли в +X
}

//func GetRayP(b, e pixel.Vec) Ray {  // not used
//	var k float64
//	if b.X != e.X {
//		k = (e.Y - b.Y) / (e.X - b.X)
//	} else {
//		k = math.Inf(1)
//	}
//	return Ray{b.X, b.Y, k, b.X <= e.X}
//}

func GetRayA(b pixel.Vec, a float64) Ray {
	for a < 0 {
		a += math.Pi * 2
	}
	for a > math.Pi*2 {
		a -= math.Pi * 2
	}
	return Ray{b.X, b.Y, math.Tan(a), math.Pi*0.5 <= a && a <= math.Pi*1.5}
}

func (ray *Ray) SetAngle(a float64) {
	for a < 0 {
		a += math.Pi * 2
	}
	for a > math.Pi*2 {
		a -= math.Pi * 2
	}
	ray.K = math.Tan(a)
	ray.P = math.Pi*0.5 < a && a < math.Pi*1.5
}

func (ray Ray) Pos() pixel.Vec {
	return pixel.V(ray.X, ray.Y)
}

func (ray Ray) RectIntersect(rect Rect) (isCollide bool, collidePoint pixel.Vec) {
	rect = Rect{X: rect.X - ray.X, Y: rect.Y - ray.Y, W: rect.W, H: rect.H}
	k := ray.K

	mw := !ray.P
	if mw {
		rect.X = -rect.X - rect.W
		if !math.IsInf(k, 1) {
			k = -k
		}
	}
	mh := k < 0
	if mh {
		rect.Y = -rect.Y - rect.H
		if !math.IsInf(k, 1) {
			k = math.Abs(k)
		}
	}

	if rect.X+rect.W < 0 || rect.Y+rect.H < 0 || k < 0 {
		goto nonCollide
	}
	if rect.X < 0 {
		if rect.Y < 0 {
			collidePoint = pixel.ZV
		} else {
			x := rect.Y / k
			if x > rect.X+rect.W {
				goto nonCollide
			}
			collidePoint = pixel.V(x, rect.Y)
		}

	} else if rect.Y < 0 {
		y := rect.X * k
		if y > rect.Y+rect.H {
			goto nonCollide
		}
		collidePoint = pixel.V(rect.X, y)

	} else {
		y := k * rect.X
		if y > rect.Y+rect.H {
			goto nonCollide
		} else if y >= rect.Y {
			collidePoint = pixel.V(rect.X, y)
		} else {
			x := rect.Y / k
			if x > rect.X+rect.W {
				goto nonCollide
			}
			collidePoint = pixel.V(x, rect.Y)
		}
	}

	if mw {
		collidePoint.X = -collidePoint.X
	}
	if mh {
		collidePoint.Y = -collidePoint.Y
	}
	return true, collidePoint.Add(pixel.V(ray.X, ray.Y))
nonCollide:
	return false, pixel.V(math.Inf(1), math.Inf(1))
	//if ray.P {
	//	if ray.K > 0 { // ◳
	//
	//	} else { //     ◲
	//
	//	}
	//} else {
	//	if ray.K > 0 { // ◱
	//
	//	} else { //     ◰
	//
	//	}
	//}
}

func (ray Ray) RectsMinDist(rects []*Rect) float64 {
	m := math.Inf(1)
	for _, rect := range rects {
		c, p := ray.RectIntersect(*rect)
		if c {
			m = math.Min(m, p.Sub(ray.Pos()).Len())
		}
	}
	return m
}
