package main

import (
	"log"
	"os"

	"github.com/evilbits/vigor/ui"
	"github.com/gdamore/tcell/v2"
)

func startEventLoop(screen tcell.Screen, grid *ui.Grid, debug bool) {
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
			grid.Draw(screen, debug)
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}

func main() {
	debug := true
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}

	screen.Clear()

	rootGrid := ui.NewGrid()
	rootBox := ui.NewBox(screen)
	rootBox.AddMultilineText([]string{"Hello world!", "Second line"})

	rootBoxTwo := ui.NewBox(screen)
	rootBoxTwo.SetBackgroundColor("red")
	rootBoxTwo.AddText("Position")

	rootGrid.
		SetRows(0, 1).
		AddItem(rootBox).
		AddItem(rootBoxTwo).
		Draw(screen, debug)

	startEventLoop(screen, rootGrid, debug)
}
