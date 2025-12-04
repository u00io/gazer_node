package system

import (
	"sync"
	"time"

	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/unit/unit000base"
	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/gomisc/logger"
)

type System struct {
	mtx    sync.Mutex
	client *U00

	units []unit000base.IUnit

	events []string
}

var Instance *System

func NewSystem() *System {
	var c System
	c.client = NewU00()
	c.units = make([]unit000base.IUnit, 0)
	return &c
}

func (c *System) Start() {
	err := config.Load()
	if err != nil {
		logger.Println("System::Start", "Cannot load config:", err)
	}

	for _, unitConfig := range config.Units() {
		unit := createUnitByType(unitConfig.Type)
		if unit == nil {
			logger.Println("System::Start", "Cannot create unit of type", unitConfig.Type)
			continue
		}
		unit.SetId(unitConfig.Id)
		unit.SetConfig(*unitConfig)
		key := utils.NewKeyFromPrivate(unitConfig.PrivateKey)
		unit.SetKey(key)
		c.units = append(c.units, unit)
	}

	c.client.Run()
	c.startAllUnits()
	go c.thWork()
}

func (c *System) Stop() {
}

func (c *System) StartUnit(unitId string) {
	for _, unit := range c.units {
		if unit.GetId() == unitId {
			unitConfig := config.UnitById(unitId)
			if unitConfig == nil {
				return
			}
			unit.SetConfig(*unitConfig)
			unit.Start()
			return
		}
	}
}

func (c *System) StopUnit(unitId string) {
	for _, unit := range c.units {
		if unit.GetId() == unitId {
			unit.Stop()
			return
		}
	}
}

func (c *System) SetUnitTranslate(unitId string, translate bool) {
	configUnit := config.UnitById(unitId)
	if configUnit == nil {
		return
	}
	configUnit.Translate = translate
	config.Save()

	Instance.EmitEvent("config_changed")
	Instance.StopUnit(unitId)
	Instance.StartUnit(unitId)
}

func (c *System) EmitEvent(event string) {
	c.mtx.Lock()
	c.events = append(c.events, event)
	c.mtx.Unlock()
}

func (c *System) GetAndClearEvents() []string {
	c.mtx.Lock()
	events := c.events
	c.events = make([]string, 0)
	c.mtx.Unlock()
	return events
}

func (c *System) startAllUnits() {
	for _, unit := range c.units {
		unit.Start()
	}
}

func (c *System) Test() {
	// c.client.WriteValue("Test Value")
}

func (c *System) thWork() {
	for {
		c.SendValues()
		time.Sleep(200 * time.Millisecond)
	}
}

func (c *System) SendValues() {
	for _, unit := range c.units {
		if unit.GetConfig().Translate == false {
			continue
		}
		values := unit.GetValues()
		//name := unit.GetParameterString("001_name_str", unit.GetType())

		for _, value := range values {
			var items []ItemToSet
			items = append(items, ItemToSet{
				Path:  value.Key,
				Name:  value.Name,
				Value: value.Value,
				Uom:   value.Uom,
			})
			key := unit.GetKey()
			if key != nil {
				c.client.Write(unit.GetKey().PrivateKey, items)
			}
		}
	}
}

func (c *System) GetState() State {
	var state State
	for _, unit := range c.units {
		typeDisplayName := ""
		if record, exists := Registry.UnitTypes[unit.GetType()]; exists {
			typeDisplayName = record.TypeDisplayName
		}

		values := make([]UnitStateDataItem, 0)
		for _, v := range unit.GetValues() {
			values = append(values, UnitStateDataItem{
				Key:   v.Key,
				Name:  v.Name,
				Value: v.Value,
				UOM:   v.Uom,
			})
		}

		unitState := UnitState{
			Id:                  unit.GetId(),
			UnitType:            unit.GetType(),
			UnitTypeDisplayName: typeDisplayName,
			Values:              values,
			Config:              unit.GetConfig(),
		}
		state.Units = append(state.Units, unitState)
	}
	return state
}

func (c *System) AddUnit(unitConfig *config.ConfigUnit) {
	id := config.AddUnit(unitConfig)

	unit := createUnitByType(unitConfig.Type)
	if unit == nil {
		logger.Println("System::AddUnit", "Cannot create unit of type", unitConfig.Type)
		return
	}
	unit.SetId(id)
	unit.SetKey(utils.NewKeyFromPrivate(unitConfig.PrivateKey))
	unit.SetConfig(*unitConfig)
	c.units = append(c.units, unit)
	unit.Start()
}

func (c *System) RemoveUnit(unitId string) {
	config.RemoveUnit(unitId)
	for i, unit := range c.units {
		if unit.GetId() == unitId {
			unit.Stop()
			c.units = append(c.units[:i], c.units[i+1:]...)
			c.EmitEvent("config_changed")
			return
		}
	}
}

func (c *System) GetUnitDefaultItemValue(unitId string) string {
	for _, unit := range c.units {
		if unit.GetId() == unitId {
			v := unit.GetValue("/")
			return v.Value
		}
	}
	return ""
}
