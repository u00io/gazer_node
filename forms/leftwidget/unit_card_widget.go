package leftwidget

import (
	"fmt"

	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type UnitCardWidget struct {
	ui.Widget
	id string

	selected bool
	OnClick  func(clickedCategory string)
}

func NewUnitCardWidget(id string) *UnitCardWidget {
	var c UnitCardWidget
	c.InitWidget()

	c.SetPanelPadding(1)
	c.SetElevation(1)
	c.SetAutoFillBackground(true)
	c.SetLayout(`
	<row padding="0" onclick="Click" cursor="pointer">
		<frame padding="1" role="primary"  cursor="pointer"/>
		<column padding="0"  cursor="pointer">
			<label id="lblName" text="Unit"  onclick="Click"  cursor="pointer"/>
			<label id="lblUnitId" onclick="Click" cursor="pointer"/>
			<label id="lblValue" onclick="Click"  cursor="pointer"/>
			<hspacer  cursor="pointer"/>
		</column>
	</row>
	`, &c, nil)

	c.id = id

	lblName := c.FindWidgetByName("lblName").(*ui.Label)
	lblName.SetText(id)

	/*c.lblName = ui.NewLabel("")
	//c.lblName.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	c.lblName.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lblName.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	c.AddWidgetOnGrid(c.lblName, 0, 0)

	unitIdShort := ""
	unit := config.UnitById(id)
	if unit != nil {
		unitIdShort = unit.PublicKey
	}

	c.lblUnitId = ui.NewLabel(unitIdShort)
	c.lblUnitId.SetForegroundColor(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	c.lblUnitId.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lblUnitId.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	c.AddWidgetOnGrid(c.lblUnitId, 1, 0)

	c.lblValue = ui.NewLabel("123")
	c.lblValue.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.lblValue.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button == nuimouse.MouseButtonLeft {
			c.Click()
		}
		return true
	})
	c.AddWidgetOnGrid(c.lblValue, 2, 0)
	*/

	c.SetYExpandable(false)
	c.SetMinWidth(300)
	c.SetMinHeight(120)
	c.SetMaxHeight(120)
	c.SetSelected(false)
	c.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		fmt.Println("UnitCardWidget clicked:", c.id)
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
		c.SetRole("primary")

	} else {
		c.SetRole("")
	}
}

func (c *UnitCardWidget) IsSelected() bool {
	return c.selected
}

func (c *UnitCardWidget) UpdateData() {
	unitConfig := config.UnitById(c.id)
	if unitConfig == nil {
		return
	}

	lblName := c.FindWidgetByName("lblName").(*ui.Label)
	lblUnitId := c.FindWidgetByName("lblUnitId").(*ui.Label)
	lblValue := c.FindWidgetByName("lblValue").(*ui.Label)
	lblName.SetText(unitConfig.GetParameterString("0000_00_name_str", unitConfig.Type))
	lblUnitId.SetText(unitConfig.PublicKey)
	lblValue.SetText(system.Instance.GetUnitDefaultItemValue(c.id))
}
