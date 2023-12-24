/*
Copyright © 2023 Mehmet Kemal Gokcay <kmlgkcy.dev@gmail.com>
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/logger"
	"github.com/kmg7/fson/internal/server"
	"github.com/kmg7/fson/internal/validator"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start serving servers",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		auth.Init()
		if err := validator.Instantiate(); err != nil {
			logger.Fatal(err.Error())
		}
		go server.ConfigServerStart()
		go server.TransferServerStart()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		func() {
			<-c
			logger.Info("Close signal")
			logger.Info("Closing the application")
			os.Exit(1)
		}()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
