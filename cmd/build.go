package cmd

import (
	"github.com/moqmar/gouda/gouda"
	"github.com/spf13/cobra"

	// Register default pipeline + plugins
	_ "github.com/moqmar/gouda/pipeline"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate the documentation HTML files",
	//Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		print(gouda.Yellow + gouda.Bold + "➜ Building documentation...\n" + gouda.Reset)
		// TODO: improve CLI - progress, spinner, colors
		gouda.RunPipeline(false)
	},
}
