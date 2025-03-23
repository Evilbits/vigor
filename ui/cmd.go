package ui

import (
	"fmt"
)

type Cmd struct {
	*Box

	CurrentCommand string
	CommandMode    bool
}

func NewCmd() *Cmd {
	cmd := &Cmd{
		CommandMode: false,
	}
	cmd.Box = NewBox()
	return cmd
}

func (cmd *Cmd) Draw(screen *Screen) {
	if cmd.CommandMode {
		cmd.AddText(fmt.Sprintf("> %s", cmd.CurrentCommand))
	}
	cmd.Box.Draw(screen)
}

func (cmd *Cmd) StartCommandMode() {
	cmd.SetBackgroundColor("")
	cmd.CommandMode = true
	cmd.AddText("")
}

func (cmd *Cmd) ExitCommandMode() {
	cmd.CommandMode = false
	cmd.ResetCurrentCommand()
}

func (cmd *Cmd) AppendRuneToCurrentCommand(char rune) {
	cmd.CurrentCommand += string(char)
}

func (cmd *Cmd) DeleteLastCharFromCommand() {
	if len(cmd.CurrentCommand) == 0 {
		return
	}
	cmd.CurrentCommand = cmd.CurrentCommand[:len(cmd.CurrentCommand)-1]
}

func (cmd *Cmd) ResetCurrentCommand() {
	cmd.CurrentCommand = ""
}

func (cmd *Cmd) SetError(error string) {
	cmd.AddText(error)
	cmd.SetBackgroundColor("red")
}
