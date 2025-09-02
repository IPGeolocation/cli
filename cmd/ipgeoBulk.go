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

var bulkIpgeoFlags common.BulkIpgeoFlags

var bulkIpgeoCmd = &cobra.Command{
	Use:   "bulk-ip-geo",
	Short: "Lookup bulk IP geolocation information using ipgeolocation.io",
	Long: `The 'bulk-ip-geo' command allows you to retrieve geolocation data for multiple IP addresses 
in a single request using the ipgeolocation.io Bulk IPGeo API.

You can query multiple IPv4 or IPv6 addresses at once by passing them with --ips 
or providing a text file with --file. Each line in the file should contain one IP address.

Examples:

  # Lookup 3 IP addresses
  ipgeolocation bulk-ip-geo --ips 8.8.8.8,1.1.1.1,192.30.253.112
  
  # Lookup IPs from a file
  ipgeolocation bulk-ip-geo --file=ips.txt

  # Lookup from file and include location/timezone
  ipgeolocation bulk-ip-geo --file=ips.txt --include=location,time_zone
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		// If --file is provided, read IPs from file
		if bulkIpgeoFlags.File != "" {
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
					bulkIpgeoFlags.IPs = append(bulkIpgeoFlags.IPs, ip)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
		}

		if len(bulkIpgeoFlags.IPs) == 0 {
			fmt.Println("Please provide at least one IP address using --ips or --file.")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/ipgeo-bulk"
		url := baseURL + "?apiKey=" + cfg.ApiKey

		if len(bulkIpgeoFlags.Include) > 0 {
			url += "&include=" + strings.Join(bulkIpgeoFlags.Include, ",")
		}
		if len(bulkIpgeoFlags.Excludes) > 0 {
			url += "&excludes=" + strings.Join(bulkIpgeoFlags.Excludes, ",")
		}
		if len(bulkIpgeoFlags.Fields) > 0 {
			url += "&fields=" + strings.Join(bulkIpgeoFlags.Fields, ",")
		}
		if bulkIpgeoFlags.Language != "" {
			url += "&lang=" + bulkIpgeoFlags.Language
		}

		payload := map[string]interface{}{
			"ips": bulkIpgeoFlags.IPs,
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

		switch bulkIpgeoFlags.Output {
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
		if bulkIpgeoFlags.OutputFile != "" {
			pretty, _ := json.MarshalIndent(result, "", "  ")
			err := os.WriteFile(bulkIpgeoFlags.OutputFile+".json", pretty, 0644)
			if err != nil {
				fmt.Println("Error writing JSON to file:", err)
				return
			}
			fmt.Println("Output saved to file:", bulkIpgeoFlags.OutputFile+".json")
		}
	},
}

func init() {
	bulkIpgeoCmd.Flags().StringSliceVar(&bulkIpgeoFlags.IPs, "ips", []string{}, "Comma-separated list of IPs")
	bulkIpgeoCmd.Flags().StringSliceVar(&bulkIpgeoFlags.Include, "include", []string{}, "To include additional values in the output")
	bulkIpgeoCmd.Flags().StringSliceVar(&bulkIpgeoFlags.Excludes, "exclude", []string{}, "Fields to exclude from the output")
	bulkIpgeoCmd.Flags().StringSliceVar(&bulkIpgeoFlags.Fields, "fields", []string{}, "Get Specific Fields to include in the output")
	bulkIpgeoCmd.Flags().StringVar(&bulkIpgeoFlags.Language, "lang", "", "Language for the output")
	bulkIpgeoCmd.Flags().StringVar(&bulkIpgeoFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml")
	bulkIpgeoCmd.Flags().StringVar(&bulkIpgeoFlags.File, "file", "", "Path to a text file containing IPs (one per line)")
	bulkIpgeoCmd.Flags().StringVar(&bulkIpgeoFlags.OutputFile, "output-file", "", "Save output to a file (JSON only)")

	rootCmd.AddCommand(bulkIpgeoCmd)
}
