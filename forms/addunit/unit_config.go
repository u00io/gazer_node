package addunit

import "github.com/u00io/nuiforms/ui"

type UnitConfig struct {
	ui.Widget

	lvItems *ui.Table
}

func NewUnitConfig() *UnitConfig {
	var c UnitConfig
	c.InitWidget()
	c.SetAllowScroll(true, true)
	c.SetYExpandable(true)
	c.lvItems = ui.NewTable()
	c.lvItems.SetColumnCount(2)
	c.lvItems.SetColumnWidth(0, 200)
	c.lvItems.SetColumnWidth(1, 300)
	c.lvItems.SetColumnName(0, "Property")
	c.lvItems.SetColumnName(1, "Value")
	c.lvItems.SetEditTriggerEnter(true)
	c.lvItems.SetEditTriggerF2(true)
	c.lvItems.SetEditTriggerDoubleClick(true)
	c.AddWidgetOnGrid(c.lvItems, 0, 0)
	return &c
}
