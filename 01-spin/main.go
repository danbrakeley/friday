package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}

func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.AntiAlias = true
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

func drawShape(dst *ebiten.Image, shape []Point, strokeWidth float32, clr color.Color) {
	if len(shape) == 0 {
		return
	}

	var path vector.Path

	end := len(shape) - 1
	path.MoveTo(shape[end].X, shape[end].Y)
	for _, v := range shape {
		path.LineTo(v.X, v.Y)
	}

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = strokeWidth
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)

	drawVerticesForUtil(dst, vs, is, clr)
}

type Point struct {
	X, Y float32
}

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

type Game struct {
	center  Point
	rot     float64
	shape   []Point
	rotGoal float64
}

func NewGame() *Game {
	return &Game{
		center:  Point{X: screenWidth / 2, Y: screenHeight / 2},
		rotGoal: float64(rand.Intn(sliceCount)) * twoPi / float64(sliceCount),
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

func (g *Game) Update() error {
	prevRot := g.rot

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.rot -= sliceRad
		if g.rot < -minRot {
			g.rot = twoPi - sliceRad
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.rot += sliceRad
		if g.rot > maxRot {
			g.rot = 0
		}
	}

	if g.rot != prevRot || len(g.shape) == 0 {
		// regenerate vertices from shape
		g.shape = make([]Point, len(shapeSrc))
		cosRot := float32(math.Cos(g.rot))
		sinRot := float32(math.Sin(g.rot))
		for i, v := range shapeSrc {
			g.shape[i] = Point{
				X: scale*(v.X*cosRot-v.Y*sinRot) + g.center.X,
				Y: scale*(v.X*sinRot+v.Y*cosRot) + g.center.Y,
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawShape(screen, g.shape, 1, color.White)

	// msg := fmt.Sprintf("TPS: %0.2f\nRot: %.3f", ebiten.ActualTPS(), g.rot)
	msg := "Spin the dial with left and right arrows"
	if math.Abs(g.rot-g.rotGoal) < 0.01 {
		msg += "\n\nCLICK!"
	}
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("spin click (friday 01)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
