package system

import (
	"sort"

	"github.com/u00io/gazer_node/unit/unit000base"
	"github.com/u00io/gazer_node/unit/unit001demosignal"
)

type UnitCategory struct {
	Name        string
	Description string
}

type UnitTypeRecord struct {
	TypeName        string
	TypeDisplayName string

	Categories  []string
	Constructor func() unit000base.IUnit
}

type UnitsRegistry struct {
	UnitCategories []*UnitCategory
	UnitTypes      map[string]*UnitTypeRecord
}

func (r *UnitsRegistry) RegisterUnitType(unitType string, displayName string, constructor func() unit000base.IUnit, categories ...string) {
	var record UnitTypeRecord
	record.TypeName = unitType
	record.TypeDisplayName = displayName
	record.Constructor = constructor
	record.Categories = categories
	r.UnitTypes[unitType] = &record
}

var Registry UnitsRegistry

func init() {
	Registry.UnitTypes = make(map[string]*UnitTypeRecord)

	// General units
	Registry.RegisterUnitType("unit001demosignal", "Demo Signal", unit001demosignal.New, "General")

	// Computer units
	/*Registry.RegisterUnitType("unit101networkadapters", "Network Adapters", unit02currenttime.New, "Computer")
	Registry.RegisterUnitType("unit102process", "Process", unit02currenttime.New, "Computer")
	Registry.RegisterUnitType("unit103storage", "Storage", unit02currenttime.New, "Computer")
	Registry.RegisterUnitType("unit104memory", "Memory", unit02currenttime.New, "Computer")

	// File units
	Registry.RegisterUnitType("unit201filesize", "File Size", unit02currenttime.New, "Files")
	Registry.RegisterUnitType("unit202filecontent", "File Content", unit02currenttime.New, "Files")
	Registry.RegisterUnitType("unit203filetail", "File Tail", unit02currenttime.New, "Files")

	// Serial Port units
	Registry.RegisterUnitType("unit301serialportkeyvalue", "Serial Port Key=Value", unit02currenttime.New, "SerialPort")
	Registry.RegisterUnitType("unit302serialportline", "Serial Port Lines", unit02currenttime.New, "SerialPort")

	// HTTP units
	Registry.RegisterUnitType("unit401httpvalue", "Http Get Value", unit02currenttime.New, "HTTP/Web")
	Registry.RegisterUnitType("unit402httpjson", "Http Get Json", unit02currenttime.New, "HTTP/Web")

	// Network units
	Registry.RegisterUnitType("unit501networkping", "Ping", unit02currenttime.New, "Network")
	Registry.RegisterUnitType("unit502networktcpconnect", "Tcp Connect", unit02currenttime.New, "Network")

	// Raspberry Pi units
	Registry.RegisterUnitType("unit601raspberrypiTemp", "Temperature", unit02currenttime.New, "Raspberry Pi")
	Registry.RegisterUnitType("unit602raspberrypiGpioPins", "GPIO Pins", unit02currenttime.New, "Raspberry Pi")

	// Console units
	Registry.RegisterUnitType("unit701plainvalue", "Plain Value", unit02currenttime.New, "Console")
	Registry.RegisterUnitType("unit702keyvalue", "Key=Value", unit02currenttime.New, "Console")*/

	Registry.UpdateUnitCategories()
}

func createUnitByType(unitType string) unit000base.IUnit {
	if record, exists := Registry.UnitTypes[unitType]; exists {
		return record.Constructor()
	}
	return nil
}

func (r *UnitsRegistry) UpdateUnitCategories() {
	categoriesMap := make(map[string]*UnitCategory)
	for _, record := range r.UnitTypes {
		for _, category := range record.Categories {
			if _, exists := categoriesMap[category]; !exists {
				categoriesMap[category] = &UnitCategory{Name: category}
			}
		}
	}
	r.UnitCategories = make([]*UnitCategory, 0, len(categoriesMap))
	for _, category := range categoriesMap {
		r.UnitCategories = append(r.UnitCategories, category)
	}

	sort.Slice(r.UnitCategories, func(i, j int) bool {
		return r.UnitCategories[i].Name < r.UnitCategories[j].Name
	})
}

func (r *UnitsRegistry) GetUnitTypeDefaultParameters(unitType string) map[string]string {
	result := make(map[string]string)
	if record, exists := r.UnitTypes[unitType]; exists {
		unit := record.Constructor()
		result = unit.GetConfig()
	}
	result["000_name_str"] = unitType + " Instance"
	return result
}
