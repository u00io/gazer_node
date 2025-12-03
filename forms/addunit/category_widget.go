package addunit

import (
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type CategoryWidget struct {
	ui.Widget
	categoryName string
	selected     bool
	OnClick      func(clickedCategory string)
}

func NewCategoryWidget(categoryName string) *CategoryWidget {
	var c CategoryWidget
	c.InitWidget()
	c.SetAutoFillBackground(true)
	c.categoryName = categoryName
	lbl := ui.NewLabel(categoryName)
	lbl.SetMouseCursor(nuimouse.MouseCursorPointer)
	lbl.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			if c.OnClick != nil {
				c.OnClick(c.categoryName)
			}
		}
		return true
	})
	c.AddWidgetOnGrid(lbl, 0, 0)
	c.SetYExpandable(false)
	c.SetMinHeight(90)
	c.SetMaxHeight(90)
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
	} else {
		c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 40, A: 255})
	}
}
