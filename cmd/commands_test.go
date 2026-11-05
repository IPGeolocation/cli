package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const testKey = "TEST_KEY &+?=#%"

type fixtureTransport struct{}

func (fixtureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	key := r.URL.Query().Get("apiKey")
	status := map[string]int{"BAD_KEY": 401, "FORBIDDEN": 403, "RATE_LIMIT": 429, "SERVER_ERROR": 503}[key]
	body := "{"
	if status != 0 {
		body = fmt.Sprintf(`{"message":"Test failure for %s"}`, key)
	} else if key != "BAD_JSON" {
		status = 200
		var payload interface{}
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				return nil, err
			}
		}
		result := map[string]interface{}{
			"ip": "2001:db8::1", "location": map[string]string{"city": "Test City"},
			"_request": map[string]interface{}{"method": r.Method, "path": r.URL.Path, "query": r.URL.Query(), "payload": payload},
		}
		if r.URL.Query().Get("lang") == "ar" {
			result["location"] = map[string]string{"city": "مدينة"}
		}
		var response interface{} = result
		if strings.HasSuffix(r.URL.Path, "-bulk") {
			response = []interface{}{result}
		}
		encoded, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		body = string(encoded)
	} else {
		status = 200
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
}

// Each command runs in its own process so Cobra's global flags and saved keys cannot leak between cases.
func TestCLIProcess(t *testing.T) {
	if os.Getenv("CLI_TEST_PROCESS") != "1" {
		return
	}
	if os.Getenv("CLI_HTTP_TEST") == "1" {
		certificate, err := os.ReadFile(os.Getenv("SSL_CERT_FILE"))
		if err != nil {
			t.Fatal(err)
		}
		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM(certificate) {
			t.Fatal("invalid simulation CA")
		}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{RootCAs: roots}
		http.DefaultClient.Transport = transport
	} else {
		http.DefaultClient.Transport = fixtureTransport{}
	}
	for i, arg := range os.Args {
		if arg == "--" {
			rootCmd.SetArgs(os.Args[i+1:])
			Execute()
			os.Exit(0)
		}
	}
	t.Fatal("missing CLI arguments")
}

func newCLIProcess(home string, args ...string) *exec.Cmd {
	process := exec.Command(os.Args[0], append([]string{"-test.run=^TestCLIProcess$", "--"}, args...)...)
	if binary := os.Getenv("CLI_TEST_BINARY"); binary != "" {
		process = exec.Command(binary, args...)
	}
	process.Dir = home
	process.Env = append(os.Environ(), "CLI_TEST_PROCESS=1", "HOME="+home, "USERPROFILE="+home, "GORACE=atexit_sleep_ms=0")
	return process
}

