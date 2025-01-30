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

	fmt.Println("Flags: ", flags.Index, flags.Duration, flags.Name)
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
		listTimers()
		break
	case "stop":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		timer := timer.GetTimer(flags.Index)
		fmt.Println(timer.Id, timer.Name, timer.Remaining, timer.State)
		timer.Stop()
		break
	case "pause":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		timer := timer.GetTimer(flags.Index)
		fmt.Println(timer.Id, timer.Name, timer.Remaining, timer.State)
		timer.Pause()
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
				encoder.Encode(&types.Timer{
					Id:        timer.Id,
					Name:      timer.Name,
					Remaining: timer.Remaining,
					Started:   timer.Started,
					State:     int(timer.State),
				})
			}
		} else {
			timer := timer.GetTimer(flags.Index)
			encoder := json.NewEncoder(conn)
			encoder.Encode(&types.Timer{
				Id:        timer.Id,
				Name:      timer.Name,
				Remaining: timer.Remaining,
				Started:   timer.Started,
				State:     int(timer.State),
			})
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

func listTimers() {
	timers := timer.GetTimers()
	for k, timer := range timers {
		var remaining time.Duration
		if timer.State == 1 {
			remaining = timer.Remaining - time.Since(timer.Started)
		} else {
			remaining = timer.Remaining
		}
		var response string
		if timer.Name == "" {
			response = fmt.Sprintf("Key: %s - Timer %s: Remaining: %v Status: %v\n", k, timer.Id, remaining, timer.State)
		} else {
			response = fmt.Sprintf("Key %s - %s: %s Timer - Remaining: %v Status: %v\n", k, timer.Id, timer.Name, remaining, timer.State)
		}
		fmt.Println(response)
	}
}
