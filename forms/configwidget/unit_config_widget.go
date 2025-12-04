package configwidget

import (
	"fmt"
	"sort"

	"github.com/u00io/gazer_node/config"
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
	c.AddWidgetOnGrid(c.lvItems, 0, 0)
	return &c
}

func (c *UnitConfigWidget) SetUnitType(unitType string, parameters map[string]string) {
	// Load properties from unit type
	c.unitType = unitType
	type Item struct {
		key   string
		value string
	}
	var items []Item

	for k, v := range parameters {
		items = append(items, Item{key: k, value: v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].key < items[j].key
	})

	c.lvItems.SetRowCount(len(items))

	for row := 0; row < len(items); row++ {
		key := items[row].key
		propName := config.PropName(key)
		c.lvItems.SetCellData2(row, 0, key)
		c.lvItems.SetCellText2(row, 0, propName)
		c.lvItems.SetCellText2(row, 1, items[row].value)
	}
}

func (c *UnitConfigWidget) GetParameters() map[string]string {
	parameters := make(map[string]string)
	for row := 0; row < c.lvItems.RowCount(); row++ {
		key := fmt.Sprint(c.lvItems.GetCellData2(row, 0))
		value := c.lvItems.GetCellText2(row, 1)
		if key != "" {
			parameters[key] = value
		}
	}
	return parameters
}
