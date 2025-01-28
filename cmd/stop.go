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
	Run: func(cmd *cobra.Command, args []string) {
		id := cmd.Flag("id").Value.String()
		conn, err := net.Dial("unix", "/tmp/pomodoro.sock")
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}
		fmt.Fprintln(conn, "stop", id)
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
