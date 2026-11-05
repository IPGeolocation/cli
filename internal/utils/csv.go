package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ArraySeparator joins the elements of an array of scalars inside a single CSV cell.
const ArraySeparator = "|"

// ConvertToCSV renders a JSON API response as CSV text.
//
// A top-level JSON array becomes one row per element, which is what the bulk
// endpoints return; any other response becomes a single row. Nested objects are
// flattened into dotted column names ("location.city"). An array whose elements
// are all scalars is joined into one cell with ArraySeparator; an array holding
// objects keeps a column per element instead ("routes.0.prefix"), so its fields
// stay addressable. Columns are the union of the keys of every row, sorted, so
// the same response always produces the same file.
func ConvertToCSV(body []byte) (string, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	// Numbers stay exactly as the API sent them; float64 would turn a Unix
	// timestamp into 1.772535675329e+09 and lose digits off long IDs.
	decoder.UseNumber()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	records, isList := data.([]interface{})
	if !isList {
		records = []interface{}{data}
	}

	rows := make([]map[string]string, 0, len(records))
	columns := map[string]bool{}
	for _, record := range records {
		row := flattenRecord(record)
		for column := range row {
			columns[column] = true
		}
		rows = append(rows, row)
	}

	if len(columns) == 0 {
		return "", nil
	}

	headers := make([]string, 0, len(columns))
	for column := range columns {
		headers = append(headers, column)
	}
	sort.Strings(headers)

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	if err := writer.Write(headers); err != nil {
		return "", err
	}
	for _, row := range rows {
		line := make([]string, len(headers))
		for i, header := range headers {
			line[i] = escapeFormula(row[header])
		}
		if err := writer.Write(line); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// flattenRecord turns one record into its CSV cells, keyed by column name.
func flattenRecord(record interface{}) map[string]string {
	row := map[string]string{}
	flattenValue("", record, row)
	// A record that is a bare scalar, an empty object or an empty array leaves
	// its cell unnamed; give it a column rather than writing a blank header.
	if value, unnamed := row[""]; unnamed {
		delete(row, "")
		row["value"] = value
	}
	return row
}

func flattenValue(prefix string, value interface{}, row map[string]string) {
	switch typed := value.(type) {
	case map[string]interface{}:
		if len(typed) == 0 {
			row[prefix] = ""
			return
		}
		for key, nested := range typed {
			flattenValue(joinColumn(prefix, key), nested, row)
		}
	case []interface{}:
		if len(typed) == 0 {
			row[prefix] = ""
			return
		}
		if scalarsOnly(typed) {
			parts := make([]string, len(typed))
			for i, item := range typed {
				parts[i] = scalarToString(item)
			}
			row[prefix] = strings.Join(parts, ArraySeparator)
			return
		}
		for i, item := range typed {
			flattenValue(joinColumn(prefix, strconv.Itoa(i)), item, row)
		}
	default:
		row[prefix] = scalarToString(typed)
	}
}

func joinColumn(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func scalarsOnly(values []interface{}) bool {
	for _, value := range values {
		switch value.(type) {
		case map[string]interface{}, []interface{}:
			return false
		}
	}
	return true
}

func scalarToString(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case bool:
		return strconv.FormatBool(typed)
	case json.Number:
		return typed.String()
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", typed)
	}
}

// escapeFormula stops a spreadsheet from running a cell as a formula. Response
// values such as organization names and whois text come from third parties, so
// a leading =, @ or control character is quoted off.
//
// A leading + or - is only quoted when a non-digit follows it. Signed numbers
// and signed data are ordinary values here, such as a "-122.08421" longitude,
// a "+1" calling code, the "+1.00H" and "-1.00H" DST durations and the lone
// "-" the astronomy endpoint uses for an absent sun or moon status; quoting
// those would corrupt real columns. A sign followed by anything else is the
// shape a formula takes, as in "+SUM(A1)" or "-cmd|' /c calc'!A0".
func escapeFormula(value string) string {
	if value == "" {
		return value
	}
	switch value[0] {
	case '=', '@', '\t', '\r':
		return "'" + value
	case '+', '-':
		// A lone sign is a value in its own right: the astronomy endpoint
		// reports sun_status and moon_status as "-". It cannot be a formula.
		if len(value) == 1 || value[1] == '.' || (value[1] >= '0' && value[1] <= '9') {
			return value
		}
		return "'" + value
	}
	return value
}
