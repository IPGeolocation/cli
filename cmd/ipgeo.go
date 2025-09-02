package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/IPGeolocation/cli/internal/common"
	"github.com/IPGeolocation/cli/internal/config"
	"github.com/IPGeolocation/cli/internal/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ipgeoFlags common.IpgeoFlags

var ipgeoCmd = &cobra.Command{
	Use:   "ipgeo",
	Short: "Lookup IP geolocation information using ipgeolocation.io",

	Long: `The 'ipgeo' command retrieves IP geolocation data such as country, city, timezone, ISP, 
and coordinates using the ipgeolocation.io IPGeo API (v2).

You can query any public IP address or domain name, and customize the output using include/exclude/fields 
parameters as supported by the API.

Available output formats include: pretty (default), raw (JSON), table, and YAML.

Examples:

  # Get info about your current IP
  ipgeolocation ipgeo

  # Lookup a specific IP
  ipgeolocation ipgeo --ip 8.8.8.8

  # Lookup a domain with extra fields
  ipgeolocation ipgeo --ip google.com --include security,timezone,currency

  # Get only specific fields in YAML
  ipgeolocation ipgeo --ip 1.1.1.1 --fields ip,organization --output yaml

  # Exclude unnecessary fields
  ipgeolocation ipgeo --excludes currency,time_zone

Notes:
  - You must have a valid API key configured using: ipgeolocation config --apikey=<your_key>
  - If no --ip flag is provided, it defaults to your current IP address.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/ipgeo"
		url := baseURL + "?apiKey=" + cfg.ApiKey
		if ipgeoFlags.IP != "" {
			url += "&ip=" + ipgeoFlags.IP
		}

		if len(ipgeoFlags.Include) > 0 {
			url += "&include=" + strings.Join(ipgeoFlags.Include, ",")
		}

		if len(ipgeoFlags.Excludes) > 0 {
			url += "&excludes=" + strings.Join(ipgeoFlags.Excludes, ",")
		}

		if len(ipgeoFlags.Fields) > 0 {
			url += "&fields=" + strings.Join(ipgeoFlags.Fields, ",")
		}

		if ipgeoFlags.Language != "" {
			url += "&lang=" + ipgeoFlags.Language
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching IP Geolocation info:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			fmt.Printf("Error: %s\n", string(body))
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		switch ipgeoFlags.Output {
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
	},
}

func init() {
	ipgeoCmd.Flags().StringVar(&ipgeoFlags.IP, "ip", "", "IPv4 or IPv6 address or domain (e.g. 8.8.8.8, google.com)")
	ipgeoCmd.Flags().StringSliceVar(&ipgeoFlags.Include, "include", []string{}, "To include additional values in the output")
	ipgeoCmd.Flags().StringSliceVar(&ipgeoFlags.Excludes, "excludes", []string{}, "Fields to exclude from the output")
	ipgeoCmd.Flags().StringSliceVar(&ipgeoFlags.Fields, "fields", []string{}, "Get Specific Fields to include in the output")
	ipgeoCmd.Flags().StringVar(&ipgeoFlags.Language, "lang", "", "Language for the output")
	ipgeoCmd.Flags().StringVar(&ipgeoFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml")

	rootCmd.AddCommand(ipgeoCmd)
}
