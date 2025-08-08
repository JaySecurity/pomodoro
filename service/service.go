package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"pomodoro/timer"
	"pomodoro/types"
	"strings"
	"time"
)

var SocketPath = "/tmp/pomodoro.sock"

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func RunService() error {
	ctx, cancel = context.WithCancel(context.Background())
	// defer cancel()

	os.Remove(SocketPath)
	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		log.Fatalf("Unable to listen on socket: %v", err)
	}

	defer listener.Close()

	go handleNotify()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				log.Printf("Unable to accept connection: %v", err)
				continue
			}
			go handleConnection(conn)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down listener...")
	// close(timer.TimerCh)
	return nil
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

	// fmt.Println("Flags: ", flags.Index, flags.Duration, flags.Name)
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
	case "stop":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		timer := timer.GetTimer(flags.Index)
		// fmt.Println(timer.Id, timer.Name, timer.Remaining, timer.State)
		timer.Stop()
	case "pause":
		if flags.Index == "0" {
			flags.Index = "1"
		}
		timer := timer.GetTimer(flags.Index)
		// fmt.Println(timer.Id, timer.Name, timer.Remaining, timer.State)
		timer.Pause()
	case "restart":
		if flags.Index == "0" {
			flags.Index = "1"
		}
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
	case "shutdown":
		if cancel != nil {
			cancel()
		}
	}
}

func handleNotify() {
	for {
		t := <-timer.TimerCh
		// fmt.Printf("Timer %s: %v\n", t.Id, t)
		var stdout bytes.Buffer
		cmd := exec.Command(
			"yad",
			"--form",
			"--title=Choose an action",
			"--text=Timer Elapsed:",
			"--field=Choose:CB", "Break\\!Restart\\!Dismiss",
			":--buttons-layout=center",
			"--button=yad-ok",
			"--separator= ",
		)
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			fmt.Println("Yad failed:", err)
			return
		}
		choice := strings.TrimSpace(stdout.String())
		if choice == "Restart" {
			// fmt.Println("Restart")
			t.Remaining = t.Duration
			t.Start()
		} else if choice == "Break" {
			// fmt.Println("Yes")
			t.Remaining = time.Second * 5
			t.Start()
		} else {
			timer.DeleteTimer(t.Id)
		}
	}
}

func listTimers() {
	timers := timer.GetTimers()
	idx := 1
	for _, timer := range timers {
		var remaining time.Duration
		if timer.State == 1 {
			remaining = timer.Remaining - time.Since(timer.Started)
		} else {
			remaining = timer.Remaining
		}
		var response string
		if timer.Name == "" {
			response = fmt.Sprintf(
				"Key: %d - Timer %s: Remaining: %v Status: %v\n",
				idx,
				timer.Id,
				remaining,
				timer.State,
			)
		} else {
			response = fmt.Sprintf("Key %d: %s Timer - Remaining: %v Status: %v\n", idx, timer.Name, remaining, timer.State)
		}
		fmt.Println(response)
		idx++
	}
}
