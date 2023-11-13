package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sc",
	Short: "欢迎你",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("伊甸园：https://setcreed.github.io")
	},
}

func init() {
	RootCmd.AddCommand(runCommand, execCommand)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
