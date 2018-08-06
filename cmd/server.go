package cmd

import (
	"github.com/andresoro/kval/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  "Start a TCP on port 7765 and host kval",
	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}
