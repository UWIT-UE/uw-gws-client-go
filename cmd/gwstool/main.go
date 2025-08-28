package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/uwit-ue/uw-gws-client-go/gws"
)

// Config represents the configuration for gwstool
type Config struct {
	APIUrl     string
	CAFile     string
	ClientCert string
	ClientKey  string
	Timeout    int
}

var (
	configFile   string
	outputFormat string
	interactive  bool
	config       *Config
	gwsClient    *gws.Client
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gwstool",
	Short: "CLI tool for University of Washington Groups Web Service",
	Long: `gwstool is a command-line interface for the University of Washington Groups Web Service.
It provides access to group management, membership operations, and search functionality.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for commands that don't need it
		if cmd.Name() == "help" || cmd.Name() == "version" ||
			(cmd.Parent() != nil && cmd.Parent().Name() == "config") {
			return nil
		}
		return initializeClient()
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.config/gwstool/config)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "output format (text|json)")
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "enable interactive prompts")

	// Add subcommands
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(memberCmd)
	rootCmd.AddCommand(searchCmd)
}

func initializeClient() error {
	var err error
	config, err = loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	gwsConfig := &gws.Config{
		APIUrl:     config.APIUrl,
		CAFile:     config.CAFile,
		ClientCert: config.ClientCert,
		ClientKey:  config.ClientKey,
		Timeout:    30,
	}

	if config.Timeout > 0 {
		gwsConfig.Timeout = time.Duration(config.Timeout) * time.Second
	}

	gwsClient, err = gws.NewClient(gwsConfig)
	if err != nil {
		return fmt.Errorf("failed to create GWS client: %w", err)
	}

	if gwsClient == nil {
		return fmt.Errorf("GWS client is nil after creation")
	}

	return nil
}

func loadConfig() (*Config, error) {
	configPath := getActiveConfigPath()

	// Default config with just API URL and timeout
	cfg := &Config{
		APIUrl:  "https://groups.uw.edu/group_sws/v3",
		Timeout: 30,
	}

	// If no config file found, return error since credentials are required
	if configPath == "" {
		return nil, fmt.Errorf("no configuration file found. Please create a config file at ~/.config/gwstool/config or /etc/gwstool.conf with your GWS credentials. Use 'gwstool config init' to create one")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "api_url":
			cfg.APIUrl = value
		case "ca_file":
			cfg.CAFile = value
		case "client_cert":
			cfg.ClientCert = value
		case "client_key":
			cfg.ClientKey = value
		case "timeout":
			if timeout, err := strconv.Atoi(value); err == nil {
				cfg.Timeout = timeout
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Validate that required credentials are present
	if cfg.ClientCert == "" || cfg.ClientKey == "" {
		return nil, fmt.Errorf("configuration file %s must contain both client_cert and client_key for authentication", configPath)
	}

	return cfg, nil
}

func getConfigPath() string {
	if configFile != "" {
		return configFile
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to find home directory")
	}

	return filepath.Join(homeDir, ".config", "gwstool", "config")
}

func getSystemConfigPath() string {
	return "/etc/gwstool.conf"
}

func getActiveConfigPath() string {
	// Priority order:
	// 1. Command line specified config file
	// 2. User config file (~/.config/gwstool/config)
	// 3. System config file (/etc/gwstool/config)

	if configFile != "" {
		return configFile
	}

	userConfigPath := getConfigPath()
	if fileExists(userConfigPath) {
		return userConfigPath
	}

	systemConfigPath := getSystemConfigPath()
	if fileExists(systemConfigPath) {
		return systemConfigPath
	}

	// Return empty string if no config file found
	return ""
}

func getConfigInfo() map[string]interface{} {
	activeConfig := getActiveConfigPath()
	userConfig := getConfigPath()
	systemConfig := getSystemConfigPath()

	info := map[string]interface{}{
		"active_config": activeConfig,
		"config_sources": map[string]interface{}{
			"command_line": configFile,
			"user_config": map[string]interface{}{
				"path":   userConfig,
				"exists": fileExists(userConfig),
			},
			"system_config": map[string]interface{}{
				"path":   systemConfig,
				"exists": fileExists(systemConfig),
			},
		},
	}

	if activeConfig != "" {
		if activeConfig == configFile {
			info["source"] = "command_line"
		} else if activeConfig == userConfig {
			info["source"] = "user"
		} else if activeConfig == systemConfig {
			info["source"] = "system"
		}
	} else {
		info["source"] = "none"
		info["error"] = "No configuration file found"
	}

	return info
}

func formatEntityList(entities gws.EntityList) string {
	if len(entities) == 0 {
		return "None"
	}

	var result []string
	for _, entity := range entities {
		if entity.Name != "" {
			result = append(result, fmt.Sprintf("%s (%s)", entity.ID, entity.Name))
		} else {
			result = append(result, entity.ID)
		}
	}
	return strings.Join(result, ", ")
}

func outputResult(data interface{}) {
	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "text":
		switch v := data.(type) {
		case *gws.Group:
			fmt.Printf("Group: %s\n", v.ID)
			fmt.Printf("Display Name: %s\n", v.DisplayName)
			fmt.Printf("Description: %s\n", v.Description)
			fmt.Printf("Contact: %s\n", v.Contact)
			fmt.Printf("GID: %d\n", v.Gid)
			if v.DependsOn != "" {
				fmt.Printf("Membership dependency group: %s\n", v.DependsOn)
			} else {
				fmt.Printf("Membership dependency group: None\n")
			}
			fmt.Printf("Administrators: %s\n", formatEntityList(v.Admins))
			fmt.Printf("Member managers: %s\n", formatEntityList(v.Updaters))
			fmt.Printf("Subgroup creators: %s\n", formatEntityList(v.Creators))
			fmt.Printf("Membership viewers: %s\n", formatEntityList(v.Readers))
		case *gws.History:
			if v != nil && len(v.Data) > 0 {
				// Define fixed column widths
				dateWidth := 24 // Width for the date column
				userWidth := 16 // Width for the user column
				actWidth := 15  // Width for the activity column
				descWidth := 70 // Width for description column

				if groupHistoryLongOutput {
					dateWidth = 24
					userWidth = 30
					actWidth = 20
					descWidth = 0 // No truncation
				}

				// Print header
				fmt.Printf("%-*s %-*s %-*s %s\n",
					dateWidth, "TIMESTAMP",
					userWidth, "USER",
					actWidth, "ACTIVITY",
					"DESCRIPTION")
				fmt.Printf("%s %s %s %s\n",
					strings.Repeat("-", dateWidth),
					strings.Repeat("-", userWidth),
					strings.Repeat("-", actWidth),
					strings.Repeat("-", descWidth))

				// Print each history entry
				for _, entry := range v.Data {
					// Convert timestamp (in milliseconds) to human-readable format
					ts := time.UnixMilli(int64(entry.Timestamp)).Format("2006-01-02 15:04:05 MST")

					// Prepare user field (user + actAs if present)
					user := entry.User
					if entry.ActAs != "" {
						user = fmt.Sprintf("%s (%s)", entry.User, entry.ActAs)
					}

					activity := entry.Activity
					description := strings.ReplaceAll(entry.Description, "|", "")
					description = strings.ReplaceAll(description, "\n", " ")

					if !groupHistoryLongOutput {
						// Truncate fields if too long
						if len(user) > userWidth {
							user = user[:userWidth-3] + "..."
						}

						if len(activity) > actWidth {
							activity = activity[:actWidth-3] + "..."
						}

						// Wrap description to fit in the terminal
						if len(description) > descWidth {
							description = description[:descWidth-3] + "..."
						}
					}

					fmt.Printf("%-*s %-*s %-*s %s\n",
						dateWidth, ts,
						userWidth, user,
						actWidth, activity,
						description)
				}
			} else {
				fmt.Println("No history entries found")
			}
		case *gws.MemberList:
			if v != nil {
				for _, member := range *v {
					fmt.Println(member.ID)
				}
			}
		case []string:
			for _, item := range v {
				fmt.Println(item)
			}
		case string:
			fmt.Println(v)
		case bool:
			fmt.Println(v)
		case int:
			fmt.Println(v)
		default:
			fmt.Printf("%+v\n", v)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", outputFormat)
		os.Exit(1)
	}
}

func promptForInput(prompt string) string {
	fmt.Print(prompt + ": ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
