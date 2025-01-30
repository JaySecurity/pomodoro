package cmd

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current timer",
	Run: func(cmd *cobra.Command, args []string) {
		flags := setFlags(cmd)
		conn, err := net.Dial("unix", "/tmp/pomodoro.sock")
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}
		fmt.Fprintln(conn, "stop")
		encoder := json.NewEncoder(conn)
		encoder.Encode(flags)
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
