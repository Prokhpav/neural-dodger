package run

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
)

// Universal pixelgl update/draw system

var (
	Updaters []interface{ Update() }
	Drawers  []interface{ Draw() }
	Rects    []*Rect

	TurboMode = false
)

const PlayerCount = 20

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

func UpdateWindow() {
	if Draw.Win.JustPressed(pixelgl.KeySpace) {
		TurboMode = !TurboMode
	}
	if Draw.Win.Pressed(pixelgl.KeyMinus) {
		MuteChance -= 0.02
		fmt.Println("MuteChance:", MuteChance)
	}
	if Draw.Win.Pressed(pixelgl.KeyEqual) {
		MuteChance += 0.02
		fmt.Println("MuteChance:", MuteChance)
	}
	if Draw.Win.Pressed(pixelgl.KeyLeftBracket) {
		MuteAmplitude -= 0.02
		fmt.Println("MuteAmplitude:", MuteAmplitude)
	}
	if Draw.Win.Pressed(pixelgl.KeyRightBracket) {
		MuteAmplitude += 0.02
		fmt.Println("MuteAmplitude:", MuteAmplitude)
	}
	if Draw.Win.Pressed(pixelgl.KeyO) {
		CarSpawnChance -= 0.0005
		fmt.Println("CarSpawnChance:", CarSpawnChance)
	}
	if Draw.Win.Pressed(pixelgl.KeyP) {
		CarSpawnChance += 0.0005
		fmt.Println("CarSpawnChance:", CarSpawnChance)
	}
	if Draw.Win.Pressed(pixelgl.KeyI) {
		sumFitness := 0.
		maxFitness := -1000.
		for _, n := range AllNetworks {
			if n.Fitness > maxFitness {
				maxFitness = n.Fitness
			}
			sumFitness += n.Fitness
		}
		fmt.Println("Average fitness:", sumFitness/PlayerCount, "     Max:", maxFitness)
	}

	for _, u := range Updaters {
		u.Update()
	}
	for _, d := range Drawers {
		d.Draw()
	}
	Draw.Update(pixel.RGBA{A: 1})
}

func Run() {
	Draw.Init("Traffic racer")

	CarManager := carManager{}
	Updaters = append(Updaters, &CarManager)

	AddRect(&Rect{X: -5000., Y: -1000, W: 5000, H: 2000 + Draw.H}) // Bounds
	AddRect(&Rect{X: Draw.W, Y: -1000, W: 5000, H: 2000 + Draw.H})
	AddRect(&Rect{X: 0, Y: Draw.H, W: Draw.W, H: 1000})

	for i := 0; i < PlayerCount; i++ {
		AllNetworks = append(AllNetworks, GetRandomNetwork())
		player := Player{
			Rect:     Rect{400, 5 + rand.Float64()*50, 50, 75, HsvToRgb(rand.Float64(), 1, 1)},
			NetIndex: i,
		}
		player.Init()
		AddUpDraw(&player, true, true)
	}

	for !Draw.Win.Closed() {
		if TurboMode {
			select {
			case <-Draw.Timer:
				UpdateWindow()
			default:
				for _, u := range Updaters {
					u.Update()
				}
			}
		} else {
			select {
			case <-Draw.Timer:
				UpdateWindow()
			}
		}

	}
}
