package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gwstool configuration",
	Long:  "Manage gwstool configuration including viewing and initializing config files",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		configInfo := getConfigInfo()

		if outputFormat == "json" {
			result := configInfo

			if config != nil {
				result["current_config"] = map[string]interface{}{
					"api_url":     config.APIUrl,
					"ca_file":     config.CAFile,
					"client_cert": config.ClientCert,
					"client_key":  config.ClientKey,
					"timeout":     config.Timeout,
				}
			}
			outputResult(result)
		} else {
			fmt.Printf("Configuration Sources:\n")
			fmt.Printf("  System config: %s (exists: %t)\n",
				configInfo["config_sources"].(map[string]interface{})["system_config"].(map[string]interface{})["path"],
				configInfo["config_sources"].(map[string]interface{})["system_config"].(map[string]interface{})["exists"])
			fmt.Printf("  User config:   %s (exists: %t)\n",
				configInfo["config_sources"].(map[string]interface{})["user_config"].(map[string]interface{})["path"],
				configInfo["config_sources"].(map[string]interface{})["user_config"].(map[string]interface{})["exists"])
			if configFile != "" {
				fmt.Printf("  Command line:  %s\n", configFile)
			}

			fmt.Printf("\nActive Configuration:\n")
			if activeConfig := configInfo["active_config"].(string); activeConfig != "" {
				fmt.Printf("  Source: %s (%s)\n", configInfo["source"], activeConfig)
			} else {
				fmt.Printf("  Source: %s - %s\n", configInfo["source"], configInfo["error"])
			}

			if config != nil {
				fmt.Printf("  API URL: %s\n", config.APIUrl)
				fmt.Printf("  CA File: %s\n", config.CAFile)
				fmt.Printf("  Client Cert: %s\n", config.ClientCert)
				fmt.Printf("  Client Key: %s\n", config.ClientKey)
				fmt.Printf("  Timeout: %d seconds\n", config.Timeout)
			}
		}
		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize user configuration file",
	Long:  "Initialize a user configuration file at ~/.config/gws-mod/config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := getConfigPath() // Always create user config

		if fileExists(configPath) {
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			if !overwrite {
				if interactive {
					response := promptForInput(fmt.Sprintf("User config file already exists at %s. Overwrite? (yes/no)", configPath))
					if response != "yes" {
						fmt.Println("Configuration initialization cancelled")
						return nil
					}
				} else {
					return fmt.Errorf("user config file already exists at %s, use --overwrite to replace it", configPath)
				}
			}
		}

		// Create directory if it doesn't exist
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Default config content
		configContent := `# gwstool user configuration file
# This file takes precedence over system configuration (/etc/gwstool.conf)
# Edit these values with your actual credentials and settings

# GWS API URL
api_url=https://groups.uw.edu/group_sws/v3

# Path to CA certificate file
ca_file=/path/to/ca.cert

# Path to client certificate file
client_cert=/path/to/client.cert

# Path to client private key file
client_key=/path/to/client.key

# Request timeout in seconds
timeout=30
`

		if interactive {
			fmt.Printf("Creating user config file at: %s\n", configPath)
			if apiURL := promptForInput("API URL [https://groups.uw.edu/group_sws/v3]"); apiURL != "" {
				configContent = fmt.Sprintf("api_url=%s\n", apiURL)
			}
			if caFile := promptForInput("CA File path"); caFile != "" {
				configContent += fmt.Sprintf("ca_file=%s\n", caFile)
			}
			if clientCert := promptForInput("Client certificate path"); clientCert != "" {
				configContent += fmt.Sprintf("client_cert=%s\n", clientCert)
			}
			if clientKey := promptForInput("Client key path"); clientKey != "" {
				configContent += fmt.Sprintf("client_key=%s\n", clientKey)
			}
			if timeout := promptForInput("Timeout in seconds [30]"); timeout != "" {
				configContent += fmt.Sprintf("timeout=%s\n", timeout)
			}
		}

		if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		if outputFormat == "json" {
			result := map[string]string{
				"status":      "created",
				"config_file": configPath,
				"type":        "user",
			}
			outputResult(result)
		} else {
			fmt.Printf("User configuration file created at: %s\n", configPath)
			fmt.Println("Please edit the file with your actual credentials")
			fmt.Println("This user config will take precedence over system configuration")
		}
		return nil
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		configInfo := getConfigInfo()
		activeConfigPath := configInfo["active_config"].(string)

		if activeConfigPath == "" {
			if outputFormat == "json" {
				result := map[string]interface{}{
					"valid":  false,
					"source": "default",
					"errors": []string{"no configuration file found"},
				}
				outputResult(result)
			} else {
				fmt.Println("No configuration file found")
				fmt.Printf("Checked locations:\n")
				if configFile != "" {
					fmt.Printf("  Command line: %s\n", configFile)
				}
				fmt.Printf("  User config:   %s\n", getConfigPath())
				fmt.Printf("  System config: %s\n", getSystemConfigPath())
				return fmt.Errorf("no configuration file found")
			}
			return nil
		}

		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Validate required fields
		errors := []string{}

		if cfg.ClientCert == "" {
			errors = append(errors, "client_cert is required")
		} else if !fileExists(cfg.ClientCert) {
			errors = append(errors, fmt.Sprintf("client_cert file not found: %s", cfg.ClientCert))
		}

		if cfg.ClientKey == "" {
			errors = append(errors, "client_key is required")
		} else if !fileExists(cfg.ClientKey) {
			errors = append(errors, fmt.Sprintf("client_key file not found: %s", cfg.ClientKey))
		}

		if cfg.CAFile != "" && !fileExists(cfg.CAFile) {
			errors = append(errors, fmt.Sprintf("ca_file not found: %s", cfg.CAFile))
		}

		if outputFormat == "json" {
			result := map[string]interface{}{
				"valid":  len(errors) == 0,
				"config": activeConfigPath,
				"source": configInfo["source"],
			}
			if len(errors) > 0 {
				result["errors"] = errors
			}
			outputResult(result)
		} else {
			fmt.Printf("Validating configuration from: %s (%s)\n", activeConfigPath, configInfo["source"])
			if len(errors) == 0 {
				fmt.Println("Configuration is valid")
			} else {
				fmt.Printf("Configuration has errors:\n")
				for _, err := range errors {
					fmt.Printf("  - %s\n", err)
				}
				return fmt.Errorf("configuration validation failed")
			}
		}
		return nil
	},
}

func init() {
	// Add subcommands to config command
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configValidateCmd)

	// Add flags
	configInitCmd.Flags().Bool("overwrite", false, "Overwrite existing config file")

	// Add config command to root
	rootCmd.AddCommand(configCmd)
}
