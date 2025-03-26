package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/forkbombeu/didimo/pkg/OpenID4VP"
	"github.com/spf13/cobra"
)

func main() {
	var input string
	var defaultPath string
	var configPath string

	// Define the root command using Cobra
	var rootCmd = &cobra.Command{
		Use:   "parse-input",
		Short: "Parses the input string using OpenID4VP",
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure the input is provided
			if input == "" || defaultPath == "" || configPath == "" {
				fmt.Println("Error: Missing required arguments.")
				cmd.Usage()
				return
			}

			// Call the ParseInput function with the provided arguments
			result, err := OpenID4VP.ParseInput(input, defaultPath, configPath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Marshal the result to JSON and print it
			output, _ := json.MarshalIndent(result, "", "    ")
			fmt.Println(string(output))
		},
	}

	// Define the flags for the command
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "Input string (required)")
	rootCmd.Flags().StringVarP(&defaultPath, "default", "d", "", "Path to the default JSON file (required)")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the config JSON file (required)")

	// Mark the flags as required
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("default")
	rootCmd.MarkFlagRequired("config")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