func runCLI(t *testing.T, home string, args ...string) string {
	t.Helper()
	process := newCLIProcess(home, args...)
	output, err := process.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI %v failed: %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func TestCLIArabic(t *testing.T) {
	home := configuredHome(t, testKey)
	for _, args := range [][]string{
		{"ipgeo", "--ip", "8.8.8.8"},
		{"bulk-ip-geo", "--ips", "8.8.8.8,1.1.1.1"},
		{"timezone", "--ip", "8.8.8.8"},
		{"astronomy", "--ip", "8.8.8.8"},
		{"astronomy-timeseries", "--ip", "8.8.8.8", "--start-date", "2026-11-01", "--end-date", "2026-11-02"},
	} {
		t.Run(args[0], func(t *testing.T) {
			output := runCLI(t, home, append(args, "--lang", "ar", "--output", "raw")...)
			var response interface{}
			if err := json.Unmarshal([]byte(output), &response); err != nil {
				t.Fatalf("Arabic response invalid: %v", err)
			}
			if args[0] == "bulk-ip-geo" {
				response = response.([]interface{})[0]
			}
			result := response.(map[string]interface{})
			query := result["_request"].(map[string]interface{})["query"].(map[string]interface{})
			if !reflect.DeepEqual(query["lang"], []interface{}{"ar"}) || result["location"].(map[string]interface{})["city"] != "مدينة" {
				t.Fatal("Arabic request or response changed")
			}
		})
	}
}

func TestCLIRequestLimitFlags(t *testing.T) {
	home := configuredHome(t, testKey)
	if output := runCLI(t, home, "my-ip", "--timeout", "10m", "--max-response-mb", "512", "--output", "raw"); !strings.Contains(output, "2001:db8::1") {
		t.Fatalf("limit overrides failed: %s", output)
	}
	for _, args := range [][]string{
		{"ipgeo", "--timeout", "0s"},
		{"ipgeo", "--timeout=-1s"},
		{"ipgeo", "--timeout", "not-a-duration"},
		{"ipgeo", "--timeout", "999999999999999h"},
		{"ipgeo", "--max-response-mb", "0"},
		{"ipgeo", "--max-response-mb=-1"},
		{"ipgeo", "--max-response-mb", "1.5"},
		{"ipgeo", "--max-response-mb", "9223372036855"},
	} {
		process := newCLIProcess(home, args...)
		output, err := process.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "Error") || strings.Contains(string(output), `"_request"`) || strings.Contains(string(output), testKey) {
			t.Fatalf("invalid limits not rejected: args=%v err=%v output=%s", args, err, output)
		}
	}
	for _, command := range []string{"ipgeo", "bulk-ip-geo", "my-ip", "config"} {
		output := runCLI(t, home, command, "--help")
		if !strings.Contains(output, "--timeout") || !strings.Contains(output, "--max-response-mb") || !strings.Contains(output, "default 5m0s") || !strings.Contains(output, "default 256") {
			t.Fatalf("request controls absent from help for %s", command)
		}
	}
}

func configuredHome(t *testing.T, key string) string {
	t.Helper()
	home := t.TempDir()
	if output := runCLI(t, home, "config", "--apikey", key); output != "✅ API key saved securely." {
		t.Fatalf("config failed: %s", output)
	}
	return home
}

