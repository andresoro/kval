package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config Config

var rootCmd = cobra.Command{
	Use:   "kval",
	Short: "Simple key-value store",
	Long:  "Kval is a simple key-value store that supports conncurrent read/write access",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run as either a server or client")
	},
}

func init() {
	config = DefaultConfig()
}

// Execute is the root cobra command
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
