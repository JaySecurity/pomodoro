package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current timer",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop called")
		conn, err := net.Dial("unix", "/tmp/pomodoro.sock")
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}

		fmt.Fprintln(conn, "stop")
		var response string
		fmt.Fscan(conn, &response)
		fmt.Println(response)
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
