package utils

import "testing"

func TestOutputFilePath(t *testing.T) {
	for _, tc := range []struct {
		name, extension, want string
	}{
		// A bare name gains the extension.
		{"results", ".csv", "results.csv"},
		{"results", ".json", "results.json"},
		{"out/results", ".csv", "out/results.csv"},

		// A name that already carries it is left alone, so it never doubles up.
		{"results.csv", ".csv", "results.csv"},
		{"results.json", ".json", "results.json"},
		{"out/results.json", ".json", "out/results.json"},

		// Matching ignores case, which matters on Windows.
		{"results.CSV", ".csv", "results.CSV"},
		{"results.Json", ".json", "results.Json"},

		// A different extension is kept rather than swapped: the caller's name
		// survives and no unrelated file is targeted.
		{"results.json", ".csv", "results.json.csv"},
		{"results.csv", ".json", "results.csv.json"},
		{"results.tar.gz", ".csv", "results.tar.gz.csv"},

		// Dots inside a name are not an extension.
		{"my.results", ".csv", "my.results.csv"},
		{"2026-01-01.report", ".json", "2026-01-01.report.json"},

		// Degenerate shapes still produce a usable path.
		{".csv", ".csv", ".csv"},
		{"results.", ".csv", "results..csv"},
	} {
		if got := OutputFilePath(tc.name, tc.extension); got != tc.want {
			t.Fatalf("OutputFilePath(%q, %q) = %q, want %q", tc.name, tc.extension, got, tc.want)
		}
	}
}
