package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause the current timer",
	Run: func(cmd *cobra.Command, args []string) {
		id := cmd.Flag("id").Value.String()
		conn, err := net.Dial("unix", "/tmp/pomodoro.sock")
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}
		fmt.Fprintln(conn, "pause", id)
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
