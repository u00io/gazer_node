package unit000base

import (
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
