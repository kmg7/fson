package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	netutils "github.com/kmg7/fson/pkg/netutils"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists available network devices",
	Long: `List of available network addresses.
Name stands for the network interfaces name
IP stands for the address.

Local adresses such as "lo" or "localhost" are not accessible with
other devices in your network.
Make sure choosing correct address for your work.

`,
	Run: func(cmd *cobra.Command, args []string) {
		logAvailableIps()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func logAvailableIps() {
	ips, err := netutils.AvailableInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw)
	fmt.Fprintf(tw, "Name\tIP Adress\n")

	for _, ip := range ips {
		fmt.Fprintf(tw, "* %v\t%v\n",
			ip.Name,
			ip.Ip.String(),
		)
	}
	fmt.Fprintln(tw)
}
