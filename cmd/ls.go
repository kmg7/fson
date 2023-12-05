/*
Copyright © 2023 Mehmet Kemal Gokcay <kmlgkcy.dev@gmail.com>
*/
package cmd

import (
	"fmt"
	"strings"

	netutils "github.com/kmg7/fson/pkg/netutils"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists available network devices",
	Long: `List of available network address.
Name stands for the network interfaces name
IP stands for the address.
You can also see the ip protocol version.

Local adresses such as "lo" or "localhost" are not accessible with
other devices in your network.
Make sure choosing correct address for your work.

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ls called")
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
	maxName, maxIp := 4, 7
	for _, ip := range ips {
		i := len(ip.Ip.String())
		n := len(ip.Name)
		if maxName < n {
			maxName = n
		}
		if maxIp < i {
			maxIp = i
		}
	}
	maxName += 4
	maxIp += 4
	fmt.Printf("%v\n", strings.Repeat("-", maxIp+maxName+9))
	fmt.Printf("Name%vIP%vis IPv4?\n", strings.Repeat(" ", maxName-3), strings.Repeat(" ", maxIp-2))
	fmt.Printf("%v\n", strings.Repeat("-", maxIp+maxName+9))
	for _, ip := range ips {
		fmt.Printf("-%v%v%v\n",
			ip.Name+strings.Repeat(" ", maxName-len(ip.Name)),
			ip.Ip.String()+strings.Repeat(" ", maxIp-len(ip.Ip.String())),
			ip.IsIPv4)
	}

}
