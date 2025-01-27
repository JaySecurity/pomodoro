/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pomodoro/timer"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var socketPath = "/tmp/pomodoro.sock"

func runService() {
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Unable to listen on socket: %v", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var cmd string
	_, err := fmt.Fscan(conn, &cmd)
	if err != nil {
		fmt.Println(err)
	}
	switch cmd {
	case "start":
		fmt.Println("Start Called")
		err := timer.NewTimer(time.Minute * 25)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("New Timer Created: %v\n", timer.Timers)
		break
	case "stop":
		fmt.Println("Stop Called")
	case "restart":
		fmt.Println("Restart Called")
	default:
		fmt.Println("Unknown Command")
	}
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the background pomodoro process",
	Long:  "Run the background pomodoro process",
	Run: func(cmd *cobra.Command, args []string) {
		// Create Signal Channel
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		go runService()

		sig := <-sigs
		fmt.Printf("Signal: %v\n", sig)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
