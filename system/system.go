package system

import (
	"time"

	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/unit/unit000base"
	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/gomisc/logger"
)

type System struct {
	client *U00

	units []unit000base.IUnit
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
		unit.SetConfig(unitConfig.Parameters)
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
		value := unit.GetValue("value")
		if value != "" {
			var items []ItemToSet
			items = append(items, ItemToSet{
				Path:  "/",
				Name:  "value",
				Value: value,
				Uom:   "Value",
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

		unitState := UnitState{
			Id:                  unit.GetId(),
			UnitType:            unit.GetType(),
			UnitTypeDisplayName: typeDisplayName,
			Value:               unit.GetValue("value"),
		}
		state.Units = append(state.Units, unitState)
	}
	return state
}

func (c *System) AddUnit(unitType string, parameters map[string]string) {
	id := config.AddUnit(unitType, parameters)

	unitConfig := config.UnitById(id)
	if unitConfig == nil {
		logger.Println("System::AddUnit", "Cannot find added unit in config with id", id)
		return
	}

	unit := createUnitByType(unitType)
	if unit == nil {
		logger.Println("System::AddUnit", "Cannot create unit of type", unitType)
		return
	}
	unit.SetId(id)
	unit.SetKey(utils.NewKeyFromPrivate(unitConfig.PrivateKey))
	unit.SetConfig(parameters)
	c.units = append(c.units, unit)
	unit.Start()
}

func (c *System) RemoveUnit(unitId string) {
	config.RemoveUnit(unitId)
	for i, unit := range c.units {
		if unit.GetId() == unitId {
			unit.Stop()
			c.units = append(c.units[:i], c.units[i+1:]...)
			return
		}
	}
}

func (c *System) GetUnitDefaultItemValue(unitId string) string {
	for _, unit := range c.units {
		if unit.GetId() == unitId {
			return unit.GetValue("value")
		}
	}
	return ""
}
