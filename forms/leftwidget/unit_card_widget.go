package leftwidget

import (
	"image/color"

	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type UnitCardWidget struct {
	ui.Widget
	id string

	categoryName        string
	categoryDisplayName string

	selected bool
	OnClick  func(clickedCategory string)

	lblCategory *ui.Label
	lblUnitId   *ui.Label
	lblValue    *ui.Label
}

func NewUnitCardWidget(categoryName string, categoryDisplayName string, id string) *UnitCardWidget {
	var c UnitCardWidget
	c.InitWidget()

	c.SetPanelPadding(1)
	c.SetAutoFillBackground(true)

	c.id = id
	c.categoryName = categoryName
	c.categoryDisplayName = categoryDisplayName

	c.lblCategory = ui.NewLabel(categoryDisplayName)
	c.lblCategory.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lblCategory.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	c.AddWidgetOnGrid(c.lblCategory, 0, 0)

	unitIdShort := ""
	unit := config.UnitById(id)
	if unit != nil {
		unitIdShort = unit.PublicKey
	}
	// input: 1234-------3240
	// format: 0x1234...3240
	if len(id) == 64 {
		unitIdShort = "0x" + id[:4] + "..." + id[len(id)-4:]
	}

	c.lblUnitId = ui.NewLabel(unitIdShort)
	c.lblUnitId.SetForegroundColor(color.RGBA{R: 150, G: 150, B: 250, A: 255})
	c.lblUnitId.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lblUnitId.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	c.AddWidgetOnGrid(c.lblUnitId, 1, 0)

	c.lblValue = ui.NewLabel("123")
	c.lblValue.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	c.AddWidgetOnGrid(c.lblValue, 2, 0)

	c.SetYExpandable(false)
	c.SetMinWidth(300)
	c.SetMinHeight(120)
	c.SetMaxHeight(120)
	c.SetSelected(false)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	return &c
}

func (c *UnitCardWidget) Click() {
	if c.OnClick != nil {
		c.OnClick(c.id)
	}
}

func (c *UnitCardWidget) SetSelected(selected bool) {
	c.selected = selected
	if selected {
		backColor := c.BackgroundColorAccent2()
		c.SetBackgroundColor(backColor)
		c.lblCategory.SetBackgroundColor(backColor)
		c.lblUnitId.SetBackgroundColor(backColor)
	} else {
		backColor := c.BackgroundColorAccent1()
		c.SetBackgroundColor(backColor)
		c.lblCategory.SetBackgroundColor(backColor)
		c.lblUnitId.SetBackgroundColor(backColor)
	}
}

func (c *UnitCardWidget) IsSelected() bool {
	return c.selected
}

func (c *UnitCardWidget) UpdateData() {
	c.lblValue.SetText(system.Instance.GetUnitDefaultItemValue(c.id))
}
