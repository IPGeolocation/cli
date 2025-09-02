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

var timeConversionFlags common.TimeConversionFlags

var timeConversionCmd = &cobra.Command{
	Use:   "time-conversion",
	Short: "Convert time between different timezones or locations using ipgeolocation.io",

	Long: `The 'time-conversion' command allows you to convert time from one location or timezone 
to another using the ipgeolocation.io Timezone API.

You can specify the source and destination using:
  - Timezones (e.g., "America/New_York")
  - City names (location_from/location_to)
  - Latitude/longitude coordinates
  - Airport codes (IATA, ICAO)
  - UN/LOCODEs

If no source time is provided, the current time is used by default.

Examples:

  # Convert time between timezones
  ipgeolocation time-conversion --tz_from "America/New_York" --tz_to "Asia/Tokyo" --time "2025-08-07 15:00:00"

  # Convert time based on locations
  ipgeolocation time-conversion --location_from "Berlin" --location_to "Sydney"

  # Convert time using coordinates
  ipgeolocation time-conversion --lat_from 52.52 --long_from 13.405 --lat_to -33.8688 --long_to 151.2093

  # Convert time using airport codes
  ipgeolocation time-conversion --iata_from "JFK" --iata_to "LHR"

Notes:
  - You must configure your API key first using: ipgeolocation config --apikey=<your_key>
  - If multiple location inputs are provided, precedence may depend on the API's logic.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/timezone/convert"
		url := baseURL + "?apiKey=" + cfg.ApiKey

		if timeConversionFlags.TimezoneFrom != "" {
			url += "&tz_from=" + timeConversionFlags.TimezoneFrom
		}

		if timeConversionFlags.TimezoneTo != "" {
			url += "&tz_to=" + timeConversionFlags.TimezoneTo
		}

		if timeConversionFlags.LocationFrom != "" {
			url += "&location_from=" + encoding.QueryEscape(timeConversionFlags.LocationFrom)
		}

		if timeConversionFlags.LocationTo != "" {
			url += "&location_to=" + encoding.QueryEscape(timeConversionFlags.LocationTo)
		}

		if timeConversionFlags.LatitudeFrom != 0 {
			url += "&lat_from=" + fmt.Sprintf("%f", timeConversionFlags.LatitudeFrom)
		}

		if timeConversionFlags.LongitudeFrom != 0 {
			url += "&long_from=" + fmt.Sprintf("%f", timeConversionFlags.LongitudeFrom)
		}

		if timeConversionFlags.LatitudeTo != 0 {
			url += "&lat_to=" + fmt.Sprintf("%f", timeConversionFlags.LatitudeTo)
		}

		if timeConversionFlags.LongitudeTo != 0 {
			url += "&long_to=" + fmt.Sprintf("%f", timeConversionFlags.LongitudeTo)
		}

		if timeConversionFlags.IataCodeFrom != "" {
			url += "&iata_from=" + timeConversionFlags.IataCodeFrom
		}

		if timeConversionFlags.IataCodeTo != "" {
			url += "&iata_to=" + timeConversionFlags.IataCodeTo
		}

		if timeConversionFlags.IcaoCodeFrom != "" {
			url += "&icao_from=" + timeConversionFlags.IcaoCodeFrom
		}

		if timeConversionFlags.IcaoCodeTo != "" {
			url += "&icao_to=" + timeConversionFlags.IcaoCodeTo
		}

		if timeConversionFlags.LoCodeFrom != "" {
			url += "&locode_from=" + timeConversionFlags.LoCodeFrom
		}

		if timeConversionFlags.LoCodeTo != "" {
			url += "&locode_to=" + timeConversionFlags.LoCodeTo
		}

		if timeConversionFlags.Time != "" {
			url += "&time=" + timeConversionFlags.Time
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching time info:", err)
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

		switch timeConversionFlags.Output {
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
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.TimezoneFrom, "tz_from", "", "Timezone from")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.TimezoneTo, "tz_to", "", "Timezone to")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.LocationFrom, "location_from", "", "Location from")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.LocationTo, "location_to", "", "Location to")
	timeConversionCmd.Flags().Float64Var(&timeConversionFlags.LatitudeFrom, "lat_from", 0, "Latitude from")
	timeConversionCmd.Flags().Float64Var(&timeConversionFlags.LongitudeFrom, "long_from", 0, "Longitude from")
	timeConversionCmd.Flags().Float64Var(&timeConversionFlags.LatitudeTo, "lat_to", 0, "Latitude to")
	timeConversionCmd.Flags().Float64Var(&timeConversionFlags.LongitudeTo, "long_to", 0, "Longitude to")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.IataCodeFrom, "iata_from", "", "IATA code from")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.IataCodeTo, "iata_to", "", "IATA code to")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.IcaoCodeFrom, "icao_from", "", "ICAO code from")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.IcaoCodeTo, "icao_to", "", "ICAO code to")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.LoCodeFrom, "lo_from", "", "LO code from")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.LoCodeTo, "lo_to", "", "LO code to")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.Time, "time", "", "Time")
	timeConversionCmd.Flags().StringVar(&timeConversionFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml")

	rootCmd.AddCommand(timeConversionCmd)
}
