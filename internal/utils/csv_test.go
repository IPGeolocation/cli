package utils

import (
	"encoding/csv"
	"strings"
	"testing"
)

func csvOf(t *testing.T, body string) string {
	t.Helper()
	output, err := ConvertToCSV([]byte(body))
	if err != nil {
		t.Fatalf("ConvertToCSV(%s) failed: %v", body, err)
	}
	return output
}

func TestConvertToCSVShapes(t *testing.T) {
	for _, tc := range []struct {
		name, body, want string
	}{
		{
			name: "nested objects become dotted columns",
			body: `{"ip":"8.8.8.8","location":{"city":"Mountain View","is_eu":false}}`,
			want: "ip,location.city,location.is_eu\n8.8.8.8,Mountain View,false\n",
		},
		{
			// The bulk endpoints answer with an array; each element is a record.
			name: "top-level array becomes one row per element",
			body: `[{"ip":"8.8.8.8"},{"ip":"1.1.1.1"}]`,
			want: "ip\n8.8.8.8\n1.1.1.1\n",
		},
		{
			// Records need not agree on their keys; every column still appears.
			name: "columns are the union across rows",
			body: `[{"ip":"8.8.8.8","city":"MV"},{"ip":"1.1.1.1","asn":"AS13335"}]`,
			want: "asn,city,ip\n,MV,8.8.8.8\nAS13335,,1.1.1.1\n",
		},
		{
			name: "array of scalars joins into one cell",
			body: `{"languages":["en-US","es-US","haw"]}`,
			want: "languages\nen-US|es-US|haw\n",
		},
		{
			// Keeping a column per field beats hiding objects in a joined cell.
			name: "array of objects keeps a column per element",
			body: `{"routes":[{"prefix":"8.8.8.0/24"},{"prefix":"8.8.4.0/24"}]}`,
			want: "routes.0.prefix,routes.1.prefix\n8.8.8.0/24,8.8.4.0/24\n",
		},
		{
			name: "array mixing scalars and objects is indexed",
			body: `{"items":["a",{"b":1}]}`,
			want: "items.0,items.1.b\na,1\n",
		},
		{
			name: "empty containers leave an empty cell",
			body: `{"ip":"8.8.8.8","tags":[],"meta":{}}`,
			want: "ip,meta,tags\n8.8.8.8,,\n",
		},
		{
			name: "null becomes an empty cell",
			body: `{"ip":"8.8.8.8","city":null}`,
			want: "city,ip\n,8.8.8.8\n",
		},
		{
			name: "booleans keep their JSON spelling",
			body: `{"is_eu":true,"is_dst":false}`,
			want: "is_dst,is_eu\nfalse,true\n",
		},
		{
			name: "deep nesting keeps the whole path",
			body: `{"time_zone":{"dst_end":{"gap":false,"utc_time":"2026-11-01 09:00"}}}`,
			want: "time_zone.dst_end.gap,time_zone.dst_end.utc_time\nfalse,2026-11-01 09:00\n",
		},
		{
			name: "a bare scalar response gets a value column",
			body: `"8.8.8.8"`,
			want: "value\n8.8.8.8\n",
		},
		{
			name: "an empty object yields a single empty row",
			body: `{}`,
			want: "value\n\n",
		},
		{
			name: "unicode survives unchanged",
			body: `{"city":"مدينة","flag":"🇺🇸"}`,
			want: "city,flag\nمدينة,🇺🇸\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := csvOf(t, tc.body); got != tc.want {
				t.Fatalf("CSV changed:\ngot  %q\nwant %q", got, tc.want)
			}
		})
	}
}

// An empty bulk result has no columns to describe, so it writes nothing rather
// than a stray header line.
func TestConvertToCSVEmptyArray(t *testing.T) {
	if got := csvOf(t, `[]`); got != "" {
		t.Fatalf("empty array produced %q", got)
	}
}

