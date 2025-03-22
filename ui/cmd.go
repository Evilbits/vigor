package ui

type Cmd struct {
	*Box
}

func NewCmd() *Cmd {
	cmd := &Cmd{}
	cmd.Box = NewBox()
	return cmd
}

func (cmd *Cmd) Draw(screen *Screen) {
	cmd.AddText("CMD")
	cmd.Box.Draw(screen)
}
