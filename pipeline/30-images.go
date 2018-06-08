package pipeline

import (
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/moqmar/gouda/gouda"
	"github.com/op/go-logging"
)

type imagesPlugin struct{}

var imagesLogger = logging.MustGetLogger("images")

// Init is called once on every build.
func (p *imagesPlugin) Init(step gouda.Step) {}

// Each is called for every file that is processed in a build.
func (p *imagesPlugin) Each(step gouda.Step, file *gouda.File) {
	if strings.HasPrefix(file.Type, "image/") {
		imagesLogger.Noticef("Passing through image file: %s", file.InputPath)
		file.OutputPath = file.InputPath
	}
}

// Info returns the name and semantic version of the plugin.
func (p *imagesPlugin) Info() (string, *semver.Version) {
	return "Images", semver.New("1.0.0")
}

func init() {
	gouda.Register(&imagesPlugin{}, gouda.OnParse, 0.1)
}
