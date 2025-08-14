package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if outputFormat == "json" {
			result := map[string]string{
				"version": version,
				"commit":  commit,
				"date":    date,
			}
			outputResult(result)
		} else {
			fmt.Printf("gwstool version %s\n", version)
			fmt.Printf("commit: %s\n", commit)
			fmt.Printf("built: %s\n", date)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
