package system

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

	Values []UnitStateDataItem
}

type State struct {
	Units []UnitState
}
