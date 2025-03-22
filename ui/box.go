package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Box struct {
	Text string

	x, y, width, height int

	backgroundColor tcell.Color
}

func NewBox() *Box {
	box := &Box{
		width:  15,
		height: 15,
	}
	return box
}

// If background color isn't found it will be set to tcell.ColorDefault
func (b *Box) SetBackgroundColor(color string) {
	mappedColor := tcell.GetColor(color)
	b.backgroundColor = mappedColor
}

func (b *Box) SetRect(width int, height int, y int, x int) {
	b.width = width
	b.height = height
	b.x = x
	b.y = y
}

func (b *Box) GetSize() (width int, height int) {
	return b.width, b.height
}

func (b *Box) GetXY() (x int, y int) {
	return b.x, b.y
}

func (b *Box) Draw(screen *Screen) {
	background := tcell.StyleDefault.Background(b.backgroundColor)
	for y := b.y; y < b.y+b.height; y++ {
		for x := b.x; x < b.x+b.width; x++ {
			screen.WriteChar(x, y, ' ', nil, background)
		}
	}

	if b.Text != "" {
		y := b.y
		x := b.x

		for _, char := range b.Text {
			if LineFeed(char) == LF {
				y++
				x = 0
				continue
			}
			screen.WriteChar(x, y, char, nil, background)
			x++
		}
	}
}

func (b *Box) AddText(text string) {
	b.Text = text
}

func (b *Box) AppendRune(char rune) {
	b.Text += string(char)
}
