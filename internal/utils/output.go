package utils

import (
	"fmt"
	"path/filepath"
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

// OutputFilePath works out the file an --output-file value should write to.
// The extension is appended only when the name does not already end in it, so
// "results" and "results.csv" both save to results.csv and neither doubles up.
//
// An extension that does not match is left in place rather than swapped, so
// --output-file results.json with --output csv writes results.json.csv. That
// keeps the name the caller gave and never silently retargets a different file.
// The comparison ignores case, so results.CSV is accepted as it stands.
func OutputFilePath(name, extension string) string {
	if strings.EqualFold(filepath.Ext(name), extension) {
		return name
	}
	return name + extension
}
