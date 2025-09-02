package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IPGeolocation/cli/internal/common"
	"github.com/IPGeolocation/cli/internal/config"
	"github.com/IPGeolocation/cli/internal/utils"

	encoding "net/url"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var timezoneFlags common.TimezoneFlags

var timezoneCmd = &cobra.Command{
	Use:   "timezone",
	Short: "Lookup timezone information using ipgeolocation.io",

	Long: `The 'timezone' command retrieves timezone information for a given IP address, 
location name, GPS coordinates, or airport code using the ipgeolocation.io API.

You can specify the location by one of the following:
  - IP address
  - City name
  - Latitude and longitude
  - Timezone name
  - IATA, ICAO, or UN/LOCODE identifiers

You can also optionally specify the language of the output (e.g. 'en', 'es', etc).

Examples:

  # Lookup timezone info by IP
  ipgeolocation timezone --ip 8.8.8.8

  # Lookup timezone by location
  ipgeolocation timezone --location "New York"

  # Lookup timezone using coordinates
  ipgeolocation timezone --latitude 37.7749 --longitude -122.4194

  # Lookup timezone using timezone name
  ipgeolocation timezone --tz "Asia/Tokyo"

  # Lookup timezone using IATA airport code
  ipgeolocation timezone --iata DXB

Note: 
  - You must have a valid API key configured using: ipgeolocation config --apikey=<your-key>
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/timezone"
		url := baseURL + "?apiKey=" + cfg.ApiKey
		if timezoneFlags.IP != "" {
			url += "&ip=" + timezoneFlags.IP
		}

		if timezoneFlags.IataCode != "" {
			url += "&iata_code=" + timezoneFlags.IataCode
		}

		if timezoneFlags.IcaoCode != "" {
			url += "&icao_code=" + timezoneFlags.IcaoCode
		}

		if timezoneFlags.LoCode != "" {
			url += "&lo_code=" + timezoneFlags.LoCode
		}

		if timezoneFlags.Language != "" {
			url += "&lang=" + timezoneFlags.Language
		}

		if timezoneFlags.Tz != "" {
			url += "&tz=" + timezoneFlags.Tz
		}

		if timezoneFlags.Location != "" {
			url += "&location=" + encoding.QueryEscape(timezoneFlags.Location)
		}

		if timezoneFlags.Latitude != 0 {
			url += "&lat=" + fmt.Sprintf("%f", timezoneFlags.Latitude)
		}

		if timezoneFlags.Longitude != 0 {
			url += "&long=" + fmt.Sprintf("%f", timezoneFlags.Longitude)
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching timezone info:", err)
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

		switch timezoneFlags.Output {
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
	timezoneCmd.Flags().StringVar(&timezoneFlags.IP, "ip", "", "IPv4 or IPv6 address (e.g. 8.8.8.8)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.Tz, "tz", "", "Timezone name (e.g. America/New_York)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.Location, "location", "", "Location name (e.g. New York)")
	timezoneCmd.Flags().Float64Var(&timezoneFlags.Latitude, "latitude", 0, "Latitude (e.g. 37.7749)")
	timezoneCmd.Flags().Float64Var(&timezoneFlags.Longitude, "longitude", 0, "Longitude (e.g. -122.4194)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.IataCode, "iata", "", "IATA code (e.g. DXB)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.IcaoCode, "icao", "", "ICAO code (e.g. KATL)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.LoCode, "lo", "", "LO code (e.g. DEBER)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.Language, "lang", "", "Language code (e.g. en)")
	timezoneCmd.Flags().StringVar(&timezoneFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml")

	rootCmd.AddCommand(timezoneCmd)
}
