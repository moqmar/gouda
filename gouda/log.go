package gouda

import (
	"os"

	"github.com/op/go-logging"
)

var leveled logging.LeveledBackend

// SetLogLevel changes the default log level
func SetLogLevel(level logging.Level) {
	leveled.SetLevel(level, "")
}

func init() {

	backend := logging.NewLogBackend(os.Stdout, "", 0)

	format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{module}/%{level:.4s} %{color:reset} %{message}`)
	formatter := logging.NewBackendFormatter(backend, format)
	leveled = logging.AddModuleLevel(formatter)
	leveled.SetLevel(logging.NOTICE, "")

	logging.SetBackend(leveled)

}
