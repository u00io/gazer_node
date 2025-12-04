package addunitwidget

import (
	"image/color"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type UnitTypeWidget struct {
	ui.Widget
	unitTypeName        string
	unitTypeDisplayName string

	lbl *ui.Label

	OnClick func(clickedItem string)

	selected bool
}

func NewUnitTypeWidget(unitTypeName string, unitTypeDisplayName string) *UnitTypeWidget {
	var c UnitTypeWidget
	c.InitWidget()
	c.SetAutoFillBackground(true)
	c.unitTypeName = unitTypeName
	c.unitTypeDisplayName = unitTypeDisplayName

	c.lbl = ui.NewLabel(unitTypeDisplayName)
	c.lbl.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lbl.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			if c.OnClick != nil {
				c.OnClick(c.unitTypeName)
			}
		}
		return true
	})

	c.AddWidgetOnGrid(c.lbl, 0, 0)
	c.SetMinHeight(75)
	c.SetMaxHeight(75)
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
		c.lbl.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	} else {
		c.SetBackgroundColor(color.RGBA{R: 40, G: 40, B: 40, A: 255})
		c.lbl.SetForegroundColor(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	}
}
