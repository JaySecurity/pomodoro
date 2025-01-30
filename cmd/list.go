package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"pomodoro/service"
	"pomodoro/types"
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
		decoder := json.NewDecoder(conn)
		for {
			var timer types.Timer
			var remaining time.Duration
			err := decoder.Decode(&timer)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
			}
			if timer.State == 1 {
				remaining = timer.Remaining - time.Since(timer.Started)
			} else {
				remaining = timer.Remaining
			}
			var response string
			if timer.Name == "" {
				response = fmt.Sprintf("Timer %s: Remaining: %v Status: %v\n", timer.Id, remaining, types.Status[timer.State])
			} else {
				response = fmt.Sprintf("%s: %s Timer - Remaining: %v Status: %v\n", timer.Id, timer.Name, remaining, types.Status[timer.State])
			}
			fmt.Println(response)
		}
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
