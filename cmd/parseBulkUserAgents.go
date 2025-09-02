/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
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

var bulkUserAgentsFlags common.ParseBulkUserAgentFlags

var parseBulkUserAgentsCmd = &cobra.Command{
	Use:   "parse-bulk-user-agents",
	Short: "Parse multiple user agent strings using ipgeolocation.io",

	Long: `The 'parse-bulk-user-agents' command allows you to parse multiple user-agent strings in a single request 
using the ipgeolocation.io Bulk User-Agent Parsing API.

It returns structured metadata for each user agent string, including device type, operating system, browser, 
browser version, and more. This is useful for analyzing logs, headers, and traffic patterns at scale.

Examples:

  # Parse multiple user agent strings
  ipgeolocation parse-bulk-user-agents --user-agents "Mozilla/5.0 (Windows NT 10.0; Win64; x64)","curl/7.64.1"

  # Output the parsed data as YAML
  ipgeolocation parse-bulk-user-agents --user-agents "Mozilla/5.0..." --output yaml

  # Display results in a table
  ipgeolocation parse-bulk-user-agents --user-agents "..." --output table

Note: 
  - You must have a valid API key configured using: ipgeolocation config --apikey=<your_key>

  `,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		if len(bulkUserAgentsFlags.UserAgents) == 0 {
			fmt.Println("Please provide at least one user agent.")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/user-agent-bulk"
		url := baseURL + "?apiKey=" + cfg.ApiKey

		payload := map[string]interface{}{
			"uaStrings": bulkUserAgentsFlags.UserAgents,
		}
		body, err := utils.PostJSON(url, payload, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			fmt.Println("Error fetching user agents info:", err)
			return
		}

		var result interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Invalid JSON:", err)
			return
		}
		switch bulkUserAgentsFlags.Output {
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
	parseBulkUserAgentsCmd.Flags().StringSliceVar(&bulkUserAgentsFlags.UserAgents, "user-agents", []string{}, "User Agents")
	parseBulkUserAgentsCmd.Flags().StringVar(&bulkUserAgentsFlags.Output, "output", "", "Output format:  raw, table, yaml")
	rootCmd.AddCommand(parseBulkUserAgentsCmd)

}
