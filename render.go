package curl2go

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/lonjaju/curl2go/templates"
)

type Render struct {
	cmd *ParsedFlags
}

func NewRender() *Render {
	return &Render{}
}

func Curl2go(curl string) (string, error) {
	return NewRender().Curl2Go(curl)
}

func toTitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

func ParseHeaders(stringHeaders []string) map[string]string {
	headers := make(map[string]string)
	for _, stringHeader := range stringHeaders {
		split := strings.Index(stringHeader, ":")
		if split == -1 {
			continue
		}
		name := strings.TrimSpace(stringHeader[:split])
		value := strings.TrimSpace(stringHeader[split+1:])
		headers[toTitleCase(name)] = value
	}
	return headers
}

func (rd *Render) Curl2Go(curl string) (string, error) {
	flags := FlagParse(curl)
	r, err := ExtractRelevant(flags)
	if err != nil {
		return "", err
	}

	return rd.Render(r)
}

func (rd *Render) Render(r *Relevant) (string, error) {
	// Parse the "go.tpl" template from the templates.FS file system
	t, err := template.ParseFS(templates.FS, "go.tpl")
	if err != nil {
		// If there is an error parsing the template, return an empty string and the error
		return "", err
	}

	// Create a new buffer to store the rendered template
	buf := new(bytes.Buffer)

	// Execute the template, passing in the relevant data and write the output to the buffer
	err = t.Execute(buf, r)
	if err != nil {
		// If there is an error executing the template, return an empty string and the error
		return "", err
	}

	// Convert the buffer to a string and return it along with no error
	return buf.String(), nil
}
