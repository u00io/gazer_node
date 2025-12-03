package addpage

import (
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type UnitTypeWidget struct {
	ui.Widget
	unitTypeName string

	OnClick func(clickedItem string)

	selected bool
}

func NewUnitTypeWidget(unitTypeName string) *UnitTypeWidget {
	var c UnitTypeWidget
	c.InitWidget()
	c.SetAutoFillBackground(true)
	c.unitTypeName = unitTypeName

	lbl := ui.NewLabel(unitTypeName)
	lbl.SetMouseCursor(nuimouse.MouseCursorPointer)
	lbl.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			if c.OnClick != nil {
				c.OnClick(c.unitTypeName)
			}
		}
		return true
	})

	c.AddWidgetOnGrid(lbl, 0, 0)
	c.SetMinHeight(200)
	c.SetMaxHeight(200)
	c.SetSelected(false)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)

	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			if c.OnClick != nil {
				c.OnClick(c.unitTypeName)
			}
		}
		return true
	})

	return &c
}

func (c *UnitTypeWidget) SetSelected(selected bool) {
	c.selected = selected
	if selected {
		c.SetBackgroundColor(color.RGBA{R: 60, G: 60, B: 60, A: 255})
	} else {
		c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 40, A: 255})
	}
}
