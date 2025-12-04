package unit104memory

import (
	"strconv"

	"github.com/shirou/gopsutil/mem"
	"github.com/u00io/gazer_node/unit/unit000base"
)

type Unit104Memory struct {
	unit000base.Unit

	counter int
}

func New() unit000base.IUnit {
	var c Unit104Memory
	c.SetType("unit104memory")
	c.Init(&c)

	return &c
}

func (c *Unit104Memory) Tick() {
	v, _ := mem.VirtualMemory()
	percents := (float64(v.Used) / float64(v.Total)) * 100.0
	c.SetValue("/", "Mem Used Percent", strconv.FormatFloat(percents, 'f', 2, 64), "%")
	c.SetValue("/total", "Mem Total", strconv.FormatUint(v.Total/1048576, 10), "MB")
	c.SetValue("/used", "Mem Used", strconv.FormatUint(v.Used/1048576, 10), "MB")
	c.SetValue("/free", "Mem Free", strconv.FormatUint(v.Free/1048576, 10), "MB")
	c.counter++
}
