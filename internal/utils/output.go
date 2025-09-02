package utils

import (
	"fmt"
	"strings"
)

func PrintAsTable(data interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch val := data.(type) {
	case map[string]interface{}:
		for key, value := range val {
			switch value.(type) {
			case map[string]interface{}, []interface{}:
				fmt.Printf("%s%s:\n", indentStr, ToTitle(key))
				PrintAsTable(value, indent+1)
			default:
				fmt.Printf("%s%-20s: %v\n", indentStr, ToTitle(key), value)
			}
		}
	case []interface{}:
		for i, item := range val {
			fmt.Printf("%s[%d]:\n", indentStr, i)
			PrintAsTable(item, indent+1)
		}
	default:
		fmt.Printf("%s%v\n", indentStr, val)
	}
}

func ToTitle(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

func ConvertToXML(data interface{}) string {
	return ""
}
