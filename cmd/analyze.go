/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("analyze called")
	},
}

func init() {
	analyzeCmd.Flags().StringP("config", "c", "", "Chemin du fichier de configuration JSON")
	analyzeCmd.Flags().StringP("output", "o", "report.json", "Chemin du fichier de sortie JSON")
	rootCmd.AddCommand(analyzeCmd)
}
