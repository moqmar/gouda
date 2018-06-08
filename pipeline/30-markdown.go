package pipeline

import (
	"github.com/coreos/go-semver/semver"
	"github.com/microcosm-cc/bluemonday"
	"github.com/moqmar/gouda/gouda"
	"github.com/op/go-logging"
	"github.com/russross/blackfriday"
)

type markdownPlugin struct{}

var markdownLogger = logging.MustGetLogger("markdown")

// Init is called once on every build.
func (p *markdownPlugin) Init(step gouda.Step) {}

// Each is called for every file that is processed in a build.
func (p *markdownPlugin) Each(step gouda.Step, file *gouda.File) {
	if file.Type == "text/markdown" {
		markdownLogger.Noticef("Parsing markdown file: %s", file.InputPath)
		unsafe := blackfriday.Run([]byte(file.Content))
		file.Content = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	}
}

// Info returns the name and semantic version of the plugin.
func (p *markdownPlugin) Info() (string, *semver.Version) {
	return "Markdown", semver.New("1.0.0")
}

func init() {
	gouda.Register(&markdownPlugin{}, gouda.OnParse, 0.1)
}
