package ui

type Drawable interface {
	Draw(screen *Screen)
	SetRect(width int, height int, y int, x int)
	AddText(text string)
}
