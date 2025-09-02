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

var securityFlags common.IPSecurityFlags

var ipSecurityCmd = &cobra.Command{
	Use:   "ip-security",
	Short: "Lookup IP security and threat intelligence information using ipgeolocation.io",

	Long: `The 'ip-security' command queries the IPGeolocation.io IP Security API to retrieve 
security-related details about a given IP address.

The data includes threat intelligence such as proxy detection, Tor usage, VPN status, 
threat types (malware, phishing), bot activity, and more.

You can tailor the output using the --include, --exclude, --fields, and --lang flags.

Examples:

  # Get security info of a public IP
  ipgeolocation ip-security --ip 8.8.8.8

  # Get security info of a public IP with additional fields
  ipgeolocation ip-security --ip 8.8.8.8 --include=location,time_zone

Notes: 
  - You must have a valid API key configured using: ipgeolocation config --apikey=<your_key>
  - If no --ip flag is provided, it defaults to your current IP address.
  `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/security"
		url := baseURL + "?apiKey=" + cfg.ApiKey
		if securityFlags.IP != "" {
			url += "&ip=" + securityFlags.IP
		}

		if len(securityFlags.Include) > 0 {
			url += "&include=" + strings.Join(securityFlags.Include, ",")
		}

		if len(securityFlags.Excludes) > 0 {
			url += "&excludes=" + strings.Join(securityFlags.Excludes, ",")
		}

		if len(securityFlags.Fields) > 0 {
			url += "&fields=" + strings.Join(securityFlags.Fields, ",")
		}

		if securityFlags.Language != "" {
			url += "&lang=" + securityFlags.Language
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching ip security info:", err)
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

		switch securityFlags.Output {
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
	ipSecurityCmd.Flags().StringVar(&securityFlags.IP, "ip", "", "IPv4 or IPv6 address (e.g. 8.8.8.8)")
	ipSecurityCmd.Flags().StringSliceVar(&securityFlags.Include, "include", []string{}, "To include additional values in the output")
	ipSecurityCmd.Flags().StringSliceVar(&securityFlags.Excludes, "exclude", []string{}, "Fields to exclude from the output")
	ipSecurityCmd.Flags().StringSliceVar(&securityFlags.Fields, "fields", []string{}, "Get Specific Fields to include in the output")
	ipSecurityCmd.Flags().StringVar(&securityFlags.Language, "lang", "", "Language for the output")
	ipSecurityCmd.Flags().StringVar(&securityFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml")

	rootCmd.AddCommand(ipSecurityCmd)
}
