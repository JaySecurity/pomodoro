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
	_, err := fmt.Fscan(conn, &cmd)
	if err != nil {
		fmt.Println(err)
	}
	switch cmd {
	case "start":
		var dur string
		_, err = fmt.Fscan(conn, &dur)
		d, _ := time.ParseDuration(dur)
		_, err = timer.NewTimer(d)
		if err != nil {
			fmt.Println(err)
		}
		break
	case "stop":
		fmt.Println("Stop Called")
		break
	case "restart":
		fmt.Println("Restart Called")
		break
	case "list":
		fmt.Println("List Called")
		timers := timer.GetTimers()
		encoder := json.NewEncoder(conn)
		for _, timer := range timers {
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
		fmt.Printf("Timer %d: %v\n", timer.Id, timer)
		cmd := exec.Command("zenity", "--question", fmt.Sprintf("--text=Timer %d has elapsed:\n Would you like to take a break?", timer.Id))
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
