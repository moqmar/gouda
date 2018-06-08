package pipeline

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreos/go-semver/semver"
	"github.com/moqmar/gouda/gouda"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

type outputPlugin struct{}

var outputLogger = logging.MustGetLogger("output")

// Init is called once on every build.
func (p *outputPlugin) Init(step gouda.Step) {}

// Each is called for every file that is processed in a build.
func (p *outputPlugin) Each(step gouda.Step, file *gouda.File) {
	if file.OutputPath != "" {
		outputPath := filepath.Join(gouda.Root, viper.GetString("output"), file.OutputPath)
		outputLogger.Noticef("Writing %s to %s (%d bytes)", file.InputPath, outputPath, len(file.Content))
		err := os.MkdirAll(filepath.Dir(outputPath), 0755)
		gouda.AssertNil(err)
		err = ioutil.WriteFile(outputPath, []byte(file.Content), 0644)
		gouda.AssertNil(err)
	}
}

// Info returns the name and semantic version of the plugin.
func (p *outputPlugin) Info() (string, *semver.Version) {
	return "Write", semver.New("1.0.0")
}

func init() {
	gouda.Register(&outputPlugin{}, gouda.OnWrite, 0.1)
}
