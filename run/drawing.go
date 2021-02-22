package run

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
	"time"
)

//                    Window size
var Draw = drawVars{W: 800, H: 600}

type drawVars struct {
	W, H  float64
	Imd   *imdraw.IMDraw
	Win   *pixelgl.Window
	Timer <-chan time.Time
}

func (draw *drawVars) Init(title string) {
	rand.Seed(time.Now().UnixNano())

	var err error
	draw.Win, err = pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, draw.W, draw.H),
	})
	if err != nil {
		panic(err)
	}

	draw.Timer = time.Tick(time.Second / 30)
	draw.Imd = imdraw.New(nil)
}

func (draw *drawVars) Update(fillColor pixel.RGBA) {
	if fillColor.A > 0 {
		draw.Win.Clear(fillColor)
	}
	draw.Imd.Draw(Draw.Win)
	draw.Imd.Clear()
	draw.Win.Update()
}

func (draw *drawVars) Rect(x, y, w, h float64, color pixel.RGBA) {
	draw.Imd.Color = color
	draw.Imd.Push(pixel.V(x, y), pixel.V(x+w, y+h))
	draw.Imd.Rectangle(0)
}

func (draw *drawVars) Circle(x, y, r float64, color pixel.RGBA) {
	draw.Imd.Color = color
	draw.Imd.Push(pixel.V(x, y))
	draw.Imd.Circle(r, 0)
}

func (draw *drawVars) Segment(x0, y0, x1, y1, thickness float64, color pixel.RGBA) {
	draw.Imd.Color = color
	draw.Imd.Push(pixel.V(x0, y0), pixel.V(x1, y1))
	draw.Imd.Line(thickness)
}
