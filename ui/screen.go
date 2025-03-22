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

func (screen *Screen) RenderCursor(x int, y int) {
	cursorColor := tcell.GetColor(screen.cursorColor)
	cursorStyle := tcell.StyleDefault.Background(cursorColor)
	primary, combining, _, _ := screen.innerScreen.GetContent(x, y)
	screen.innerScreen.SetContent(x, y, primary, combining, cursorStyle)
}

func (screen *Screen) RenderCursorMove(x int, y int, prevX int, prevY int) {
	cursorColor := tcell.GetColor(screen.cursorColor)
	cursorStyle := tcell.StyleDefault.Background(cursorColor)

	prevPrimary, prevCombining, prevStyle, _ := screen.innerScreen.GetContent(prevX, prevY)
	primary, combining, _, _ := screen.innerScreen.GetContent(x, y)

	// Update new cursor location and reset previous location
	screen.innerScreen.SetContent(x, y, primary, combining, cursorStyle)
	screen.innerScreen.SetContent(prevX, prevY, prevPrimary, prevCombining, prevStyle)
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
