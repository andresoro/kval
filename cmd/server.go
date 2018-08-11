package cmd

import (
	"time"

	"github.com/andresoro/kval/server"
	"github.com/spf13/cobra"
)

var shardNum int
var duration time.Duration

func init() {
	rootCmd.AddCommand(serverCmd)
	shardNum = 4
	duration = 3 * time.Minute
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  "Start a TCP on port 7765 and host kval",
	Run: func(cmd *cobra.Command, args []string) {
		r := server.NewRPC(":8080", shardNum, duration)
		r.Start()
	},
}
