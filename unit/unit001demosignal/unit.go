package unit001demosignal

import (
	"math"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/u00io/gazer_node/unit/unit000base"
)

type Unit01FileContent struct {
	unit000base.Unit

	counter int
}

func New() unit000base.IUnit {
	var c Unit01FileContent
	c.SetType("unit001demosignal")
	c.Init(&c)

	c.SetParameterFloat64("100_offset_num", 0.0)

	return &c
}

func (c *Unit01FileContent) Tick() {
	offset := c.GetParameterFloat64("100_offset_num", 0.0)

	demoData := ""
	//demoData += time.Now().Format("15:04:05")
	//rnd := rand.Int31() % 100
	sinValue := math.Sin(float64(time.Now().Unix()%60)/60.0*2.0*math.Pi)*100 + 100
	// add slow sin wave
	sinValue += math.Sin(float64(time.Now().Unix()%300)/300.0*2.0*math.Pi)*50 + 50
	// add fast sin wave
	sinValue += math.Sin(float64(time.Now().Unix()%10)/10.0*2.0*math.Pi)*20 + 20
	// add some noise
	sinValue += (rand.Float64() - 0.5) * 10

	// add offset
	sinValue += float64(offset)

	demoData = strconv.FormatFloat(sinValue, 'f', 1, 64)

	c.SetValue("value", demoData)
	c.counter++
}
