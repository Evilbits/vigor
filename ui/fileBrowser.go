package ui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

type FileBrowser struct {
	*Box

	Files   []os.DirEntry
	CurrDir string

	cursorY int
}

func NewFileBrowser(currDir string) *FileBrowser {
	fileBrowser := &FileBrowser{
		CurrDir: currDir,
		cursorY: 0,
	}
	fileBrowser.Box = NewBox()
	return fileBrowser
}

func (fb *FileBrowser) MoveCursor(moveY int) {
	fb.cursorY += moveY
	if fb.cursorY < 0 {
		fb.cursorY = 0
	}
	maxYPos := len(fb.Files) - 1
	if fb.cursorY > maxYPos {
		fb.cursorY = maxYPos
	}
}

func (fb *FileBrowser) GetCurrentFile() os.DirEntry {
	return fb.Files[fb.cursorY]
}

func (fb *FileBrowser) Update() {
	entries, err := os.ReadDir(fb.CurrDir)
	if err != nil {
		panic(err)
	}

	// Order so that directories come first
	var dirs []os.DirEntry
	var files []os.DirEntry

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}
	allEntries := append(dirs, files...)

	fb.Files = allEntries
}

func (fb *FileBrowser) Draw(screen *Screen) {
	if len(fb.Files) == 0 {
		fb.Update()
	}

	fb.Box.AddText(fb.buildTextContent())
	fb.Box.Draw(screen)
	screen.RenderCursor(0, fb.cursorY, tcell.CursorStyleDefault)
}

func (fb *FileBrowser) buildTextContent() string {
	var output string
	for _, dir := range fb.Files {
		output += fb.renderEntry(dir)
	}
	return output
}

func (fb *FileBrowser) renderEntry(entry os.DirEntry) string {
	var output string
	output += entry.Name()
	if entry.IsDir() {
		output += "/"
	}
	output += fmt.Sprint(LF)
	return output
}
