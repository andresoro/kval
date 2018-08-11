package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andresoro/kval/server"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var port string

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Cli for kval server",
	Long:  "Start a prompt to interact with kval server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Connecting to server")
		run()
	},
}

var client *server.Client

func init() {
	rootCmd.AddCommand(cliCmd)
	port = "7741"
}

func run() {
	client = &server.Client{
		Port: port,
	}
	err := client.Init()

	if err != nil {
		log.Fatal(err)
		return
	}

	p := prompt.New(
		execute,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("kval"),
	)
	p.Run()
}

func execute(input string) {
	in := strings.Fields(input)
	cmd := in[0]

	// handle commands
	switch cmd {

	case "add":
		key := in[1]
		val := strings.Join(in[2:], " ")
		if val == "" {
			fmt.Println("Need a value to go with a key")
			return
		}
		msg, err := client.Add(key, val)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)

	case "get":
		key := in[1]
		msg, err := client.Get(key)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(msg)

	case "delete":
		key := in[1]
		msg, err := client.Delete(key)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(msg)

	case "exit":
		os.Exit(0)

	default:
		fmt.Println("Not a proper command")
	}

}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}
