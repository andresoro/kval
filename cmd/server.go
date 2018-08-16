package cmd

import (
	"github.com/andresoro/kval/kval"

	"github.com/andresoro/kval/server"
	"github.com/spf13/cobra"
)

var (
	store *kval.Store
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  "Start a TCP on port 7765 and host kval",
	Run: func(cmd *cobra.Command, args []string) {
		r := server.NewRPC(config.rpcPort, store)
		h := server.NewHTTP(config.httpPort, store)
		go h.Start()
		r.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	store, _ = kval.New(config.shardNum, config.duration)

}
