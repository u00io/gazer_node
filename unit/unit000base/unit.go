package unit000base

import (
	"sync"
	"time"

	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/gomisc/logger"
)

type ItemValue struct {
	Key   string
	Name  string
	Value string
	Uom   string
}

type Unit struct {
	id      string
	mtx     sync.Mutex
	key     *utils.Key
	tp      string
	started bool
	stoping bool
	config  *config.ConfigUnit
	values  map[string]ItemValue
	iUnit   IUnit
}

type IUnit interface {
	Start()
	Stop()

	SetId(id string)
	GetId() string

	GetKey() *utils.Key
	SetKey(key *utils.Key)

	SetConfig(config config.ConfigUnit)
	GetConfig() config.ConfigUnit

	Config() *config.ConfigUnit

	GetType() string
	GetValue(key string) ItemValue
	GetValues() map[string]ItemValue
	SetValue(key string, name string, value string, uom string)
	Tick()
}

func (c *Unit) Init(iUnit IUnit) {
	c.values = make(map[string]ItemValue)
	c.config = config.NewConfigUnit()
	c.config.Type = c.GetType()
	c.config.Parameters = make(map[string]string)
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

func (c *Unit) SetConfig(config config.ConfigUnit) {
	c.mtx.Lock()
	c.config = &config
	c.mtx.Unlock()
}

func (c *Unit) GetConfig() config.ConfigUnit {
	c.mtx.Lock()
	result := *c.config
	defer c.mtx.Unlock()
	return result
}

func (c *Unit) Config() *config.ConfigUnit {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.config
}

func (c *Unit) GetValue(key string) ItemValue {
	c.mtx.Lock()
	value, exists := c.values[key]
	c.mtx.Unlock()
	if !exists {
		return ItemValue{}
	}
	return value
}

func (c *Unit) GetValues() map[string]ItemValue {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	result := make(map[string]ItemValue)
	for k, v := range c.values {
		result[k] = v
	}
	return result
}

func (c *Unit) SetValue(key string, name string, value string, uom string) {
	itemValue := ItemValue{
		Key:   key,
		Name:  name,
		Value: value,
		Uom:   uom,
	}
	c.mtx.Lock()
	c.values[key] = itemValue
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
