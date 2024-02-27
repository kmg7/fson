package cmd

import (
	"os"

	"github.com/kmg7/fson/env"
	"github.com/kmg7/fson/internal/auth"
	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/profiles"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fson",
	Short: "File Sharing Over Network",
	Long: `fson is a program for making file transfer over http.
	More information can be found on https://github.com/kmg7/fson`,

	Run: func(cmd *cobra.Command, args []string) {
		setRootFlagsToEnv()
		config.Instance()
		auth.Instance()
		profiles.Instance()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setRootFlagsToEnv() error {
	if debugMode {
		if err := env.SetModeDebug(); err != nil {
			return err
		}
	}
	return nil
}

var debugMode bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "")
}
