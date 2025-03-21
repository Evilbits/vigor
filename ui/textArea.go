package ui

type TextArea struct {
	*Box
}

func NewTextArea() *TextArea {
	textArea := &TextArea{}
	textArea.Box = NewBox()

	return textArea
}

func (ta *TextArea) Draw(screen *Screen) {
	ta.Box.Draw(screen)
}

func (ta *TextArea) AddText(text string) {
	ta.Box.Text = text
}
