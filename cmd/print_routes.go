package cmd

import (
	"github.com/spf13/cobra"
)

var printRoutesCmd = &cobra.Command{
	Use:   "print-routes",
	Short: "Prints out all registered routes",
	Run: func(cmd *cobra.Command, args []string) {
		initServer(true)
	},
}
