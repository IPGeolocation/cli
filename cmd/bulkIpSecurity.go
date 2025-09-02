package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/IPGeolocation/cli/internal/common"
	"github.com/IPGeolocation/cli/internal/config"
	"github.com/IPGeolocation/cli/internal/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var bulkSecurityFlags common.BulkIPSecurityFlags
var bulkIpSecurityCmd = &cobra.Command{
	Use:   "bulk-ip-security",
	Short: "Lookup Bulk IP Security information using ipgeolocation.io",
	Long: `The 'bulk-ip-security' command allows you to retrieve security information for multiple IP addresses 
using the IPGeolocation.io Bulk IP Security API. This is useful for bulk processing of IP addresses, 
retrieving security information for multiple IP addresses at once, and analyzing security data for a large number of IP addresses.
Use the --ips flag to specify the list of IP addresses you want to retrieve security information for.

Example Use Cases:

Retrieving security information for a list of IP addresses: ipgeolocation bulk-ip-security --ips "8.8.8.8,8.8.4.4"

Retrieving security information for a list of IP addresses with including additional fields: ipgeolocation bulk-ip-security --ips "8.8.8.8,8.8.4.4" --include=location,time_zone
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		// If --file is provided, read IPs from file
		if bulkSecurityFlags.File != "" {
			file, err := os.Open(bulkIpgeoFlags.File)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				ip := strings.TrimSpace(scanner.Text())
				if ip != "" {
					bulkSecurityFlags.IPs = append(bulkSecurityFlags.IPs, ip)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
		}

		if len(bulkSecurityFlags.IPs) == 0 {
			fmt.Println("Please provide at least one IP address using --ips or --file.")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/security-bulk"
		url := baseURL + "?apiKey=" + cfg.ApiKey

		if len(bulkSecurityFlags.Include) > 0 {
			url += "&include=" + strings.Join(bulkSecurityFlags.Include, ",")
		}

		if len(bulkSecurityFlags.Excludes) > 0 {
			url += "&excludes=" + strings.Join(bulkSecurityFlags.Excludes, ",")
		}

		if len(bulkSecurityFlags.Fields) > 0 {
			url += "&fields=" + strings.Join(bulkSecurityFlags.Fields, ",")
		}

		if bulkSecurityFlags.Language != "" {
			url += "&lang=" + bulkSecurityFlags.Language
		}
		payload := map[string]interface{}{
			"ips": bulkSecurityFlags.IPs,
		}

		body, err := utils.PostJSON(url, payload, map[string]string{"Content-Type": "application/json"})
		if err != nil {
			fmt.Println("Error fetching Bulk IP Security info:", err)
			return
		}

		var result interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Invalid JSON:", err)
			return
		}

		switch bulkSecurityFlags.Output {
		case "raw":
			fmt.Println(string(body))
		case "table":
			utils.PrintAsTable(result, 0)
		case "yaml":
			yamlData, err := yaml.Marshal(result)
			if err != nil {
				fmt.Println("Error converting to YAML:", err)
				return
			}
			fmt.Println(string(yamlData))
		default:
			pretty, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(pretty))
		}
		if bulkSecurityFlags.OutputFile != "" {
			pretty, _ := json.MarshalIndent(result, "", "  ")
			err := os.WriteFile(bulkSecurityFlags.OutputFile+".json", pretty, 0644)
			if err != nil {
				fmt.Println("Error writing JSON to file:", err)
				return
			}
			fmt.Println("Output saved to file:", bulkSecurityFlags.OutputFile+".json")
		}
	},
}

func init() {
	bulkIpSecurityCmd.Flags().StringSliceVar(&bulkSecurityFlags.IPs, "ips", []string{}, "IPs")
	bulkIpSecurityCmd.Flags().StringSliceVar(&bulkSecurityFlags.Include, "include", []string{}, "To include additional values in the output")
	bulkIpSecurityCmd.Flags().StringSliceVar(&bulkSecurityFlags.Excludes, "exclude", []string{}, "Fields to exclude from the output")
	bulkIpSecurityCmd.Flags().StringSliceVar(&bulkSecurityFlags.Fields, "fields", []string{}, "Get Specific Fields to include in the output")
	bulkIpSecurityCmd.Flags().StringVar(&bulkSecurityFlags.Language, "lang", "", "Language for the output")
	bulkIpSecurityCmd.Flags().StringVar(&bulkSecurityFlags.Output, "output", "pretty", "Output format: pretty, raw, table")
	bulkIpSecurityCmd.Flags().StringVar(&bulkSecurityFlags.File, "file", "", "Path to a text file containing IPs (one per line)")
	bulkIpSecurityCmd.Flags().StringVar(&bulkSecurityFlags.OutputFile, "output-file", "", "Save output to a file (JSON only)")

	rootCmd.AddCommand(bulkIpSecurityCmd)
}
