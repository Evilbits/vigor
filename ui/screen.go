package ui

import (
	"github.com/gdamore/tcell/v2"
	"log"
)

type Screen struct {
	*Grid

	innerScreen tcell.Screen

	cursorColor string
	OnKeyPress  func(event *tcell.EventKey)
}

func NewScreen() *Screen {
	tScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	screen := &Screen{
		innerScreen: tScreen,
		cursorColor: "white",
	}
	return screen
}

func (screen *Screen) initTScreen() {
	if err := screen.innerScreen.Init(); err != nil {
		log.Fatal(err)
	}
}

func (screen *Screen) WriteChar(x int, y int, char rune, combining []rune, background tcell.Style) {
	screen.innerScreen.SetContent(x, y, char, nil, background)
}

func (screen *Screen) RenderCursor(x int, y int, cursorStyle tcell.CursorStyle) {
	screen.innerScreen.ShowCursor(x, y)
	screen.innerScreen.SetCursorStyle(cursorStyle)
}

func (screen *Screen) Size() (width int, height int) {
	return screen.innerScreen.Size()
}

func (screen *Screen) StartEventLoop(focusedArea *TextArea) {
	if screen.Grid == nil {
		log.Fatal("Screen must have a Grid to render to")
	}
	tScreen := screen.innerScreen

	screen.initTScreen()
	tScreen.Clear()

	for {
		event := tScreen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventResize:
			tScreen.Sync()
			screen.Grid.Draw(screen)
		case *tcell.EventKey:
			if screen.OnKeyPress == nil {
				panic("No event handler registered")
			}
			screen.OnKeyPress(event)
		}
		screen.Grid.Draw(screen)
		tScreen.Show()
	}
}

func (screen *Screen) Fini() {
	screen.innerScreen.Fini()
}
