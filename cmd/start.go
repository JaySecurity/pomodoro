package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"pomodoro/service"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer.",
	Run: func(cmd *cobra.Command, args []string) {
		flags := setFlags(cmd)
		conn, err := net.Dial("unix", service.SocketPath)
		if err != nil {
			if strings.Contains(
				err.Error(),
				"dial",
			) {
				fmt.Println("Starting Service...")
				StartServiceDetached(args)
				time.Sleep(time.Millisecond * 200)
				conn, err = net.Dial("unix", service.SocketPath)
				if err != nil {

					fmt.Println("Attempt 2: Failed to connect to service.", err)
					return
				}
			} else {
				fmt.Println("Attempt 1: Failed to connect to service.", err)
				return
			}
		}

		fmt.Fprintln(conn, "start")
		encoder := json.NewEncoder(conn)
		encoder.Encode(flags)
		conn.Close()
	},
}

func init() {
	startCmd.Flags().StringP("duration", "d", "25m", "Duration of the timer")
	rootCmd.AddCommand(startCmd)
}
