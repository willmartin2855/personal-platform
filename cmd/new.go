package cmd

import "github.com/spf13/cobra"

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Scaffold new resources",
	Long: `Create new projects or connected resources.

Subcommands target specific scaffolds (for example GitHub Pages sites).`,
}
