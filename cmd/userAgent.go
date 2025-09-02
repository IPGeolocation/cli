/*
Copyright © 2025 IPGeolocation.io <support@ipgeolocation.iio>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/IPGeolocation/cli/internal/common"
	"github.com/IPGeolocation/cli/internal/config"
	"github.com/IPGeolocation/cli/internal/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var userAgentFlags common.ParseUserAgentFlags

// userAgentCmd represents the userAgent command
var userAgentCmd = &cobra.Command{
	Use:   "parse-user-agent",
	Short: "Parse User-Agent string using ipgeolocation.io",
	Long: `Parse and extract device, browser, engine, and OS details 
from a User-Agent string using ipgeolocation.io API.

Examples:
  ipgeolocation parse-user-agent --user-agent "Mozilla/5.0 ..."
  ipgeolocation parse-user-agent --user-agent "<UA>" --output yaml
  ipgeolocation parse-user-agent --user-agent "<UA>" --output table
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		if userAgentFlags.UserAgent == "" {
			fmt.Println("Please provide a user agent string using --user-agent")
			return
		}
		baseURL := "https://api.ipgeolocation.io/v2/user-agent"
		url := baseURL + "?apiKey=" + cfg.ApiKey

		payload := map[string]interface{}{
			"uaString": userAgentFlags.UserAgent,
		}
		body, err := utils.PostJSON(url, payload, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			fmt.Println("Error fetching user agent info:", err)
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		switch userAgentFlags.Output {
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
	userAgentCmd.Flags().StringVar(&userAgentFlags.UserAgent, "user-agent", "", "User Agent")
	userAgentCmd.Flags().StringVar(&userAgentFlags.Output, "output", "", "Output format:  raw, table, yaml")
	rootCmd.AddCommand(userAgentCmd)

}
