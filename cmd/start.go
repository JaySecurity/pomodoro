package cmd

import (
	"fmt"
	"net"
	"pomodoro/service"

	"github.com/spf13/cobra"
)

var duration string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer.",
	Run: func(cmd *cobra.Command, args []string) {
		d := cmd.Flag("duration").Value.String()
		conn, err := net.Dial("unix", service.SocketPath)
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}

		fmt.Fprintln(conn, "start", d)
		conn.Close()
	},
}

func init() {
	startCmd.Flags().StringVarP(&duration, "duration", "d", "25m", "Duration of the timer")
	rootCmd.AddCommand(startCmd)
}
