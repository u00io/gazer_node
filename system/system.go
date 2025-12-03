package system

import (
	"time"

	"github.com/u00io/gazer_node/config"
	unit00base "github.com/u00io/gazer_node/unit/unit_00_base"
	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/gomisc/logger"
)

type System struct {
	client *U00

	units []unit00base.IUnit
}

var Instance *System

func NewSystem() *System {
	var c System
	c.client = NewU00()
	c.units = make([]unit00base.IUnit, 0)
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

func (c *System) LoadDefaultConfig() {
	u01 := createUnitByType("unit01filecontent")
	c.units = append(c.units, u01)
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
			c.client.Write(unit.GetKey().PrivateKey, items)
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
	unit := createUnitByType(unitType)
	if unit == nil {
		logger.Println("System::AddUnit", "Cannot create unit of type", unitType)
		return
	}
	unit.SetId(id)
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
