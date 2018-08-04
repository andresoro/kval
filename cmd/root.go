package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "kval",
	Short: "Simple key-value store",
	Long:  "Kval is a simple key-value store that supports conncurrency",
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

// Execute is the root cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
