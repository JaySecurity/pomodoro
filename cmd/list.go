package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"pomodoro/service"
	"pomodoro/timer"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Timers",
	Long:  `List All Active Timers.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", service.SocketPath)
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}

		fmt.Fprintln(conn, "list")
		var timer timer.Timer
		decoder := json.NewDecoder(conn)
		for {
			err := decoder.Decode(&timer)
			if err == io.EOF {
				break
			}
			response := fmt.Sprintf("Timer %d: Duration: %s\n", timer.Id, timer.Duration)
			fmt.Println(response)
		}
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
