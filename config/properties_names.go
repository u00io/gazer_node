package config

var propertyNames = map[string]string{
	"0000_00_name_str":         "Unit Name",
	"0100_00_offset_num":       "Offset",
	"0102_00_process_name_str": "Process Name",
	"0102_01_process_id_int":   "Process ID",
}

func PropName(propCode string) string {
	result := propCode
	r, ok := propertyNames[propCode]
	if ok {
		result = r
	}
	return result
}
