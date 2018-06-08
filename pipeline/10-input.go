package pipeline

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/moqmar/go-gitignore"
	"github.com/op/go-logging"
	"github.com/spf13/viper"

	"github.com/coreos/go-semver/semver"
	"github.com/moqmar/gouda/gouda"
)

// inputPlugin is the OnRead plugin responsible for adding and reading files from the disk to the memory
type inputPlugin struct {
	// Include contains the files to *include* in gitignore format (so, the files that would be ignored by git get added by gouda, and the other way around).
	// The default is []string{"!.git/", "*.md", "*.svg", "*.png", "*.jpg", "*.jpeg"}
	// To add a file from another plugin, you can use `pipeline.Include("*.txt")` in gouda.BeforeRead, e.g. using an Init event handler with gouda.On()
	Include []string
}

var inputLogger = logging.MustGetLogger("input")

// Include adds a file to the inclusion list
func Include(what string) {
	gouda.Plugins["Input"].(*inputPlugin).Include = append(gouda.Plugins["Input"].(*inputPlugin).Include, what)
}

// Source: https://golangcode.com/get-the-content-type-of-file/
func getContentType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// Init is called once on every build.
func (p *inputPlugin) Init(step gouda.Step) {
	viper.SetDefault("output", "./gouda-output")

	gistr := strings.Join(p.Include, "\n")
	if gistr == "" {
		gistr = "*"
	}
	gistr += "\n!/" + gouda.Unwrap(filepath.Rel(gouda.Root, gouda.Unwrap(filepath.Abs(viper.GetString("output"))).(string))).(string)

	gi := gitignore.NewGitIgnoreFromReader(gouda.Root, strings.NewReader(gistr))
	inputLogger.Debugf("Using include file:\n%s", gistr)

	inputLogger.Debugf("Reading directory: %s", gouda.Root)
	err := filepath.Walk(gouda.Root, func(path string, f os.FileInfo, err error) error {
		gouda.AssertNil(err)

		if gi.Match(path, f.IsDir()) {

			// Get type
			mimetype := mime.TypeByExtension(filepath.Ext(path))
			if mimetype == "" {
				f, err := os.Open(path)
				gouda.AssertNil(f)
				defer f.Close()
				mimetype, err = getContentType(f)
				gouda.AssertNil(err)
			}
			mimetype = strings.Split(mimetype, ";")[0]

			// Add to gouda.Files
			inputPath := gouda.Unwrap(filepath.Rel(gouda.Root, path)).(string)
			inputLogger.Noticef("Adding file %s (%s)", inputPath, mimetype)
			gouda.Files = append(gouda.Files, &gouda.File{
				InputPath: inputPath,
				Type:      mimetype,
			})

		} else {
			inputLogger.Debugf("Not matching include file - not adding file: %s", gouda.Unwrap(filepath.Rel(gouda.Root, path)).(string))
		}
		return nil
	})
	gouda.AssertNil(err)
}

// Each is called for every file that is processed in a build.
func (p *inputPlugin) Each(step gouda.Step, file *gouda.File) {
	// Read the actual file
	buf, err := ioutil.ReadFile(file.InputPath)
	gouda.AssertNil(err)
	file.Content = string(buf)
}

// Info returns the name and semantic version of the plugin.
func (p *inputPlugin) Info() (string, *semver.Version) {
	return "Input", semver.New("1.0.0")
}

func init() {
	gouda.Register(&inputPlugin{
		Include: []string{"!.git/", "*.md", "*.svg", "*.png", "*.jpg", "*.jpeg"},
	}, gouda.OnRead, 0.1)
}
