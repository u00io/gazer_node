package unit103storage

import (
	"github.com/u00io/gazer_node/unit/unit000base"
)

type Unit103Storage struct {
	unit000base.Unit

	disk string
}

type item struct {
	key   string
	name  string
	value string
	uom   string
}

func New() unit000base.IUnit {
	var c Unit103Storage
	c.SetType("unit103storage")
	c.Init(&c)

	c.Config().SetParameterString("0000_00_name_str", "Storage")

	return &c
}

func (c *Unit103Storage) Tick() {
	items := c.processTick()
	for _, it := range items {
		c.SetValue(it.key, it.name, it.value, it.uom)
	}
}
