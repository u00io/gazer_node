package system

type UnitState struct {
	Id string

	UnitType            string
	UnitTypeDisplayName string

	Value string
}

type State struct {
	Units []UnitState
}
