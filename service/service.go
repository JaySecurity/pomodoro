package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"pomodoro/timer"
	"pomodoro/types"
	"time"
)

var SocketPath = "/tmp/pomodoro.sock"

func RunService() {
	os.Remove(SocketPath)
	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		log.Fatalf("Unable to listen on socket: %v", err)
	}

	defer listener.Close()

	go handleNotify()

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
	var flags types.Flags
	_, err := fmt.Fscan(conn, &cmd)
	if err != nil {
		fmt.Println(err)
	}
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&flags)
	if err != nil {
		fmt.Println(err)
	}

	switch cmd {
	case "start":
		if flags.Index == "0" {
			d, err := time.ParseDuration(flags.Duration)
			_, err = timer.NewTimer(d, flags)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Get Timer and Start
			timer := timer.GetTimer(flags.Index)
			timer.Start()
		}
		break
	case "stop":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		break
	case "pause":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		timer.UpdateCh <- "pause"
		break
	case "restart":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		break
	case "list":
		if flags.Index == "0" {
			timers := timer.GetTimers()
			encoder := json.NewEncoder(conn)
			for _, timer := range timers {
				encoder.Encode(timer)
			}
		} else {
			timer := timer.GetTimer(flags.Index)
			encoder := json.NewEncoder(conn)
			encoder.Encode(timer)
		}
		break
	default:
		fmt.Println("Unknown Command")
		break
	}
}

func handleNotify() {
	for {
		timer := <-timer.TimerCh
		fmt.Printf("Timer %s: %v\n", timer.Id, timer)
		cmd := exec.Command("zenity", "--question", fmt.Sprintf("--text=Timer %s has elapsed:\n Would you like to take a break?", timer.Id))
		if errors.Is(cmd.Err, exec.ErrDot) {
			cmd.Err = nil
		}
		err := cmd.Run()
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			switch exitCode {
			case 0:
				fmt.Println("Yes")
			case 1:
				fmt.Println("No")
			}
		}

	}
}
