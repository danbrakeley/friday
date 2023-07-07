package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var shapeSrc []Point = []Point{
	{X: 0, Y: 15},
	{X: 2, Y: 20},
	{X: 9, Y: 20},
	{X: 20, Y: 9},
	{X: 20, Y: -9},
	{X: 9, Y: -20},
	{X: -9, Y: -20},
	{X: -20, Y: -9},
	{X: -20, Y: 9},
	{X: -9, Y: 20},
	{X: -2, Y: 20},
}

type GameScene struct {
	midiMgr *MidiMgr
	center  Point
	rot     int // [0,127]
	shape   []Point
	rotGoal int // [0,127]
}

func NewGameScene(midiMgr *MidiMgr) *GameScene {
	return &GameScene{
		midiMgr: midiMgr,
		center:  Point{X: screenWidth / 2, Y: screenHeight / 2},
		rotGoal: rand.Intn(128),
	}
}

const (
	sliceCount = 64 // number of keypresses to complete one full rotation
	twoPi      = math.Pi * 2
	sliceRad   = twoPi / sliceCount       // radian offset per keypress
	minRot     = twoPi / (sliceCount * 2) // anything under this rounds down to 0 (aka 2π)
	maxRot     = twoPi - minRot           // anything over this rounds up to 2π (aka 0)

	scale = 3
)

func (g *GameScene) Update(mgr *SceneMgr) error {
	g.midiMgr.Update()

	prevRot := g.rot
	g.rot = g.midiMgr.Knob1()

	// if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
	// 	g.rot -= sliceRad
	// 	if g.rot < -minRot {
	// 		g.rot = twoPi - sliceRad
	// 	}
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
	// 	g.rot += sliceRad
	// 	if g.rot > maxRot {
	// 		g.rot = 0
	// 	}
	// }

	if g.rot != prevRot || len(g.shape) == 0 {
		// regenerate vertices from shape
		g.shape = make([]Point, len(shapeSrc))
		rads := float64(g.rot) * twoPi / 127.0
		cosRot := float32(math.Cos(rads))
		sinRot := float32(math.Sin(rads))
		for i, v := range shapeSrc {
			g.shape[i] = Point{
				X: scale*(v.X*cosRot-v.Y*sinRot) + g.center.X,
				Y: scale*(v.X*sinRot+v.Y*cosRot) + g.center.Y,
			}
		}
	}

	return nil
}

var (
	colorBG = color.RGBA{0x56, 0x55, 0x54, 0xff}
	colorFG = color.RGBA{0xf6, 0xf1, 0x93, 0xee}
)

func (g *GameScene) Draw(mgr *SceneMgr, screen *ebiten.Image) {
	screen.Fill(colorBG)
	drawShape(screen, g.shape, 1, colorFG)

	// msg := fmt.Sprintf("TPS: %0.2f\nRot: %.3f", ebiten.ActualTPS(), g.rot)
	msg := "Spin the dial with left and right arrows"
	if g.rot == g.rotGoal {
		msg += "\n\nCLICK!"
	}
	ebitenutil.DebugPrint(screen, msg)
}
