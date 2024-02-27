package cmd

import (
	"fmt"
	"log"

	"github.com/kmg7/fson/internal/config"
	"github.com/kmg7/fson/internal/crypt"
	"github.com/kmg7/fson/internal/profiles"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes application for first use.",
	Long:  `Initializes application for first use.`,
	Run: func(cmd *cobra.Command, args []string) {
		setRootFlagsToEnv()
		hash, err := crypt.Instance(crypt.Options{BcryptCost: 8}).Bcrypt([]byte("admin"))
		if err != nil {
			log.Fatal("Init crypt error")
		}

		if err := config.Setup(); err != nil {
			fmt.Printf("Something went wrong while setting config up. %v\n", err.Error())
		}
		if err := profiles.Setup("admin", string(*hash)); err != nil {
			fmt.Printf("Something went wrong while setting profiles up. %v\n", err.Error())
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
