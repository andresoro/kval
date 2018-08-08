package cmd

import (
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Cli for kval server",
	Long:  "Start a prompt to interact with kval server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)
}

func run() {
	p := prompt.New(
		execute,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("kval"),
	)
	p.Run()
}

func execute(input string) {

}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "add"},
		{Text: "get"},
		{Text: "delete"},
	}
}
