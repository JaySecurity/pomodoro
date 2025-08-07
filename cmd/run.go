package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"pomodoro/service"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use: "run",
	RunE: func(cmd *cobra.Command, args []string) error {
		return service.RunService()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the background pomodoro process",
	Long:  "Run the background pomodoro process",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("MYAPP_BACKGROUND") != "1" {
			// Not yet daemonized — fork ourselves
			execPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %w", err)
			}

			proc := exec.Command(execPath, append([]string{"run"}, args...)...)
			proc.Env = append(os.Environ(), "MYAPP_BACKGROUND=1")
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr

			if err := proc.Start(); err != nil {
				return fmt.Errorf("failed to fork background process: %w", err)
			}

			fmt.Printf("Started background process with PID %d\n", proc.Process.Pid)
			return nil
		}
		return nil
	},
}

// Shutdown command: stops the listener
var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown the background pomodoro process",
	Long:  "Shutdown the background pomodoro process",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", service.SocketPath)
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}
		fmt.Fprintln(conn, "shutdown")
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(runCmd, shutdownCmd, serveCmd)
}

func StartServiceDetached(args []string) error {
	if os.Getenv("MYAPP_BACKGROUND") != "1" {
		// Not yet daemonized — fork ourselves
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		proc := exec.Command(execPath, append([]string{"run"}, args...)...)
		proc.Env = append(os.Environ(), "MYAPP_BACKGROUND=1")
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr

		if err := proc.Start(); err != nil {
			return fmt.Errorf("failed to fork background process: %w", err)
		}

		fmt.Printf("Started background process with PID %d\n", proc.Process.Pid)
	}
	return nil
}
