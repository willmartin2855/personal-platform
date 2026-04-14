package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "personal-platform",
	Short: "Personal IDP CLI",
	Long: `Scaffolds and manages personal infrastructure (GitHub Pages sites, secrets, and more).`,
}

// Execute runs the root command tree.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
}
