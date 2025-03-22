package ui

import "github.com/gdamore/tcell/v2"

type Drawable interface {
	Draw(screen *Screen)
	SetRect(width int, height int, y int, x int)
	AddText(text string)
	// TODO: Find a way to remove HandleKey()
	HandleKey(event *tcell.EventKey)
}
