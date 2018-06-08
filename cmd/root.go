package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/moqmar/gouda/gouda"
)

var cliLogger = logging.MustGetLogger("cli")

// TODO: Custom, more colorful help template
var rootCmd = &cobra.Command{
	Use:   "gouda",
	Short: gouda.Bold + gouda.Yellow + "The no-worries documentation tool" + gouda.Reset,
	Long:  gouda.Cheese,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute starts the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cliLogger.Critical(err.Error())
		os.Exit(1)
	}
}

var cfgFile string
var debug bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file (default \"gouda.yaml\" in any parent directory)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print more detailed output")
	rootCmd.PersistentFlags().String("target", "./html", "target directory for generated documentation")
	viper.BindPFlag("target", rootCmd.PersistentFlags().Lookup("target"))
	//viper.SetDefault("target", "./html")
}

func initConfig() {
	if debug {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Walk up the file tree
		cwd, err := os.Getwd()
		if err == nil {
			n := strings.Count(strings.TrimSuffix(cwd, "/"), "/")
			for ; n > 0; n-- {
				p := ""
				for i := 0; i < n; i++ {
					p += "/.."
				}
				p = strings.TrimLeft(p, "/")
				viper.AddConfigPath(p)
			}
		}

		viper.AddConfigPath(".")
		viper.SetConfigName("gouda.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		cliLogger.Info("Couldn't find a configuration file, using the current directory with the default config.")
	} else {
		cliLogger.Noticef("Using the following configuration file: %s", viper.ConfigFileUsed())
		root, err := filepath.Abs(path.Dir(viper.ConfigFileUsed()))
		gouda.AssertNil(err)
		gouda.Root = root
		os.Chdir(gouda.Root)
	}

}
