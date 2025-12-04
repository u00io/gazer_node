package addunit

import (
	"github.com/u00io/nuiforms/ui"
)

type UnitConfigWidget struct {
	ui.Widget

	unitType string

	lvItems *ui.Table
}

func NewUnitConfigWidget() *UnitConfigWidget {
	var c UnitConfigWidget
	c.InitWidget()
	c.SetAllowScroll(true, true)
	c.SetYExpandable(true)
	c.lvItems = ui.NewTable()
	c.lvItems.SetColumnCount(2)
	c.lvItems.SetColumnWidth(0, 150)
	c.lvItems.SetColumnWidth(1, 200)
	c.lvItems.SetColumnName(0, "Property")
	c.lvItems.SetColumnName(1, "Value")
	c.lvItems.SetEditTriggerEnter(true)
	c.lvItems.SetEditTriggerF2(true)
	c.lvItems.SetEditTriggerDoubleClick(true)
	c.lvItems.SetRowCount(10)
	c.AddWidgetOnGrid(c.lvItems, 0, 0)
	return &c
}

func (c *UnitConfigWidget) SetUnitType(unitType string) {
	// Load properties from unit type
	c.unitType = unitType
}

func (c *UnitConfigWidget) GetParameters() map[string]string {
	parameters := make(map[string]string)
	for row := 0; row < c.lvItems.RowCount(); row++ {
		key := c.lvItems.GetCellText2(row, 0)
		value := c.lvItems.GetCellText2(row, 1)
		if key != "" {
			parameters[key] = value
		}
	}
	return parameters
}
