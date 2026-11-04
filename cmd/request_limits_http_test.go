package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIRequestLimitsHTTP(t *testing.T) {
	if os.Getenv("CLI_HTTP_TEST") != "1" {
		t.Skip("requires the local HTTP simulation")
	}
	for _, key := range []string{"SLOW_HEADERS", "SLOW_BODY"} {
		home := configuredHome(t, key)
		for _, args := range [][]string{{"ipgeo"}, {"parse-user-agent", "--user-agent", "test"}} {
			output := runCLI(t, home, append(args, "--timeout", "100ms")...)
			if !strings.Contains(output, "request timed out") || strings.Contains(output, key) || strings.Contains(output, `"_request"`) {
				t.Fatalf("timeout did not stop safely: %s", output)
			}
		}
		output := runCLI(t, home, "ipgeo", "--timeout", "2s", "--output", "raw")
		if !json.Valid([]byte(output)) || !strings.Contains(output, "2001:db8::1") {
			t.Fatalf("longer timeout did not complete: %s", output)
		}
	}
	for _, key := range []string{"LIMIT_RESPONSE", "SLOW_BODY"} {
		home := configuredHome(t, key)
		for _, command := range []string{"bulk-ip-geo", "bulk-ip-security"} {
			for _, exists := range []bool{false, true} {
				name := fmt.Sprintf("%s-%v", command, exists)
				file := filepath.Join(home, name+".json")
				if exists {
					if err := os.WriteFile(file, []byte("keep existing results"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				args := []string{command, "--ips", "8.8.8.8", "--output-file", name}
				want := "response exceeds 1 MB"
				if key == "LIMIT_RESPONSE" {
					args = append(args, "--max-response-mb", "1")
				} else {
					args = append(args, "--timeout", "100ms")
					want = "request timed out"
				}
				output := runCLI(t, home, args...)
				if !strings.Contains(output, want) || strings.Contains(output, key) || strings.Contains(output, `"_request"`) || strings.Contains(output, "Output saved") {
					t.Fatalf("failed bulk request produced output: %s", output)
				}
				data, err := os.ReadFile(file)
				if exists {
					if err != nil || string(data) != "keep existing results" {
						t.Fatalf("failed request changed output file: %v", err)
					}
				} else if !os.IsNotExist(err) {
					t.Fatal("failed request created an output file")
				}
			}
		}
		if key == "LIMIT_RESPONSE" {
			output := runCLI(t, home, "bulk-ip-geo", "--ips", "8.8.8.8", "--max-response-mb", "2", "--output", "raw", "--output-file", "complete")
			if !strings.Contains(output, "Output saved to file:") {
				t.Fatal("larger response limit did not complete")
			}
			data, err := os.ReadFile(filepath.Join(home, "complete.json"))
			var result []struct{ Padding string }
			if err != nil || json.Unmarshal(data, &result) != nil || len(result) != 1 || len(result[0].Padding) != 1_100_000 {
				t.Fatal("larger response limit lost data")
			}
		}
	}
}

func TestCLILargeBulkHTTP(t *testing.T) {
	if os.Getenv("CLI_HTTP_TEST") != "1" {
		t.Skip("requires the local HTTP simulation")
	}
	home := configuredHome(t, "BULK_50K")
	ips, err := os.Create(filepath.Join(home, "ips.txt"))
	if err != nil {
		t.Fatal(err)
	}
	writer := bufio.NewWriter(ips)
	for i := 0; i < 50_000; i++ {
		fmt.Fprintf(writer, "2001:db8::%x\n", i)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	ips.Close()
	for _, command := range []string{"bulk-ip-geo", "bulk-ip-security"} {
		t.Run(command, func(t *testing.T) {
			output, err := os.Create(filepath.Join(home, command+".stdout"))
			if err != nil {
				t.Fatal(err)
			}
			defer output.Close()
			args := []string{command, "--file", "ips.txt", "--output", "raw"}
			city := "Test City"
			if command == "bulk-ip-geo" {
				args = append(args, "--lang", "ar")
				city = "مدينة"
			}
			process := newCLIProcess(home, args...)
			var stderr bytes.Buffer
			process.Stdout, process.Stderr = output, &stderr
			if err := process.Run(); err != nil {
				t.Fatalf("large bulk request failed: %v %s", err, stderr.String())
			}
			info, err := output.Stat()
			if err != nil || info.Size() < 200_000_000 || info.Size() > 256_000_000 {
				t.Fatalf("unexpected large response size: %v %v", info, err)
			}
			if _, err := output.Seek(0, io.SeekStart); err != nil {
				t.Fatal(err)
			}
			decoder := json.NewDecoder(output)
			if token, err := decoder.Token(); err != nil || token != json.Delim('[') {
				t.Fatalf("bulk output is not a JSON array: %v %v", token, err)
			}
			count := 0
			for decoder.More() {
				var result struct {
					IP       string
					Location struct{ City string }
					Company  struct{ Name string }
				}
				if err := decoder.Decode(&result); err != nil {
					t.Fatal(err)
				}
				if result.IP != fmt.Sprintf("2001:db8::%x", count) || result.Location.City != city || len(result.Company.Name) != 4000 {
					t.Fatalf("bulk result %d changed", count)
				}
				count++
			}
			if token, err := decoder.Token(); err != nil || token != json.Delim(']') || count != 50_000 {
				t.Fatalf("bulk output incomplete: count=%d token=%v err=%v", count, token, err)
			}
			if _, err := decoder.Token(); err != io.EOF {
				t.Fatalf("unexpected data after bulk results: %v", err)
			}
			t.Logf("Verified %d ordered results and %d output bytes with default limits", count, info.Size())
		})
	}
}
