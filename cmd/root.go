package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "Un outil CLI pour l'analyse distribuée de logs",
	Long: `GoLog Analyzer est un outil en ligne de commande permettant d'analyser en parallèle 
plusieurs fichiers de logs provenant de différentes sources et de générer des rapports personnalisés.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur : %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Ici tu peux définir des flags globaux (persistents) si besoin
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (par défaut: $HOME/.loganalyzer.yaml)")
}
