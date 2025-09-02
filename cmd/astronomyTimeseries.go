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

var astronomyTimeseriesFlags common.AstronomyTimeSeriesFlags

var astronomyTimeseriesCmd = &cobra.Command{
	Use:   "astronomy-timeseries",
	Short: "Lookup Astronomy Time Series information using ipgeolocation.io",
	Long: `The 'astronomy-timeseries' command retrieves time series data for astronomical events 
(such as sunrise, sunset, moonrise, and moonset) for a given date range 
using the IPGeolocation.io Astronomy Time Series API.

This is useful for tracking solar and lunar events over multiple days.

Examples:

  ipgeolocation astronomy-timeseries --location "New York" --start-date 2025-08-01 --end-date 2025-08-07

  ipgeolocation astronomy-timeseries --ip 8.8.8.8 --start-date 2025-08-01 --end-date 2025-08-07

  ipgeolocation astronomy-timeseries --latitude 37.7749 --longitude -122.4194 --start-date 2025-08-01 --end-date 2025-08-07


  `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil || cfg.ApiKey == "" {
			fmt.Println("API key not found. Please run: ipgeolocation config --apikey=<your-key>")
			return
		}

		baseURL := "https://api.ipgeolocation.io/v2/astronomy/timeSeries"
		url := baseURL + "?apiKey=" + cfg.ApiKey

		if astronomyTimeseriesFlags.DateStart == "" || astronomyTimeseriesFlags.DateEnd == "" {
			fmt.Println("Please provide both start and end dates.")
			return
		}

		if astronomyTimeseriesFlags.DateStart != "" {
			url += "&dateStart=" + astronomyTimeseriesFlags.DateStart
		}

		if astronomyTimeseriesFlags.DateEnd != "" {
			url += "&dateEnd=" + astronomyTimeseriesFlags.DateEnd
		}

		if astronomyTimeseriesFlags.IP != "" {
			url += "&ip=" + astronomyTimeseriesFlags.IP
		}

		if astronomyTimeseriesFlags.Language != "" {
			url += "&lang=" + astronomyTimeseriesFlags.Language
		}

		if astronomyTimeseriesFlags.Location != "" {
			url += "&location=" + encoding.QueryEscape(astronomyTimeseriesFlags.Location)
		}

		if astronomyTimeseriesFlags.Latitude != 0 {
			url += "&lat=" + fmt.Sprintf("%f", astronomyTimeseriesFlags.Latitude)
		}

		if astronomyTimeseriesFlags.Longitude != 0 {
			url += "&long=" + fmt.Sprintf("%f", astronomyTimeseriesFlags.Longitude)
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching astronomy time-series info:", err)
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

		switch astronomyTimeseriesFlags.Output {
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

	astronomyTimeseriesCmd.Flags().StringVar(&astronomyTimeseriesFlags.IP, "ip", "", "IPv4 or IPv6 address (e.g. 8.8.8.8)")
	astronomyTimeseriesCmd.Flags().StringVar(&astronomyTimeseriesFlags.Location, "location", "", "Location name (e.g. New York)")
	astronomyTimeseriesCmd.Flags().StringVarP(&astronomyTimeseriesFlags.DateStart, "start-date", "s", "", "Start date (e.g. 2023-01-01)")
	astronomyTimeseriesCmd.Flags().StringVarP(&astronomyTimeseriesFlags.DateEnd, "end-date", "e", "", "End date (e.g. 2023-12-31)")
	astronomyTimeseriesCmd.Flags().Float64Var(&astronomyTimeseriesFlags.Latitude, "latitude", 0, "Latitude (e.g. 37.7749)")
	astronomyTimeseriesCmd.Flags().Float64Var(&astronomyTimeseriesFlags.Longitude, "longitude", 0, "Longitude (e.g. -122.4194)")
	astronomyTimeseriesCmd.Flags().StringVar(&astronomyTimeseriesFlags.Language, "lang", "", "Language code (e.g. en)")
	astronomyTimeseriesCmd.Flags().StringVar(&astronomyTimeseriesFlags.Output, "output", "pretty", "Output format: pretty, raw, table, yaml")

	rootCmd.AddCommand(astronomyTimeseriesCmd)
}
