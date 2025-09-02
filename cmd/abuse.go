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

var abuseFlags common.AbuseFlags

var abuseCmd = &cobra.Command{
	Use:   "abuse",
	Short: "Lookup abuse/contact information for an IP address using ipgeolocation.io",

	Long: `The 'abuse' command allows you to retrieve abuse contact information for a given IP address 
using the IPGeolocation.io Abuse Contact API.

This can be useful for reporting malicious or suspicious IP activity to the appropriate network administrator.

Example:

  ipgeolocation abuse --ip 8.8.8.8`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/abuse"
		url := baseURL + "?apiKey=" + cfg.ApiKey
		if abuseFlags.IP != "" {
			url += "&ip=" + abuseFlags.IP
		}

		if len(abuseFlags.Excludes) > 0 {
			url += "&excludes=" + strings.Join(abuseFlags.Excludes, ",")
		}

		if len(abuseFlags.Fields) > 0 {
			url += "&fields=" + strings.Join(abuseFlags.Fields, ",")
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching abuse info:", err)
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

		switch abuseFlags.Output {
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
	abuseCmd.Flags().StringVar(&abuseFlags.IP, "ip", "", "IPv4 or IPv6 address (e.g. 8.8.8.8)")
	abuseCmd.Flags().StringSliceVar(&abuseFlags.Excludes, "exclude", []string{}, "Fields to exclude from the output")
	abuseCmd.Flags().StringSliceVar(&abuseFlags.Fields, "fields", []string{}, "Get Specific Fields to include in the output")
	abuseCmd.Flags().StringVar(&abuseFlags.Output, "output", "", "Output format: yaml, raw, table, yaml")

	rootCmd.AddCommand(abuseCmd)
}
