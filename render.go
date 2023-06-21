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
	t, err := template.ParseFS(templates.FS, "go.tpl")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
