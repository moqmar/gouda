package gouda

import (
	"html/template"
	"strings"
)

// Template describes a HTML template for Gouda
type Template struct {
	Template  *template.Template
	Metadata  map[string]interface{}
	Extension string
}

// Templates contains rendering templates mapped by MIME type (e.g. gouda.Template["text/markdown"])
var Templates = map[string]*Template{}

// Render renders the template with the provided data
func (t *Template) Render(file *File, data map[string]interface{}) string {
	// Merge the metadata from the Template with the data provided directly to the function (which has a higher priority)
	mergedData := map[string]interface{}{}
	for k, v := range t.Metadata {
		mergedData[k] = v
	}
	for k, v := range file.Metadata {
		mergedData[k] = v
	}
	for k, v := range file.Frontmatter {
		mergedData[k] = v
	}
	for k, v := range data {
		mergedData[k] = v
	}
	mergedData["content"] = template.HTML(file.Content)
	mergedData["root"] = "."
	depth := strings.Count(file.OutputPath, "/")
	for i := 0; i < depth; i++ {
		mergedData["root"] = mergedData["root"].(string) + "/.."
	}
	mergedData["root"] = strings.TrimPrefix(mergedData["root"].(string), "./")

	// Render the template as a string
	s := &strings.Builder{}
	err := t.Template.Execute(s, mergedData)
	if err != nil {
		panic(err)
	}
	return s.String()
}
