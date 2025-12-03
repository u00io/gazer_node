package system

import (
	"sort"

	unit00base "github.com/u00io/gazer_node/unit/unit_00_base"
	unit01filecontent "github.com/u00io/gazer_node/unit/unit_01_file_content"
	unit02currenttime "github.com/u00io/gazer_node/unit/unit_02_current_time"
)

type UnitCategory struct {
	Name        string
	Description string
}

type UnitTypeRecord struct {
	TypeName        string
	TypeDisplayName string

	Categories  []string
	Constructor func() unit00base.IUnit
}

type UnitsRegistry struct {
	UnitCategories []*UnitCategory
	UnitTypes      map[string]*UnitTypeRecord
}

func (r *UnitsRegistry) RegisterUnitType(unitType string, displayName string, constructor func() unit00base.IUnit, categories ...string) {
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
	Registry.RegisterUnitType("unit01filecontent", "File Content", unit01filecontent.New, "File Operations", "Content Management")
	Registry.RegisterUnitType("unit02currenttime", "Current Time", unit02currenttime.New, "Time Management", "Utilities")
	Registry.UpdateUnitCategories()
}

func createUnitByType(unitType string) unit00base.IUnit {
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
