package cmd

import (
	"crypto/tls"
	"crypto/x509"
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

func runCLI(t *testing.T, home string, args ...string) string {
	t.Helper()
	process := exec.Command(os.Args[0], append([]string{"-test.run=^TestCLIProcess$", "--"}, args...)...)
	if binary := os.Getenv("CLI_TEST_BINARY"); binary != "" {
		process = exec.Command(binary, args...)
	}
	process.Dir = home
	process.Env = append(os.Environ(), "CLI_TEST_PROCESS=1", "HOME="+home, "USERPROFILE="+home, "GORACE=atexit_sleep_ms=0")
	output, err := process.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI %v failed: %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
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
		for _, format := range []string{"raw", "pretty", "yaml", "table"} {
			t.Run(tc.name+"/"+format, func(t *testing.T) {
				args := append([]string{tc.name}, tc.flags...)
				output := runCLI(t, home, append(args, "--output", format)...)
				if format == "table" {
					if !strings.Contains(output, "2001:db8::1") || !strings.Contains(output, "Test City") {
						t.Fatalf("table lost result fields: %s", output)
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
	if output := runCLI(t, home, "--version"); output != "ipgeolocation version 2.1.1" {
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
