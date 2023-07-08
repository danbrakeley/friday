package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Shape struct {
	Points      []Vec2D
	Xfm         Transform2D
	Color       color.Color
	StrokeWidth float32
}

func (s *Shape) Draw(screen *ebiten.Image) {
	drawShape(screen, s.Points, s.StrokeWidth, s.Color)
}
