package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IPGeolocation/cli/v2/internal/common"
	"github.com/IPGeolocation/cli/v2/internal/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var myIPFlags common.MyIPFlags

var myIPCmd = &cobra.Command{
	Use:     "my-ip",
	Aliases: []string{"getip"},
	Short:   "Get your public IP address using ipgeolocation.io",
	Long: `The 'my-ip' command returns the public IP address of the machine making the request,
as seen by the IPGeolocation.io API.

This is useful in shell scripts, cron jobs, and CI/CD pipelines when you need to know
the public IP address of the current host.

Examples:

  # Get your public IP address
  ipgeolocation my-ip

  # Print only the IP address (useful in scripts)
  ipgeolocation my-ip --output text

Notes:
  - This command does not require an API key.
  `,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("https://api.ipgeolocation.io/getip")
		if err != nil {
			fmt.Println("Error fetching your IP address:", err)
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

		switch myIPFlags.Output {
		case "text":
			fmt.Println(result["ip"])
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
	myIPCmd.Flags().StringVar(&myIPFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml, text")

	rootCmd.AddCommand(myIPCmd)
}
