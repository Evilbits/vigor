package main

import (
	"github.com/gdamore/tcell/v2"
)

type Box struct {
	x, y, width, height int

	innerX, innerY, innerWidth, innerHeight int

	backgroundColor tcell.Color
}

func NewBox() *Box {
	box := &Box{
		width:           15,
		height:          15,
		backgroundColor: tcell.GetColor("grey"),
	}
	return box
}

func (b *Box) Draw(screen tcell.Screen) {
	background := tcell.StyleDefault.Background(b.backgroundColor)
	for y := b.y; y < b.y+b.height; y++ {
		for x := b.x; x < b.x+b.width; x++ {
			screen.SetContent(x, y, ' ', nil, background)
		}
	}
	screen.SetContent(1, 1, 'H', nil, background)
	screen.SetContent(2, 1, 'e', nil, background)
}
