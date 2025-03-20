package main

import (
	"github.com/gdamore/tcell/v2"
)

type gridItem struct {
	Item        *Box
	Row, Column int
}

type Grid struct {
	*Box

	items []*gridItem
	rows  []int
}

func NewGrid() *Grid {
	grid := &Grid{}
	grid.Box = NewBox()
	return grid
}

func (gr *Grid) SetRows(rows ...int) *Grid {
	gr.rows = rows
	return gr
}

func (gr *Grid) Draw(screen tcell.Screen) {
	gr.Box.Draw(screen)
}
