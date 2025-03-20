package main

import (
	"github.com/gdamore/tcell/v2"
)

type gridItem struct {
	Item *Box
	Row  int
}

type Grid struct {
	items []*gridItem
	rows  []int
}

func NewGrid() *Grid {
	grid := &Grid{}
	return grid
}

func (gr *Grid) SetRows(rows ...int) *Grid {
	gr.rows = rows
	return gr
}

func (gr *Grid) AddItem(item *Box, row int) *Grid {
	gr.items = append(gr.items, &gridItem{
		Item: item,
		Row:  row,
	})
	return gr
}

func (gr *Grid) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()

	// Split our view up into sections assuming each section is the same height
	itemHeight := screenHeight / len(gr.items)

	for idx, gridItem := range gr.items {
		item := gridItem.Item

		// Calculate y position of this item based on the total items we are rendering
		yPosition := idx * itemHeight
		item.SetRect(screenWidth, itemHeight, yPosition, 0)
		item.Draw(screen)
	}
}