func TestCLICommands(t *testing.T) {
	home := configuredHome(t, testKey)
	cases := []struct {
		name, path string
		flags      []string
		query      url.Values
		payload    interface{}
	}{
		{"my-ip", "/getip", nil, url.Values{}, nil},
		{"getip", "/getip", nil, url.Values{}, nil},
		{"ipgeo", "/v3/ipgeo", []string{"--ip", "2001:db8::1", "--include", "security,time_zone", "--exclude", "currency", "--fields", "ip,location", "--lang", "en"}, url.Values{"ip": {"2001:db8::1"}, "include": {"security,time_zone"}, "excludes": {"currency"}, "fields": {"ip,location"}, "lang": {"en"}}, nil},
		{"ipgeo", "/v3/ipgeo", []string{"--ip", "example.test&injected=yes+#?%", "--excludes", "currency"}, url.Values{"ip": {"example.test&injected=yes+#?%"}, "excludes": {"currency"}}, nil},
		{"ip-security", "/v3/security", []string{"--ip", "8.8.8.8", "--exclude", "country", "--fields", "ip,security"}, url.Values{"ip": {"8.8.8.8"}, "excludes": {"country"}, "fields": {"ip,security"}}, nil},
		{"abuse", "/v3/abuse", []string{"--ip", "example.test", "--exclude", "country", "--fields", "ip,abuse"}, url.Values{"ip": {"example.test"}, "excludes": {"country"}, "fields": {"ip,abuse"}}, nil},
		{"asn", "/v3/asn", []string{"--ip", "8.8.8.8", "--asn", "AS15169", "--include", "peers,routes", "--exclude", "whois", "--fields", "asn"}, url.Values{"ip": {"8.8.8.8"}, "asn": {"AS15169"}, "include": {"peers,routes"}, "excludes": {"whois"}, "fields": {"asn"}}, nil},
		{"timezone", "/v3/timezone", []string{"--ip", "8.8.8.8", "--iata", "JFK", "--icao", "KJFK", "--lo", "USNYC", "--lang", "en", "--tz", "America/New_York", "--location", "New York & Queens", "--latitude", "40.7", "--longitude", "-74"}, url.Values{"ip": {"8.8.8.8"}, "iata_code": {"JFK"}, "icao_code": {"KJFK"}, "lo_code": {"USNYC"}, "lang": {"en"}, "tz": {"America/New_York"}, "location": {"New York & Queens"}, "lat": {"40.700000"}, "long": {"-74.000000"}}, nil},
		{"time-conversion", "/v3/timezone/convert", []string{"--tz_from", "America/New_York", "--tz_to", "Asia/Tokyo", "--location_from", "New York & Queens", "--location_to", "Tokyo", "--lat_from", "40.7", "--long_from", "-74", "--lat_to", "35.6", "--long_to", "139.7", "--iata_from", "JFK", "--iata_to", "HND", "--icao_from", "KJFK", "--icao_to", "RJTT", "--lo_from", "USNYC", "--lo_to", "JPTYO", "--time", "2026-10-29 15:00:00+05:00"}, url.Values{"tz_from": {"America/New_York"}, "tz_to": {"Asia/Tokyo"}, "location_from": {"New York & Queens"}, "location_to": {"Tokyo"}, "lat_from": {"40.700000"}, "long_from": {"-74.000000"}, "lat_to": {"35.600000"}, "long_to": {"139.700000"}, "iata_from": {"JFK"}, "iata_to": {"HND"}, "icao_from": {"KJFK"}, "icao_to": {"RJTT"}, "locode_from": {"USNYC"}, "locode_to": {"JPTYO"}, "time": {"2026-10-29 15:00:00+05:00"}}, nil},
		{"astronomy", "/v3/astronomy", []string{"--ip", "8.8.8.8", "--tz", "America/New_York", "--lang", "en", "--location", "New York & Queens", "--latitude", "40.7", "--longitude", "-74", "--elevation", "20.5"}, url.Values{"ip": {"8.8.8.8"}, "time_zone": {"America/New_York"}, "lang": {"en"}, "location": {"New York & Queens"}, "lat": {"40.700000"}, "long": {"-74.000000"}, "elevation": {"20.500000"}}, nil},
		{"astronomy-timeseries", "/v3/astronomy/timeSeries", []string{"--ip", "2001:db8::1", "--start-date", "2026-10-29", "--end-date", "2026-10-30", "--lang", "en", "--location", "New York & Queens", "--latitude", "40.7", "--longitude", "-74"}, url.Values{"ip": {"2001:db8::1"}, "dateStart": {"2026-10-29"}, "dateEnd": {"2026-10-30"}, "lang": {"en"}, "location": {"New York & Queens"}, "lat": {"40.700000"}, "long": {"-74.000000"}}, nil},
		{"bulk-ip-geo", "/v3/ipgeo-bulk", []string{"--ips", "8.8.8.8,2001:db8::1", "--include", "security,time_zone", "--exclude", "currency", "--fields", "ip,location", "--lang", "en"}, url.Values{"include": {"security,time_zone"}, "excludes": {"currency"}, "fields": {"ip,location"}, "lang": {"en"}}, map[string]interface{}{"ips": []interface{}{"8.8.8.8", "2001:db8::1"}}},
		{"bulk-ip-security", "/v3/security-bulk", []string{"--ips", "8.8.8.8,2001:db8::1", "--exclude", "country", "--fields", "ip,security"}, url.Values{"excludes": {"country"}, "fields": {"ip,security"}}, map[string]interface{}{"ips": []interface{}{"8.8.8.8", "2001:db8::1"}}},
		{"parse-user-agent", "/v3/user-agent", []string{"--user-agent", `Agent/1.0 "quoted" + &`}, url.Values{}, map[string]interface{}{"uaString": `Agent/1.0 "quoted" + &`}},
		{"parse-bulk-user-agents", "/v3/user-agent-bulk", []string{"--user-agents", "Agent/1.0,curl/8.0"}, url.Values{}, map[string]interface{}{"uaStrings": []interface{}{"Agent/1.0", "curl/8.0"}}},
	}
	for _, tc := range cases {
		if tc.path != "/getip" {
			tc.query.Set("apiKey", testKey)
		}
		for _, format := range []string{"raw", "pretty", "yaml", "table", "csv"} {
			t.Run(tc.name+"/"+format, func(t *testing.T) {
				args := append([]string{tc.name}, tc.flags...)
				output := runCLI(t, home, append(args, "--output", format)...)
				if format == "table" {
					if !strings.Contains(output, "2001:db8::1") || !strings.Contains(output, "Test City") {
						t.Fatalf("table lost result fields: %s", output)
					}
					return
				}
				if format == "csv" {
					rows := parseCSV(t, output)
					if len(rows) != 1 {
						t.Fatalf("expected one CSV row, got %d:\n%s", len(rows), output)
					}
					// Nested fields keep their path, so location.city survives flattening.
					if rows[0]["ip"] != "2001:db8::1" || rows[0]["location.city"] != "Test City" {
						t.Fatalf("CSV lost result fields: %s", output)
					}
					if _, charted := rows[0]["_request.path"]; !charted {
						t.Fatalf("CSV lost the echoed request: %s", output)
					}
					return
				}
				var result interface{}
				var err error
				if format == "yaml" {
					err = yaml.Unmarshal([]byte(output), &result)
				} else {
					err = json.Unmarshal([]byte(output), &result)
				}
				if err != nil {
					t.Fatalf("invalid %s output: %v\n%s", format, err, output)
				}
				if strings.HasSuffix(tc.path, "-bulk") {
					result = result.([]interface{})[0]
				}
				object := result.(map[string]interface{})
				if object["ip"] != "2001:db8::1" {
					t.Fatalf("result lost IP: %v", object)
				}
				actual := object["_request"].(map[string]interface{})
				method := "GET"
				if tc.payload != nil {
					method = "POST"
				}
				expected := map[string]interface{}{"method": method, "path": tc.path, "query": tc.query, "payload": tc.payload}
				encoded, _ := json.Marshal(expected)
				var normalized map[string]interface{}
				json.Unmarshal(encoded, &normalized)
				if !reflect.DeepEqual(actual, normalized) {
					t.Fatalf("request changed:\ngot  %v\nwant %v", actual, normalized)
				}
			})
		}
	}
	if output := runCLI(t, home, "getip", "--output", "text"); output != "2001:db8::1" {
		t.Fatalf("text output changed: %s", output)
	}
}

