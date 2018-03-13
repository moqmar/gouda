package cmd

import (
	"time"

	"github.com/moqmar/gouda/gouda"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate the documentation HTML files",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		print(gouda.Yellow + gouda.Bold + "➜ Building documentation…    ")
		for true {
			gouda.Progress()
			time.Sleep(100 * time.Millisecond)
		}
	},
}
