package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Box struct {
	*Screen

	Text string

	x, y, width, height int

	innerX, innerY, innerWidth, innerHeight int

	backgroundColor tcell.Color
}

func NewBox(screen *Screen) *Box {
	box := &Box{
		width:           15,
		height:          15,
		backgroundColor: tcell.GetColor("grey"),
	}
	box.Screen = screen
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

func (b *Box) AddText(text string) *Box {
	b.Text = text
	return b
}

func (b *Box) AddMultilineText(textArr []string) *Box {
	compiledText := ""

	for _, text := range textArr {
		compiledText += LF.Add(text)
	}
	b.Text = compiledText
	return b
}
