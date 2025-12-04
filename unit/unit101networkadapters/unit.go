package unit101networkadapters

import (
	"net"
	"time"

	"github.com/u00io/gazer_node/unit/unit000base"
)

type Unit101NetworkAdapters struct {
	unit000base.Unit

	lastCounters          map[int]LastCounters
	addressesOfInterfaces map[int]string
}

type item struct {
	key   string
	name  string
	value string
	uom   string
}

func New() unit000base.IUnit {
	var c Unit101NetworkAdapters
	c.SetType("unit103storage")
	c.Init(&c)
	c.lastCounters = make(map[int]LastCounters)
	c.addressesOfInterfaces = make(map[int]string)

	c.Config().SetParameterString("0000_00_name_str", "Network Adapters")

	return &c
}

type LastCounters struct {
	DT            time.Time
	TotalIn       uint64
	TotalOut      uint64
	TotalInBytes  uint64
	TotalOutBytes uint64
}

func (c *Unit101NetworkAdapters) writeAddresses(ni net.Interface) {
	// Addresses
	addrs, err := ni.Addrs()
	if err == nil {
		addrsString := ""
		for _, a := range addrs {
			if len(addrsString) > 0 {
				addrsString += " "
			}
			addrsString += a.String()
		}
		if c.addressesOfInterfaces[ni.Index] != addrsString {
			c.addressesOfInterfaces[ni.Index] = addrsString
			// c.SetString(ni.Name+"/Addresses", addrsString, "-")
		}
	}
}

func (c *Unit101NetworkAdapters) Tick() {
	items := c.processTick()
	for _, it := range items {
		c.SetValue(it.key, it.name, it.value, it.uom)
	}
}
