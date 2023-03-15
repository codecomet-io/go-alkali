package commands

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

func parseTemplate(format string) (*template.Template, error) {
	// aliases is from https://github.com/containerd/nerdctl/blob/v0.17.1/cmd/nerdctl/fmtutil.go#L116-L126 (Apache License 2.0)
	aliases := map[string]string{
		"json": "{{json .}}",
	}
	if alias, ok := aliases[format]; ok {
		format = alias
	}
	// funcs is from https://github.com/docker/cli/blob/v20.10.12/templates/templates.go#L12-L20 (Apache License 2.0)
	funcs := template.FuncMap{
		"json": func(v interface{}) string {
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			enc.Encode(v)
			// Remove the trailing new line added by the encoder
			return strings.TrimSpace(buf.String())
		},
	}
	return template.New("").Funcs(funcs).Parse(format)
}
