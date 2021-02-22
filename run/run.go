package run

import (
	"github.com/faiface/pixel"
)

// Universal pixelgl update/draw system

var (
	Updaters []interface{ Update() }
	Drawers  []interface{ Draw() }
	Rects    []*Rect
)

func Destroy(obj interface{}) {
	for i, u := range Updaters {
		if obj == u {
			Updaters = append(Updaters[:i], Updaters[i+1:]...)
			break
		}
	}
	for i, d := range Drawers {
		if obj == d {
			Drawers = append(Drawers[:i], Drawers[i+1:]...)
			break
		}
	}
}

func Run() {
	Draw.Init("Traffic racer")

	Updaters = append(Updaters, &CarManager)

	player := Player{Rect: Rect{400, 50, 50, 50, pixel.RGB(0.5, 1, 1)}}
	player.Init()
	AddUpDraw(&player, true, true)
	AddRect(&Rect{X: -1000, Y: -1000, W: 1000, H: 2000 + Draw.H})
	AddRect(&Rect{X: Draw.W, Y: -1000, W: 1000, H: 2000 + Draw.H})
	AddRect(&Rect{X: 0, Y: Draw.H, W: Draw.W, H: 1000})

	for !Draw.Win.Closed() {
		select {
		case <-Draw.Timer:
			for _, u := range Updaters {
				u.Update()
			}
			for _, d := range Drawers {
				d.Draw()
			}
			Draw.Update(pixel.RGBA{A: 1})
		}
	}
}
