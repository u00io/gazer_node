package addunitwidget

import (
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type CategoryWidget struct {
	ui.Widget
	lbl          *ui.Label
	categoryName string
	selected     bool
	OnClick      func(clickedCategory string)
}

func NewCategoryWidget(categoryName string) *CategoryWidget {
	var c CategoryWidget
	c.InitWidget()
	c.SetAutoFillBackground(true)
	c.categoryName = categoryName
	c.lbl = ui.NewLabel(categoryName)
	c.lbl.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lbl.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			if c.OnClick != nil {
				c.OnClick(c.categoryName)
			}
		}
		return true
	})
	c.AddWidgetOnGrid(c.lbl, 0, 0)
	c.SetYExpandable(false)
	c.SetMinHeight(30)
	c.SetMaxHeight(30)
	c.SetSelected(false)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			if c.OnClick != nil {
				c.OnClick(c.categoryName)
			}
		}
		return true
	})
	return &c
}

func (c *CategoryWidget) SetSelected(selected bool) {
	c.selected = selected
	if selected {
		c.SetBackgroundColor(color.RGBA{R: 60, G: 60, B: 60, A: 255})
		c.lbl.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	} else {
		c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 40, A: 255})
		c.lbl.SetForegroundColor(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	}
}
