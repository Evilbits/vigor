package ui

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	innerScreen tcell.Screen
}

func NewScreen() *Screen {
	tScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	screen := &Screen{
		innerScreen: tScreen,
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

func (screen *Screen) Size() (width int, height int) {
	return screen.innerScreen.Size()
}

func (screen *Screen) StartEventLoop(grid *Grid) {
	tScreen := screen.innerScreen

	screen.initTScreen()
	tScreen.Clear()

	quit := func() {
		tScreen.Fini()
		os.Exit(0)
	}

	for {
		tScreen.Show()

		event := tScreen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventResize:
			tScreen.Sync()
			grid.Draw(screen)
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}
