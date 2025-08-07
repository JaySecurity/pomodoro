package cmd

import (
	"os"
	"pomodoro/types"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "A pomodoro timer",
	Long:  `A pomodoro timer with options to start and run the pomodoro process in the background. Stop a timer in progress or restart a timer.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP("id", "i", "0", "Timer ID")
	rootCmd.PersistentFlags().StringP("name", "n", "", "Timer Name (Optional)")
}

func setFlags(cmd *cobra.Command) types.Flags {
	flags := types.Flags{}
	if n, err := cmd.Flags().GetString("name"); err == nil {
		flags.Name = n
	}
	if i, err := cmd.Flags().GetString("id"); err == nil {
		flags.Index = i
	}
	if d, err := cmd.Flags().GetString("duration"); err == nil {
		flags.Duration = d
	}
	return flags
}
