package unit000base

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/gomisc/logger"
)

type Unit struct {
	id      string
	mtx     sync.Mutex
	key     *utils.Key
	tp      string
	started bool
	stoping bool
	config  map[string]string
	values  map[string]string
	iUnit   IUnit
}

type IUnit interface {
	Start()
	Stop()

	SetId(id string)
	GetId() string

	GetKey() *utils.Key
	SetKey(key *utils.Key)

	SetConfig(config map[string]string)
	GetConfig() map[string]string

	GetParameterBool(key string, defaultValue bool) bool
	GetParameterInt64(key string, defaultValue int64) int64
	GetParameterFloat64(key string, defaultValue float64) float64
	GetParameterString(key string, defaultValue string) string

	SetParameterBool(key string, value bool)
	SetParameterInt64(key string, value int64)
	SetParameterFloat64(key string, value float64)
	SetParameterString(key string, value string)

	GetType() string
	GetValue(key string) string
	SetValue(key, value string)
	Tick()
}

func (c *Unit) Init(iUnit IUnit) {
	c.config = make(map[string]string)
	c.values = make(map[string]string)
	c.iUnit = iUnit
}

func (c *Unit) GetId() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.id
}

func (c *Unit) SetId(id string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.id = id
}

func (c *Unit) GetKey() *utils.Key {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.key
}

func (c *Unit) SetKey(key *utils.Key) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.key = key
}

func (c *Unit) SetType(tp string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.tp = tp
}

func (c *Unit) GetType() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.tp
}

func (c *Unit) SetConfig(config map[string]string) {
	c.mtx.Lock()
	c.config = config
	c.mtx.Unlock()
}

func (c *Unit) GetConfig() map[string]string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.config
}

func (c *Unit) GetParameterBool(key string, defaultValue bool) bool {
	valueStr := c.GetParameterString(key, fmt.Sprint(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Unit) GetParameterInt64(key string, defaultValue int64) int64 {
	valueStr := c.GetParameterString(key, fmt.Sprint(defaultValue))
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Unit) GetParameterFloat64(key string, defaultValue float64) float64 {
	valueStr := c.GetParameterString(key, fmt.Sprint(defaultValue))
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Unit) GetParameterString(key string, defaultValue string) string {
	c.mtx.Lock()
	value, exists := c.config[key]
	c.mtx.Unlock()
	if !exists {
		return defaultValue
	}
	return value
}

func (c *Unit) SetParameterBool(key string, value bool) {
	c.SetParameterString(key, fmt.Sprint(value))
}

func (c *Unit) SetParameterInt64(key string, value int64) {
	c.SetParameterString(key, fmt.Sprint(value))
}

func (c *Unit) SetParameterFloat64(key string, value float64) {
	c.SetParameterString(key, fmt.Sprint(value))
}

func (c *Unit) SetParameterString(key string, value string) {
	c.mtx.Lock()
	c.config[key] = value
	c.mtx.Unlock()
}

func (c *Unit) GetValue(key string) string {
	c.mtx.Lock()
	value, exists := c.values[key]
	c.mtx.Unlock()
	if !exists {
		return ""
	}
	return value
}

func (c *Unit) SetValue(key, value string) {
	c.mtx.Lock()
	c.values[key] = value
	c.mtx.Unlock()
}

func (c *Unit) Start() {
	c.mtx.Lock()
	if c.started {
		c.mtx.Unlock()
		return
	}
	c.stoping = false
	go c.thWork()
	c.mtx.Unlock()
}

func (c *Unit) Stop() {
	c.mtx.Lock()
	if !c.started {
		c.mtx.Unlock()
		return
	}
	c.stoping = true
	c.mtx.Unlock()

	dtStartWaitingForStop := time.Now()

	for {
		if time.Since(dtStartWaitingForStop) > 1*time.Second {
			logger.Println("Unit stop timeout exceeded, force stopping")
			break
		}
		c.mtx.Lock()
		started := c.started
		c.mtx.Unlock()
		if !started {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *Unit) thWork() {
	c.started = true
	for !c.stoping {
		c.iUnit.Tick()
		time.Sleep(500 * time.Millisecond)
	}
	c.started = false
	c.stoping = false
}

func (c *Unit) Tick() {
}
