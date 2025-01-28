package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"pomodoro/service"
	"pomodoro/timer"
	"syscall"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the background pomodoro process",
	Long:  "Run the background pomodoro process",
	Run: func(cmd *cobra.Command, args []string) {
		// Create Signal Channel
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		go service.RunService()

		sig := <-sigs
		close(timer.TimerCh)
		fmt.Printf("Signal: %v\n", sig)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
