package cmd

import (
	"fmt"

	"github.com/IPGeolocation/cli/internal/config"
	"github.com/IPGeolocation/cli/internal/utils"

	"github.com/spf13/cobra"
)

var apikey string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set or show API key configuration",
	Long: `The 'config' command is used to set or retrieve your saved ipgeolocation.io API key.

You can securely save your API key for future use, so you don't need to pass it with every command.
If no flag is passed, the currently stored API key will be displayed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if apikey != "" {
			encrypted, err := utils.EncryptString(apikey)
			if err != nil {
				fmt.Println("❌ Failed to encrypt API key:", err)
				return
			}
			cfg := config.Config{ApiKey: encrypted}
			if err := config.Save(cfg); err != nil {
				fmt.Println("❌ Failed to save config:", err)
				return
			}
			fmt.Println("✅ API key saved securely.")
		} else {
			cfg, err := config.Load()
			if err != nil {
				fmt.Println("❌ Failed to load config:", err)
				return
			}

			if cfg.ApiKey == "" {
				fmt.Println("⚠️  No API key configured.")
				return
			}
			// Mask all but last 5 characters
			masked := "********"
			if len(cfg.ApiKey) > 5 {
				masked += cfg.ApiKey[len(cfg.ApiKey)-5:]
			} else {
				masked += cfg.ApiKey
			}

			fmt.Println("🔐 Current API key:", masked)
		}
	},
}

func init() {
	configCmd.Flags().StringVar(&apikey, "apikey", "", "Set your ipgeolocation.io API key")
	rootCmd.AddCommand(configCmd)
}
