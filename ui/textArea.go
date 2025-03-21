package ui

type TextArea struct {
	*Box

	cursorX, cursorY int
}

func NewTextArea() *TextArea {
	textArea := &TextArea{
		cursorX: 0,
		cursorY: 0,
	}
	textArea.Box = NewBox()

	return textArea
}

func (ta *TextArea) Draw(screen *Screen) {
	ta.Box.Draw(screen)

	screen.RenderCursor(ta.cursorX, ta.cursorY)
}

func (ta *TextArea) moveCursor(moveX int, moveY int, screen *Screen) {
	x, y := ta.GetXY()
	if (ta.cursorX+moveX < x) || (ta.cursorY+moveY < y) {
		return
	}
	width, height := ta.Box.GetSize()
	if (ta.cursorX+moveX >= width) || (ta.cursorY+moveY >= height) {
		return
	}
	prevX := ta.cursorX
	prevY := ta.cursorY
	ta.cursorX += moveX
	ta.cursorY += moveY
	screen.RenderCursorMove(ta.cursorX, ta.cursorY, prevX, prevY)
	screen.innerScreen.Show()
}

func (ta *TextArea) HandleKey(char rune, screen *Screen) {
	switch char {
	case 'h':
		ta.moveCursor(-1, 0, screen)
	case 'j':
		ta.moveCursor(0, 1, screen)
	case 'k':
		ta.moveCursor(0, -1, screen)
	case 'l':
		ta.moveCursor(1, 0, screen)
	}
}
