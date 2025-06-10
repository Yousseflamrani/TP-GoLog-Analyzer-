package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "LogAnalyzer - Outil d'analyse de logs distribuée",
	Long: `LogAnalyzer est un outil CLI en Go pour analyser des fichiers de logs 
provenant de diverses sources en parallèle avec gestion robuste des erreurs.

Développé pour aider les administrateurs système à centraliser l'analyse 
de multiples logs et en extraire des informations clés.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erreur: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Configuration globale si nécessaire
}
