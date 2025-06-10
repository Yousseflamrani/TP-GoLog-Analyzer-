/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/Yousseflamrani/TP-GoLog-Analyzer-/internal/analyzer"
	"github.com/Yousseflamrani/TP-GoLog-Analyzer-/internal/config"
	"github.com/Yousseflamrani/TP-GoLog-Analyzer-/internal/reporter"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: analyzeRun,
}

func analyzeRun(cmd *cobra.Command, args []string) {
	configPath, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")

	logs, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Erreur de configuration : %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	results := make(chan analyzer.Result, len(logs))

	for _, log := range logs {
		wg.Add(1)
		go analyzer.ProcessLog(log, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var reports []reporter.Report
	for res := range results {
		var errorDetails string
		if res.Error != nil {
			errorDetails = res.Error.Error()
		}
		reports = append(reports, reporter.Report{
			LogID:        res.LogID,
			Status:       res.Status,
			ErrorDetails: errorDetails,
		})
		fmt.Printf("[%s] %s : %v\n", res.Status, res.LogID, res.Error)
	}

	if err := reporter.ExportJSON(reports, outputPath); err != nil {
		fmt.Printf("Erreur d'export : %v\n", err)
		os.Exit(1)
	}
}

func init() {
	analyzeCmd.Flags().StringP("config", "c", "", "Chemin du fichier de configuration JSON")
	analyzeCmd.Flags().StringP("output", "o", "report.json", "Chemin du fichier de sortie JSON")
	rootCmd.AddCommand(analyzeCmd)
}
