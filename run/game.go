package run

import (
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

const (
	RayCount = InpCount
)

var (
	MaxSpeed       = 450.
	CarSpawnChance = 0.05
)

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
	C    pixel.RGBA
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
	Draw.Segment(ray.X, ray.Y, ray.X+x, ray.Y+ray.K*x, 1, ray.C)
}

type Car struct {
	Rect
	speed float64
}

func (c *Car) Update() {
	c.Y -= c.speed / 30
	if c.Y+c.H < 0 {
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

func (m *carManager) Update() {
	if rand.Float64() < CarSpawnChance {
		w := 30 + rand.Float64()*30
		c := &Car{
			Rect: Rect{
				X: rand.Float64() * (Draw.W - w),
				Y: Draw.H,
				W: w,
				H: 60 + rand.Float64()*60,
				C: HsvToRgb(-0.1+rand.Float64()*0.2, 1, 0.5),
			},
			speed: 150 + rand.Float64()*150,
		}
		AddRect(&c.Rect)
		AddUpDraw(c, true, true)
	}
}

type Player struct {
	Rect
	Rays      [RayCount]*DistRay
	DeathTime float64
	//Net *Network
	NetIndex int
}

func (p *Player) Init() {
	r := 0.
	for i := 0; i < RayCount; i++ {
		a := math.Pi * 2 * (0.5 + (r+math.Sin(r*2*math.Pi)/4/math.Pi)*0.5)
		ray := DistRay{Ray: GetRayA(pixel.V(p.X+p.W/2, p.Y), a), C: p.C.Scaled(0.1)}
		p.Rays[i] = &ray
		r += 1. / (RayCount - 1)
		AddUpDraw(&ray, false, true)
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
	distances := [RayCount]float64{}
	for i, ray := range p.Rays {
		ray.X = p.X + p.W/2
		ray.Update()
		distances[i] = ray.dist / 100.
	}
	g := AllNetworks[p.NetIndex].Get(distances)[0]
	//if p.NetIndex == 0 {
	//	fmt.Println(g, distances)
	//}
	p.X += MaxSpeed * (g*2 - 1) / 30.
	if p.DeathTime < 0 {
		AllNetworks[p.NetIndex].Fitness += 1. / 30
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
	//outOfBounds := p.X < 0 || p.X + p.W > Draw.W
	p.X = -p.W/2 + Draw.W/2 + (rand.Float64()*2-1)*Draw.W/3
	p.DeathTime = 2
	//for i := 0; i < len(Rects); i++ {
	//	if p.Rect.IsCollideRect(*Rects[i]) {
	//		Rects[i].X = -100 - Rects[i].W
	//	}
	//}
	AllNetworks[p.NetIndex] = NewNetwork()
	//if outOfBounds {
	//	AllNetworks[p.NetIndex].Fitness = -1000
	//} else {
	AllNetworks[p.NetIndex].Fitness = 0
	//}
}
