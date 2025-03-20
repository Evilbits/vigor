package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func startEventLoop(screen tcell.Screen) {
	quit := func() {
		screen.Fini()
		os.Exit(0)
	}

	for {
		screen.Show()

		event := screen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}

	screen.Clear()

	root := NewGrid()
	root.Draw(screen)

	startEventLoop(screen)
}
