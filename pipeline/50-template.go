package pipeline

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/moqmar/gouda/gouda"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

type templatePlugin struct{}

var templateLogger = logging.MustGetLogger("template")

// Init is called once on every build.
func (p *templatePlugin) Init(step gouda.Step) {
	viper.SetDefault("template", "github.com/moqmar/gouda-mintlook")
	templateDir := strings.TrimRight(viper.GetString("template"), "\\/")
	if !strings.HasPrefix(templateDir, ".") { // We've got a git template!
		// TODO: implement templates from git repository
		panic("Git templates are not yet implemented")
	} else {
		var err error
		templateDir, err = filepath.Abs(templateDir)
		gouda.AssertNil(err)
		templateLogger.Noticef("Using local template at %s", templateDir)
	}

	// Add asset files to gouda build output
	assetsDir, err := os.Stat(templateDir + "/.assets")
	if err == nil && assetsDir.IsDir() {
		err = filepath.Walk(templateDir+"/.assets", func(path string, f os.FileInfo, err error) error {
			gouda.AssertNil(err)

			if f.IsDir() {
				templateLogger.Debugf("Asset: ignoring directory: %s", path)
				return nil
			}

			content, err := ioutil.ReadFile(path)
			gouda.AssertNil(err)

			templateLogger.Debugf("Asset: adding file with %d bytes: %s", len(content), path)
			gouda.Files = append(gouda.Files, &gouda.File{
				InputPath:  "<template asset>",
				OutputPath: gouda.Unwrap(filepath.Rel(templateDir+"/.assets", path)).(string),
				Content:    string(content),
			})

			return nil
		})
		gouda.AssertNil(err)
	}

	// TODO: Access the config from an asset?!

	// Add templates
	err = filepath.Walk(templateDir, func(path string, f os.FileInfo, err error) error {
		gouda.AssertNil(err)

		if f.IsDir() || strings.HasPrefix(path, templateDir+string(filepath.Separator)+".assets"+string(filepath.Separator)) {
			templateLogger.Debugf("Template: ignoring directory or asset file: %s", path)
			return nil
		}

		extension := filepath.Ext(strings.TrimSuffix(path, filepath.Ext(path))) // Second-level extension is the extension of the processed file
		mimetype := strings.TrimSuffix(strings.TrimSuffix(gouda.Unwrap(filepath.Rel(templateDir, path)).(string), filepath.Ext(path)), extension)
		templateLogger.Debugf("Template: adding template for %s: %s", mimetype, path)
		gouda.Templates[mimetype] = &gouda.Template{
			Template:  gouda.Unwrap(template.ParseFiles(path)).(*template.Template),
			Metadata:  map[string]interface{}{},
			Extension: extension,
		}

		return nil
	})
	gouda.AssertNil(err)
}

// Each is called for every file that is processed in a build.
func (p *templatePlugin) Each(step gouda.Step, file *gouda.File) {
	if template, ok := gouda.Templates[file.Type]; file.InputPath != "<template asset>" && ok {
		file.OutputPath = strings.TrimSuffix(file.InputPath, filepath.Ext(file.InputPath)) + template.Extension
		templateLogger.Debugf("Found template for type %s, processing file: %s (output path: %s)", file.Type, file.InputPath, file.OutputPath)
		file.Content = template.Render(file, map[string]interface{}{})
	}
}

// Info returns the name and semantic version of the plugin.
func (p *templatePlugin) Info() (string, *semver.Version) {
	return "Template", semver.New("1.0.0")
}

func init() {
	gouda.Register(&templatePlugin{}, gouda.OnRender, 0.1)
}
