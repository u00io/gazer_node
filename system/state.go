package system

import "github.com/u00io/gazer_node/config"

type UnitStateDataItem struct {
	Key   string
	Name  string
	Value string
	UOM   string
}

type UnitState struct {
	Id string

	UnitType            string
	UnitTypeDisplayName string

	Config config.ConfigUnit

	Values []UnitStateDataItem
}

type State struct {
	Units []UnitState
}
