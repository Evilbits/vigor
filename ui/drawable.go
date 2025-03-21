package ui

// Drawable represents any component that can be drawn to the screen
type Drawable interface {
	Draw(screen Screen)
	SetRect(width, height, y, x int)
}
