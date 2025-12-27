package tview

import (
	"github.com/gdamore/tcell/v2"
)

type FormModal struct {
	// The form embedded in the modal's frame.
	*Form
	labelWidth  int
	itemsWidth  int
	itemsHeight int
}

// NewFormModal returns a new [FormModal] message window.
func NewFormModal(cb func(form *Form)) *FormModal {
	form := NewForm().
		SetButtonsAlign(AlignCenter).
		SetButtonBackgroundColor(Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(Styles.PrimaryTextColor)
	form.SetBorderPadding(0, 0, 0, 0).SetBorder(true)
	cb(form)

	return &FormModal{Form: form}
}

// SetDirtySize sets itemsWidth, itemsHeight and labelWidth as dirty, so it will be recalculated on next redraw.
// When you add/remove/change a field, call this.
func (f *FormModal) SetDirtySize() {
	f.labelWidth = 0
}

func (f *FormModal) recalcSize() {
	f.itemsWidth = 0
	f.itemsHeight = 3
	f.labelWidth = f.calcMaxLabelWidth()

	for _, item := range f.Form.items {
		itemWidth, itemHeight := item.GetFieldWidth(), item.GetFieldHeight()
		f.itemsWidth = max(f.itemsWidth, itemWidth)
		f.itemsHeight += itemHeight + 1
	}

	// resize input fields to the max width
	for _, item := range f.Form.items {
		inputField, ok := item.(*InputField)
		if !ok {
			continue
		}

		if itemWidth := item.GetFieldWidth(); itemWidth < f.itemsWidth {
			inputField.SetFieldWidth(f.itemsWidth)
		}
	}

	f.itemsWidth += 2 // padding

	buttonsWidth := -1
	for _, b := range f.buttons {
		buttonsWidth += b.width + 1
	}

	if buttonsWidth > f.itemsWidth+f.labelWidth {
		f.itemsWidth = buttonsWidth - f.labelWidth
	}
}

// Draw draws this primitive onto the screen.
func (f *FormModal) Draw(screen tcell.Screen) {
	if f.labelWidth == 0 {
		f.recalcSize()
	}

	screenWidth, screenHeight := screen.Size()
	width := min(f.itemsWidth+f.labelWidth, screenWidth)
	height := min(f.itemsHeight, screenHeight)

	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	f.SetRect(x, y, width, height)
	f.Form.Draw(screen)
}