func TestCLIHTTPErrors(t *testing.T) {
	for _, key := range []string{"BAD_KEY", "FORBIDDEN", "RATE_LIMIT", "SERVER_ERROR", "BAD_JSON"} {
		t.Run(key, func(t *testing.T) {
			home := configuredHome(t, key)
			for _, args := range [][]string{{"ipgeo"}, {"parse-user-agent", "--user-agent", "test"}} {
				output := runCLI(t, home, args...)
				if !strings.Contains(output, "Error") || strings.Contains(output, key) || strings.Contains(output, `"_request"`) {
					t.Fatalf("unsafe or missing API error: %s", output)
				}
			}
		})
	}
}

func TestCLIConfigAndBulkFiles(t *testing.T) {
	home := configuredHome(t, "short")
	if output := runCLI(t, home, "config"); strings.Contains(output, "short") || !strings.Contains(output, "********") {
		t.Fatalf("short key exposed: %s", output)
	}
	if output := runCLI(t, home, "--version"); output != "ipgeolocation version 2.1.2" {
		t.Fatalf("incorrect version: %s", output)
	}
	if output := runCLI(t, home, "--help"); !strings.Contains(output, "time-conversion") {
		t.Fatalf("help lost commands: %s", output)
	}
	if output := runCLI(t, t.TempDir(), "ipgeo"); !strings.Contains(output, "API key not found") {
		t.Fatalf("missing config not reported: %s", output)
	}
	if output := runCLI(t, home, "astronomy-timeseries"); !strings.Contains(output, "both start and end dates") {
		t.Fatalf("missing dates not reported: %s", output)
	}
	if err := os.WriteFile(filepath.Join(home, "ips.txt"), []byte(" 2001:db8::1 \n\n1.1.1.1\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, command := range []string{"bulk-ip-geo", "bulk-ip-security"} {
		output := runCLI(t, home, command, "--ips", "8.8.8.8", "--file", "ips.txt", "--output-file", command)
		if !strings.Contains(output, "Output saved to file:") {
			t.Fatalf("bulk file failed: %s", output)
		}
		data, err := os.ReadFile(filepath.Join(home, command+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var results []map[string]interface{}
		if err := json.Unmarshal(data, &results); err != nil {
			t.Fatal(err)
		}
		payload := results[0]["_request"].(map[string]interface{})["payload"].(map[string]interface{})
		if !reflect.DeepEqual(payload["ips"], []interface{}{"8.8.8.8", "2001:db8::1", "1.1.1.1"}) {
			t.Fatalf("bulk inputs changed: %v", payload)
		}
	}
}

func parseCSV(t *testing.T, output string) []map[string]string {
	t.Helper()
	records, err := csv.NewReader(strings.NewReader(output)).ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV: %v\n%s", err, output)
	}
	if len(records) == 0 {
		return nil
	}
	rows := make([]map[string]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := map[string]string{}
		for i, header := range records[0] {
			row[header] = record[i]
		}
		rows = append(rows, row)
	}
	return rows
}

// --output-file follows --output: csv saves .csv, everything else still saves
// the .json it always did.
func TestCLICSVOutputFile(t *testing.T) {
	home := configuredHome(t, testKey)
	for _, command := range []string{"bulk-ip-geo", "bulk-ip-security"} {
		t.Run(command, func(t *testing.T) {
			name := command + "-csv"
			output := runCLI(t, home, command, "--ips", "8.8.8.8,1.1.1.1", "--output", "csv", "--output-file", name)
			if !strings.Contains(output, "Output saved to file: "+name+".csv") {
				t.Fatalf("csv file not reported: %s", output)
			}
			if _, err := os.Stat(filepath.Join(home, name+".json")); !os.IsNotExist(err) {
				t.Fatalf("csv run also wrote a .json file: %v", err)
			}

			saved, err := os.ReadFile(filepath.Join(home, name+".csv"))
			if err != nil {
				t.Fatal(err)
			}
			rows := parseCSV(t, string(saved))
			if len(rows) != 1 || rows[0]["ip"] != "2001:db8::1" || rows[0]["location.city"] != "Test City" {
				t.Fatalf("saved CSV lost fields: %s", saved)
			}
			// Both requested addresses ride in one cell, joined by the separator.
			if rows[0]["_request.payload.ips"] != "8.8.8.8|1.1.1.1" {
				t.Fatalf("array cell changed: %q", rows[0]["_request.payload.ips"])
			}
			// The file and stdout must agree.
			if strings.TrimSpace(strings.TrimSuffix(output, "Output saved to file: "+name+".csv")) != strings.TrimSpace(string(saved)) {
				t.Fatalf("file and stdout disagree:\n%s\n---\n%s", output, saved)
			}

			// Without --output csv the JSON file is unchanged.
			plain := command + "-plain"
			if output := runCLI(t, home, command, "--ips", "8.8.8.8", "--output-file", plain); !strings.Contains(output, "Output saved to file: "+plain+".json") {
				t.Fatalf("json file regressed: %s", output)
			}
			data, err := os.ReadFile(filepath.Join(home, plain+".json"))
			if err != nil {
				t.Fatal(err)
			}
			var results []map[string]interface{}
			if err := json.Unmarshal(data, &results); err != nil {
				t.Fatalf("json file is no longer JSON: %v", err)
			}
			if _, err := os.Stat(filepath.Join(home, plain+".csv")); !os.IsNotExist(err) {
				t.Fatalf("plain run wrote a .csv file: %v", err)
			}
		})
	}
}

// A request that fails must leave no file behind and must not overwrite one,
// for CSV exactly as for JSON.
func TestCLICSVOutputFileNotWrittenOnFailure(t *testing.T) {
	for _, key := range []string{"BAD_KEY", "SERVER_ERROR", "BAD_JSON"} {
		t.Run(key, func(t *testing.T) {
			home := configuredHome(t, key)
			for _, command := range []string{"bulk-ip-geo", "bulk-ip-security"} {
				for _, exists := range []bool{false, true} {
					name := fmt.Sprintf("%s-%v", command, exists)
					file := filepath.Join(home, name+".csv")
					if exists {
						if err := os.WriteFile(file, []byte("keep existing results"), 0600); err != nil {
							t.Fatal(err)
						}
					}
					output := runCLI(t, home, command, "--ips", "8.8.8.8", "--output", "csv", "--output-file", name)
					if strings.Contains(output, "Output saved") || strings.Contains(output, key) {
						t.Fatalf("failed request reported a saved file: %s", output)
					}
					data, err := os.ReadFile(file)
					if exists {
						if err != nil || string(data) != "keep existing results" {
							t.Fatalf("failed request changed the output file: %v", err)
						}
					} else if !os.IsNotExist(err) {
						t.Fatalf("failed request created an output file: %v", err)
					}
				}
			}
		})
	}
}

// --output-file adds .csv or .json only when the given name lacks it, so a
// name the user already suffixed is used as it stands and never doubles up.
func TestCLIOutputFileExtension(t *testing.T) {
	for _, tc := range []struct {
		given, format, want string
	}{
		// A bare name gains the extension that matches --output.
		{"plain", "", "plain.json"},
		{"plain-csv", "csv", "plain-csv.csv"},

		// A name that already carries it is used as given.
		{"suffixed.json", "", "suffixed.json"},
		{"suffixed.csv", "csv", "suffixed.csv"},

		// Matching ignores case.
		{"upper.JSON", "", "upper.JSON"},
		{"upper.CSV", "csv", "upper.CSV"},

		// A non-matching extension is kept, not swapped.
		{"mixed.json", "csv", "mixed.json.csv"},
		{"mixed.csv", "", "mixed.csv.json"},

		// A dot that is not the extension is left alone.
		{"2026.report", "", "2026.report.json"},
	} {
		for _, command := range []string{"bulk-ip-geo", "bulk-ip-security"} {
			name := tc.given
			label := command + "/" + name + "/" + tc.format
			t.Run(label, func(t *testing.T) {
				directory := t.TempDir()
				// A fresh HOME per case so leftover files cannot mask a mistake.
				if output := runCLI(t, directory, "config", "--apikey", testKey); !strings.Contains(output, "saved") {
					t.Fatalf("config failed: %s", output)
				}
				args := []string{command, "--ips", "8.8.8.8", "--output-file", name}
				if tc.format != "" {
					args = append(args, "--output", tc.format)
				}
				output := runCLI(t, directory, args...)
				if !strings.Contains(output, "Output saved to file: "+tc.want) {
					t.Fatalf("expected %q, got: %s", tc.want, output)
				}

				// Exactly one result file exists, and it is the expected one.
				entries, err := os.ReadDir(directory)
				if err != nil {
					t.Fatal(err)
				}
				var produced []string
				for _, entry := range entries {
					if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
						produced = append(produced, entry.Name())
					}
				}
				if len(produced) != 1 || produced[0] != tc.want {
					t.Fatalf("expected only %q, found %v", tc.want, produced)
				}

				// And it holds the format that was asked for.
				data, err := os.ReadFile(filepath.Join(directory, tc.want))
				if err != nil {
					t.Fatal(err)
				}
				if tc.format == "csv" {
					if rows := parseCSV(t, string(data)); len(rows) != 1 || rows[0]["ip"] != "2001:db8::1" {
						t.Fatalf("saved file is not the expected CSV: %s", data)
					}
				} else {
					var results []map[string]interface{}
					if err := json.Unmarshal(data, &results); err != nil {
						t.Fatalf("saved file is not JSON: %v\n%s", err, data)
					}
				}
			})
		}
	}
}

func TestCLICSVHelpAndFallback(t *testing.T) {
	home := configuredHome(t, testKey)
	for _, command := range []string{"ipgeo", "bulk-ip-geo", "my-ip", "timezone", "asn", "abuse", "parse-user-agent"} {
		if output := runCLI(t, home, command, "--help"); !strings.Contains(output, "csv") {
			t.Fatalf("csv missing from %s help: %s", command, output)
		}
	}
	// An unrecognised format still falls back to pretty JSON, as it always has.
	output := runCLI(t, home, "ipgeo", "--output", "csvv")
	if !json.Valid([]byte(output)) {
		t.Fatalf("unknown format no longer falls back to pretty: %s", output)
	}
}
