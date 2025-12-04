package unit102process

import (
	"time"

	"github.com/prometheus/procfs"
	"github.com/u00io/gazer_node/unit/unit000base"
)

type Unit102Process struct {
	unit000base.Unit

	dtOperationTime time.Time
	proc            procfs.Proc

	lastCpuValid bool
	lastCpuValue float64
	lastCpuTime  time.Time

	configProcessIdActive   bool
	configProcessId         int
	configProcessNameActive bool
	configProcessName       string

	lastConfigProcessId   int
	lastConfigProcessName string

	actualProcessId   int
	actualProcessName string

	lastKernelTimeMs     int64
	lastUserTimeMs       int64
	lastReadProcessTimes time.Time
}

type item struct {
	key   string
	name  string
	value string
	uom   string
}

func New() unit000base.IUnit {
	var c Unit102Process
	c.SetType("unit103storage")
	c.Init(&c)

	c.Config().SetParameterString("0000_00_name_str", "Process")

	c.Config().SetParameterString("0102_00_process_name_str", "explorer.exe")
	c.Config().SetParameterInt64("0102_01_process_id_int", 0)

	c.dtOperationTime = time.Now().UTC()
	c.actualProcessId = int(-1)
	c.lastCpuValid = false
	c.lastCpuValue = float64(0)
	c.lastCpuTime = time.Now()

	c.lastKernelTimeMs = int64(0)
	c.lastUserTimeMs = int64(0)
	c.lastReadProcessTimes = time.Now().UTC()

	return &c
}

type LastCounters struct {
	DT            time.Time
	TotalIn       uint64
	TotalOut      uint64
	TotalInBytes  uint64
	TotalOutBytes uint64
}

func (c *Unit102Process) Tick() {
	c.configProcessName = c.Config().GetParameterString("0102_00_process_name_str", "")
	c.configProcessId = int(c.Config().GetParameterInt64("0102_01_process_id_int", 0))

	if c.configProcessId > 0 {
		c.configProcessIdActive = true
	} else {
		c.configProcessIdActive = false
	}
	if c.configProcessName != "" {
		c.configProcessNameActive = true
	} else {
		c.configProcessNameActive = false
	}

	if c.lastConfigProcessId != c.configProcessId || c.lastConfigProcessName != c.configProcessName {
		c.actualProcessId = -1
		c.lastConfigProcessId = c.configProcessId
		c.lastConfigProcessName = c.configProcessName
	}

	items := c.processTick()
	for _, it := range items {
		c.SetValue(it.key, it.name, it.value, it.uom)
	}
}
