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

var asnFlags common.ASNFlags

var asnCmd = &cobra.Command{
	Use:   "asn",
	Short: "Lookup ASN information using ipgeolocation.io",
	Long: `
	ASN API provides comprehensive details for an ASN including the as name, organization name, the country of registration, associated domain, and its type (ISP, host provider, or business). The API also shows the allocation date of provided ASN and if it is currently allocated or not. It also contains the routing information including peering, upstreams, and downstreams to help understand the relationship between different ASNs.

Example Use Cases:

Looking up ASN information for an IP address: ipgeolocation asn --ip 8.8.8.8

Retrieving ASN information for a specific ASN number ipgeolocation asn --asn 12345

Getting peering relationships for an ASN number: ipgeolocation asn --asn 12345 --include=peers
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/asn"
		url := baseURL + "?apiKey=" + cfg.ApiKey
		if asnFlags.IP != "" {
			url += "&ip=" + asnFlags.IP
		}

		if asnFlags.ASN != "" {
			url += "&asn=" + asnFlags.ASN
		}

		if len(asnFlags.Include) > 0 {
			url += "&include=" + strings.Join(asnFlags.Include, ",")
		}

		if len(asnFlags.Excludes) > 0 {
			url += "&excludes=" + strings.Join(asnFlags.Excludes, ",")
		}

		if len(asnFlags.Fields) > 0 {
			url += "&fields=" + strings.Join(asnFlags.Fields, ",")
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching ASN info:", err)
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

		switch asnFlags.Output {
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
	asnCmd.Flags().StringVar(&asnFlags.IP, "ip", "", "IPv4 or IPv6 address (e.g. 8.8.8.8)")
	asnCmd.Flags().StringVar(&asnFlags.ASN, "asn", "", "ASN number (e.g. 8075)")
	asnCmd.Flags().StringSliceVar(&asnFlags.Include, "include", []string{}, "To include additional values in the output")
	asnCmd.Flags().StringSliceVar(&asnFlags.Excludes, "exclude", []string{}, "Fields to exclude from the output")
	asnCmd.Flags().StringSliceVar(&asnFlags.Fields, "fields", []string{}, "Get Specific Fields to include in the output")
	asnCmd.Flags().StringVar(&asnFlags.Output, "output", "pretty", "Output format: pretty, raw, table")

	rootCmd.AddCommand(asnCmd)
}
