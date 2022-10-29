package cli

import (
	"fmt"
	"log"
	"nomad_coin/explorer"
	"nomad_coin/rest"
	"strings"

	"github.com/spf13/cobra"
)

var port int

var launchCmd = &cobra.Command{
	Use:   "launch <explorer/rest>",
	Short: "An application to simulate blockchain",
	Long: `
	This command can launch rest server or blockchain explorer
	Usage: launch <explorer/rest> --port=<port number>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Server to launch is not specified")
		}
		explorerOrRest := args[1]
		fmt.Println(explorerOrRest)
		switch strings.ToLower(explorerOrRest) {
		case "explorer":
			explorer.Start(port)
		case "rest":
			rest.Start(port)
		default:
			log.Fatal("Wrong input for server to launch")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Start() {
	launchCmd.Flags().IntVar(&port, "port", 4000, "port for opening server")
	cobra.CheckErr(launchCmd.Execute())
}
