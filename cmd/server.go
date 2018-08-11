package cmd

import (
	"time"

	"github.com/andresoro/kval/kval"

	"github.com/andresoro/kval/server"
	"github.com/spf13/cobra"
)

var (
	shardNum int
	duration time.Duration
	store    *kval.Store
	httpPort string
	rpcPort  string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	rpcPort = ":7741"
	httpPort = ":8080"
	shardNum = 4
	duration = 3 * time.Minute
	store = kval.New(shardNum, duration)

}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  "Start a TCP on port 7765 and host kval",
	Run: func(cmd *cobra.Command, args []string) {
		r := server.NewRPC(rpcPort, store)
		h := server.NewHTTP(httpPort, store)
		go h.Start()
		r.Start()
	},
}
