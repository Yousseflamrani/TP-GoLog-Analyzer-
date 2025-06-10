package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golog-analyzer",
	Short: "GoLog-Analyzer est un outil pour analyser les fichiers de logs.",
	Long:  `Un outil CLI en Go pour analyser, parser et extraire des informations utiles depuis des fichiers de logs configurés dans un fichier JSON.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}
}

func init() {
}
