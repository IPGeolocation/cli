/*
Copyright © 2025 IPGeolocation.io <support@ipgeolocation.io>
For more details visit: https://ipgeolocation.io
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/IPGeolocation/cli/ascii"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ipgeolocation",
	Short:   "CLI for accessing IPGeolocation.io API endpoints",
	Version: "1.0.1",
	Long: `ipgeolocation is a command-line interface (CLI) tool for accessing all major endpoints of the IPGeolocation.io API. 
It allows you to retrieve IP geolocation data, timezone information, currency details, user-agent parsing, and more, directly from your terminal.

This tool is built using the Cobra framework and supports easy querying of IP data with simple subcommands.

Examples:

  # Get geolocation info of a specific IP
  ipgeolocation ipgeo --ip 8.8.8.8

  # Get timezone info
  ipgeolocation timezone --ip 8.8.8.8

  # Parse user agent string
  ipgeolocation parse-user-agent --user-agent "Mozilla/5.0 ..."

You must have a valid API key from ipgeolocation.io to use this tool. You can set your API key using the "ipgeolocation config --apikey=<your-key>" command.
`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ascii.GetAsciiArt())
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init sets up the root command with flags.
func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
