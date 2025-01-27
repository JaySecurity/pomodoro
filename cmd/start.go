/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start Called with Args:", args)

		conn, err := net.Dial("unix", "/tmp/pomodoro.sock")
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}

		fmt.Fprintln(conn, "start")
		var response string
		fmt.Fscan(conn, &response)
		fmt.Println(response)
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
