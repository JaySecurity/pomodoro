package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restart called")
		conn, err := net.Dial("unix", "/tmp/pomodoro.sock")
		if err != nil {
			fmt.Println("Failed to connect to service:", err)
			return
		}

		fmt.Fprintln(conn, "restart")
		var response string
		fmt.Fscan(conn, &response)
		fmt.Println(response)
		conn.Close()
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
