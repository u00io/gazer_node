package unit00base

import (
	"sync"
	"time"

	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/gomisc/logger"
)

type Unit struct {
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
	SetConfig(config map[string]string)
	GetId() string
	GetKey() *utils.Key
	GetType() string
	GetValue(key string) string
	SetValue(key, value string)
	Stop()
	Tick()
}

func (c *Unit) Init(iUnit IUnit) {
	c.config = make(map[string]string)
	c.values = make(map[string]string)
	c.iUnit = iUnit
	c.key = utils.NewKey()
}

func (c *Unit) GetId() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.key.String()
}

func (c *Unit) GetKey() *utils.Key {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.key
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
