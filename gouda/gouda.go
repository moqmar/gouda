package gouda

import (
	"sort"

	"github.com/coreos/go-semver/semver"
	logging "github.com/op/go-logging"
)

// Plugin defines the interface for Gouda plugins.
type Plugin interface {
	Init(step Step)
	Each(step Step, file *File)
	Info() (string, *semver.Version)
}

// File is a representation of an actual file in the filesystem, but in its current state during the Gouda pipeline.
type File struct {
	// InputPath is the path of the source file relative to the project root.
	InputPath string
	// OutputPath is the relative path where OnWrite should write the file.Content to. If it's empty, the file will be dismissed.
	OutputPath string
	// Type contains the MIME type of the file (e.g. text/markdown)
	Type string
	// Content contains the current state of the file (e.g. the Markdown code in AfterRender).
	Content string
	// Metadata contains additional metadata for the file - plugins might write information here to communicate with other steps or plugins.
	Metadata map[string]interface{}
	// Frontmatter contains the YAML data parsed from the file.
	Frontmatter map[string]interface{}
}

// Files contains the files used during the build and deploy pipeline.
var Files = []*File{}

// Root is the absolute project root path, without trailing slash.
var Root = "."

// Step describes the step a plugin is currently running in.
type Step int

const (
	// BeforeRead is executed before any file is read. It should be used to add other files.
	BeforeRead Step = 0
	// OnRead is supposed to read files from the hard drive to gouda.Files.
	OnRead Step = 1
	// AfterRead is executed after all files are read. It should be used to pre-process file metadata.
	AfterRead Step = 2
	// BeforeParse is executed before the files are parsed. It should be used to pre-process file contents.
	BeforeParse Step = 3
	// OnParse is supposed to parse files of a specific type, e.g. a Markdown parser for `.md` files.
	OnParse Step = 4
	// AfterParse is executed after the files are parsed. It should be used to post-process file contents.
	AfterParse Step = 5
	// BeforeRender is executed before the files are rendered using the HTML template. It should be used to pre-process the template.
	BeforeRender Step = 6
	// OnRender is supposed to render the parsed content into a HTML template.
	OnRender Step = 7
	// AfterRender is executed after the files are rendered using the HTML template. It should be used to post-process the rendered HTML code.
	AfterRender Step = 8
	// BeforeWrite is executed before the rendered HTML files are written back to the hard drive. It should be used to configure output paths and metadata.
	BeforeWrite Step = 9
	// OnWrite is supposed to write the rendered HTML files to the hard drive.
	OnWrite Step = 10
	// AfterWrite is executed after all output files have been written. It should be used for cleanup actions.
	AfterWrite Step = 11
	// BeforeDeploy is executed before the documentation has been deployed. It should be used for deployment-preprocessing, e.g. to rewrite URLs for sub-directories.
	BeforeDeploy Step = 12
	// OnDeploy is supposed to deploy the documentation according to the config file.
	OnDeploy Step = 13
	// AfterDeploy is executed after everything has been deployed. It should be used for e.g. webhooks.
	AfterDeploy Step = 14
)

// plugins contains the raw and unsorted registered plugins (and seemingly has a very complicated datatype so it can be sorted relatively easily later)
var plugins = map[Step]map[float64][]Plugin{}

// Plugins contains the registered plugins by name.
var Plugins = map[string]Plugin{}

// handlers contains the basic Init handlers (functions that aren't plugins) that are executed once after all of the step's plugins.
var handlers = map[Step][]func(){}

var pluginLogger = logging.MustGetLogger("plugin")

// Register adds a plugin to the build pipeline.
func Register(plugin Plugin, order Step, priority float64) {
	name, version := plugin.Info()
	pluginLogger.Noticef("Registering plugin for step %d: %s %s", order, name, version.String())
	if _, ok := plugins[order]; !ok {
		plugins[order] = map[float64][]Plugin{}
	}
	if _, ok := plugins[order][priority]; !ok {
		plugins[order][priority] = []Plugin{plugin}
	} else {
		plugins[order][priority] = append(plugins[order][priority], plugin)
	}
	Plugins[name] = plugin
}

// On adds an Init event handler to a specified step. The handlers specified that way will be executed once after all of the step's plugins are done.
func On(order Step, fn func()) {
	if _, ok := handlers[order]; !ok {
		handlers[order] = []func(){fn}
	} else {
		handlers[order] = append(handlers[order], fn)
	}
}

var pipelineLogger = logging.MustGetLogger("pipeline")

// RunPipeline runs the build pipeline with all loaded plugins against a specific root.
func RunPipeline(deploy bool) {
	// Step 0: sort the plugins by priority
	pipelineLogger.Debugf("Sorting plugins...")
	var sortedPlugins = map[Step][]Plugin{}
	for step, pluginMap := range plugins {
		pluginCount := 0
		priorities := make([]float64, len(pluginMap))
		for priority, pluginList := range pluginMap {
			pluginCount += len(pluginList)
			priorities = append(priorities, priority)
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(priorities)))
		for _, priority := range priorities {
			if len(pluginMap[priority]) > 0 {
				if sortedPlugins[step] == nil {
					sortedPlugins[step] = []Plugin{}
				}
				sortedPlugins[step] = append(sortedPlugins[step], pluginMap[priority]...)
			}
		}
	}

	// Step 1: prepare build pipeline
	steps := []Step{BeforeRead, OnRead, AfterRead, BeforeParse, OnParse, AfterParse, BeforeRender, OnRender, AfterRender, BeforeWrite, OnWrite, AfterWrite}

	// Step 2: prepare deployment pipeline
	if deploy {
		steps = append(steps, BeforeDeploy, OnDeploy, AfterDeploy)
	}

	// Step 3: run pipeline
	pipelineLogger.Noticef("Executing pipeline - got %d plugins...", len(Plugins))
	for _, step := range steps {
		for _, plugin := range sortedPlugins[step] {
			plugin.Init(step)
		}
		for _, plugin := range sortedPlugins[step] {
			for _, file := range Files {
				plugin.Each(step, file)
			}
		}
		for _, handler := range handlers[step] {
			handler()
		}
	}
}
