package ui

type gridItem struct {
	Item Drawable
}

type Grid struct {
	items []*gridItem
	// Row array containing metadata regarding row size
	rows []int
}

func NewGrid() *Grid {
	grid := &Grid{}
	return grid
}

// SetRows acts as a way to change the default behaviour of rows within the grid.
// The integer value that is chosen indicates the height of the row in relation to the amount of rows the Box can display.
//
// A value of 0 indicates that a row should take up any remaining space that it can.
// A value of 1 indicates that the row should take up 1 height unit.
func (gr *Grid) SetRows(rows ...int) *Grid {
	gr.rows = rows
	return gr
}

// Add an item to the grid. Item order matters as we expect gr.rows[0] to be filled by gr.items[0].
// Adding an item without a corresponding row will lead to the item not being rendered.
func (gr *Grid) AddItem(item Drawable) *Grid {
	gr.items = append(gr.items, &gridItem{
		Item: item,
	})
	return gr
}

func (gr *Grid) GetItem(idx int) *gridItem {
	return gr.items[idx]
}

func (gr *Grid) Draw(screen *Screen) {
	screenWidth, screenHeight := screen.Size()

	// Start rendering at yPos 0
	nextYPos := 0

	// Calculate how many items should fill remaining space
	fillItems := 0
	for _, rowHeight := range gr.rows {
		if rowHeight == 0 {
			fillItems += 1
		}
	}

	for idx, gridItem := range gr.items {
		// Don't render items that don't have a row
		if idx >= len(gr.rows) {
			return
		}

		// If this row height is 0 calculate the height of this element
		calculatedHeight := gr.rows[idx]

		if calculatedHeight == 0 {
			// Calculate height based on how many items we have that require unknown yHeight
			calculatedHeight = screenHeight / fillItems
			// Remove space used by elements below
			for _, rowHeight := range gr.rows[idx:] {
				calculatedHeight -= rowHeight
			}
		}

		item := gridItem.Item
		item.SetRect(screenWidth, calculatedHeight, nextYPos, 0)
		item.Draw(screen)

		// Move down yPos so next item can render from there
		nextYPos += calculatedHeight
	}
}