// float64 would render a Unix timestamp as 1.772535675329e+09 and round long
// identifiers, so numbers must pass through exactly as the API sent them.
func TestConvertToCSVKeepsNumbersExact(t *testing.T) {
	body := `{"unix":1772535675.329,"big":12345678901234567890,"neg":-122.08421,"zero":0,"exp":1e3,"offset":-8}`
	want := "big,exp,neg,offset,unix,zero\n12345678901234567890,1e3,-122.08421,-8,1772535675.329,0\n"
	if got := csvOf(t, body); got != want {
		t.Fatalf("number formatting changed:\ngot  %q\nwant %q", got, want)
	}
}

// Separators, quotes and newlines inside values must survive a round trip.
func TestConvertToCSVQuotesSpecialCharacters(t *testing.T) {
	body := `{"org":"Google, LLC","quote":"say \"hi\"","line":"a\nb","tab":"a\tb"}`
	output := csvOf(t, body)

	rows, err := csv.NewReader(strings.NewReader(output)).ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v\n%s", err, output)
	}
	if len(rows) != 2 {
		t.Fatalf("expected a header and one row, got %d rows", len(rows))
	}
	values := map[string]string{}
	for i, header := range rows[0] {
		values[header] = rows[1][i]
	}
	for header, want := range map[string]string{
		"org": "Google, LLC", "quote": `say "hi"`, "line": "a\nb", "tab": "a\tb",
	} {
		if values[header] != want {
			t.Fatalf("%s round-tripped as %q, want %q", header, values[header], want)
		}
	}
}

// Response values come from third parties, so a cell must not open as a formula
// in a spreadsheet. Numbers stay untouched so coordinates and offsets survive.
func TestConvertToCSVEscapesFormulas(t *testing.T) {
	for value, want := range map[string]string{
		// Formulas are quoted off.
		`=1+1`:               `'=1+1`,
		`=cmd|' /c calc'!A0`: `'=cmd|' /c calc'!A0`,
		`+SUM(A1)`:           `'+SUM(A1)`,
		`-cmd|' /c calc'!A0`: `'-cmd|' /c calc'!A0`,
		`@import`:            `'@import`,
		"\tlead":             "'\tlead",
		"\rlead":             "'\rlead",
		// A lone sign is real data: the astronomy endpoint reports an absent
		// sun_status or moon_status as "-". It cannot be a formula.
		"+": "+",
		"-": "-",
		// Signed data is left exactly as the API sent it. The DST duration
		// field carries both signs, so the two must be treated alike.
		"-1.00H":     "-1.00H",
		"+1.00H":     "+1.00H",
		"-122.08421": "-122.08421",
		"+.5":        "+.5",
		"-8":         "-8",
		"+1":         "+1",
		// Ordinary text is untouched.
		"Google LLC":    "Google LLC",
		"":              "",
		"not=a=formula": "not=a=formula",
		"مدينة":         "مدينة",
	} {
		if got := escapeFormula(value); got != want {
			t.Fatalf("escapeFormula(%q) = %q, want %q", value, got, want)
		}
	}

	// And it reaches the rendered output.
	if got := csvOf(t, `{"org":"=1+1","lat":"-37.4"}`); got != "lat,org\n-37.4,'=1+1\n" {
		t.Fatalf("escaping missing from output: %q", got)
	}
}

func TestConvertToCSVRejectsInvalidJSON(t *testing.T) {
	for _, body := range []string{``, `{`, `not json`, `{"a":}`} {
		if _, err := ConvertToCSV([]byte(body)); err == nil {
			t.Fatalf("invalid JSON %q was accepted", body)
		}
	}
}

// Go randomises map iteration, so column order has to be imposed, not observed.
func TestConvertToCSVIsDeterministic(t *testing.T) {
	body := `{"z":1,"a":2,"m":{"y":3,"b":4},"list":["q","p"]}`
	first := csvOf(t, body)
	for i := 0; i < 50; i++ {
		if got := csvOf(t, body); got != first {
			t.Fatalf("output varies between runs:\n%q\n%q", first, got)
		}
	}
	if !strings.HasPrefix(first, "a,list,m.b,m.y,z\n") {
		t.Fatalf("columns are not sorted: %q", first)
	}
}
