/*
Copyright © 2023 Mehmet Kemal Gokcay <kmlgkcy.dev@gmail.com>
*/
package cmd

import (
	"github.com/kmg7/fson/internal/config"
	server "github.com/kmg7/fson/internal/server/config"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		server.StartConfigServer("localhost:8080")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
