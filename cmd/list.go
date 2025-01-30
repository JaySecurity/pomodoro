package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"pomodoro/service"
	"pomodoro/timer"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Timers",
	Long:  `List All Active Timers.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := setFlags(cmd)
		conn, err := net.Dial("unix", service.SocketPath)
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}

		fmt.Fprintln(conn, "list")
		encoder := json.NewEncoder(conn)
		encoder.Encode(flags)
		var timer timer.Timer
		decoder := json.NewDecoder(conn)
		for {
			fmt.Println("looping")
			err := decoder.Decode(&timer)
			if err == io.EOF {
				break
			}
			remaining := timer.Remaining - time.Since(timer.Started)
			var response string
			if timer.Name == "" {
				response = fmt.Sprintf("Timer %s: Duration: %v\n", timer.Id, remaining)
			} else {
				response = fmt.Sprintf("%s: %s Timer - Remaining: %v\n", timer.Id, timer.Name, remaining)
			}
			fmt.Println(response)
		}
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
