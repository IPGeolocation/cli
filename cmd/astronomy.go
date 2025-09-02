package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	encoding "net/url"

	"github.com/IPGeolocation/cli/internal/common"
	"github.com/IPGeolocation/cli/internal/config"
	"github.com/IPGeolocation/cli/internal/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var astronomyFlags common.AstronomyFlags

var astronomyCmd = &cobra.Command{
	Use:   "astronomy",
	Short: "Lookup astronomy information like sunrise, sunset, moon phase, etc.",
	Long: `The 'astronomy' command uses the ipgeolocation.io Astronomy API
to fetch astronomy-related data such as:

- Sunrise and sunset
- Solar noon
- Moonrise and moonset
- Moon phase
- Day length
- Timezone-based or coordinate-based location support

You can specify the location using IP, city name, coordinates, or timezone.

API Reference: https://ipgeolocation.io/astronomy-api.html
`,
	Example: `
  # Get astronomy data for an IP address
  ipgeolocation astronomy --ip=8.8.8.8

  # Use coordinates
  ipgeolocation astronomy --latitude=40.7128 --longitude=-74.0060

  # Pretty table output
  ipgeolocation astronomy --location="New York" --output=table

  # YAML output
  ipgeolocation astronomy --ip=1.1.1.1 --output=yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/astronomy"
		url := baseURL + "?apiKey=" + cfg.ApiKey
		if astronomyFlags.IP != "" {
			url += "&ip=" + astronomyFlags.IP
		}

		if astronomyFlags.Tz != "" {
			url += "&time_zone=" + astronomyFlags.Tz
		}

		if astronomyFlags.Language != "" {
			url += "&lang=" + astronomyFlags.Language
		}

		if astronomyFlags.Location != "" {
			url += "&location=" + encoding.QueryEscape(astronomyFlags.Location)
		}

		if astronomyFlags.Latitude != 0 {
			url += "&lat=" + fmt.Sprintf("%f", astronomyFlags.Latitude)
		}

		if astronomyFlags.Longitude != 0 {
			url += "&long=" + fmt.Sprintf("%f", astronomyFlags.Longitude)
		}

		if astronomyFlags.Elevation != 0 {
			url += "&elevation=" + fmt.Sprintf("%f", astronomyFlags.Elevation)
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching astronomy info:", err)
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

		switch astronomyFlags.Output {
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

	astronomyCmd.Flags().StringVar(&astronomyFlags.IP, "ip", "", "IPv4 or IPv6 address (e.g. 8.8.8.8)")
	astronomyCmd.Flags().StringVar(&astronomyFlags.Tz, "tz", "", "Timezone name (e.g. America/New_York)")
	astronomyCmd.Flags().StringVar(&astronomyFlags.Location, "location", "", "Location name (e.g. New York)")
	astronomyCmd.Flags().Float64Var(&astronomyFlags.Latitude, "latitude", 0, "Latitude (e.g. 37.7749)")
	astronomyCmd.Flags().Float64Var(&astronomyFlags.Longitude, "longitude", 0, "Longitude (e.g. -122.4194)")
	astronomyCmd.Flags().StringVar(&astronomyFlags.Language, "lang", "", "Language code (e.g. en)")
	astronomyCmd.Flags().Float64Var(&astronomyFlags.Elevation, "elevation", 0, "Elevation (e.g. 1000)")
	astronomyCmd.Flags().StringVar(&astronomyFlags.Output, "output", "pretty", "Output format: pretty, raw, table")

	rootCmd.AddCommand(astronomyCmd)
}
