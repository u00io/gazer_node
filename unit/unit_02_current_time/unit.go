package unit02currenttime

import (
	"time"

	unit00base "github.com/u00io/gazer_node/unit/unit_00_base"
)

type Unit02CurrentTime struct {
	unit00base.Unit
}

func New() unit00base.IUnit {
	var c Unit02CurrentTime
	c.SetType("unit02currenttime")
	c.Init(&c)
	return &c
}

func (c *Unit02CurrentTime) Tick() {
	//logger.Println("Unit02CurrentTime Tick")
	c.SetValue("value", "Current Time: "+time.Now().Format("2006-01-02 15:04:05"))
}
