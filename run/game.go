package run

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
)

// game classes

func AddUpDraw(obj interface {
	Update()
	Draw()
}, update, draw bool) {
	if update {
		Updaters = append(Updaters, obj)
	}
	if draw {
		Drawers = append(Drawers, obj)
	}
}

func AddRect(rect *Rect) {
	Rects = append(Rects, rect)
}

type DistRay struct {
	Ray
	dist float64
}

func (ray *DistRay) Update() {
	ray.dist = ray.RectsMinDist(Rects)
}

func (ray DistRay) Draw() {
	dist := ray.dist
	if math.IsInf(dist, 1) {
		dist = 1000
	}
	if !ray.P {
		dist = -dist
	}
	x := dist / math.Sqrt(ray.K*ray.K+1)
	Draw.Segment(ray.X, ray.Y, ray.X+x, ray.Y+ray.K*x, 4, pixel.RGB(0.5, 0.5, 0))
}

type Car struct {
	Rect
	speed float64
}

func (c *Car) Update() {
	c.Y -= c.speed / 30
	if c.Y+c.H < 10 {
		for i, r := range Rects {
			if *r == c.Rect {
				Rects = append(Rects[:i], Rects[i+1:]...)
				break
			}
		}
		Destroy(c)
	}
}

type carManager struct{}

var CarManager = carManager{}

func (m *carManager) Update() {
	if rand.Float64() < 0.05 {
		w := 50.
		c := &Car{
			Rect: Rect{
				X: rand.Float64() * (Draw.W - w),
				Y: Draw.H,
				W: w,
				H: 80,
				C: pixel.RGB(0.5, 0.7, 0),
			},
			speed: 100 + rand.Float64()*100,
		}
		AddRect(&c.Rect)
		AddUpDraw(c, true, false)
	}
}

type Player struct {
	Rect
	Rays      []*DistRay
	DeathTime float64
}

func (p *Player) Init() {
	for r := 0.; r <= 1.001; r += 0.01 {
		a := math.Pi * 2 * (0.49 + (r+math.Sin(r*2*math.Pi)/4/math.Pi)*0.52)
		ray := DistRay{GetRayA(pixel.V(p.X+p.W/2, p.Y+p.H/2), a), 10}
		p.Rays = append(p.Rays, &ray)
		Drawers = append(Drawers, &ray)
	}
}

func (p Player) Draw() {
	c := p.C
	if p.DeathTime >= 0 {
		c = c.Scaled(0.5)
	}
	Draw.Rect(p.X, p.Y, p.W, p.H, c)
}

func (p *Player) Update() {
	if Draw.Win.Pressed(pixelgl.KeyA) {
		p.X -= 200. / 30
	}
	if Draw.Win.Pressed(pixelgl.KeyD) {
		p.X += 200. / 30
	}
	for _, ray := range p.Rays {
		ray.X = p.X + p.W/2
		ray.dist = ray.RectsMinDist(Rects)
	}
	if p.DeathTime < 0 {
		for _, r := range Rects {
			if p.Rect.IsCollideRect(*r) {
				p.OnDeath()
				break
			}
		}
	} else {
		p.DeathTime -= 1. / 30
	}

}

func (p *Player) OnDeath() {
	p.X = (Draw.W - p.W) / 2
	p.DeathTime = 2
	for i := 0; i < len(Rects); i++ {
		if p.Rect.IsCollideRect(*Rects[i]) {
			Rects[i].X = -100 - Rects[i].W
		}
	}
}
