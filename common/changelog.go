package common

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"text/template"
)

type Changelog struct {
	Version      string
	Additions    []ChangelogEntry
	Changes      []ChangelogEntry
	Removals     []ChangelogEntry
	Deprecations []ChangelogEntry
	Bugfixes     []ChangelogEntry
	Breaking     []ChangelogEntry
	Repo         string
}

type ChangelogEntry struct {
	Description string
	Link        string
}

// Careful changing this, as it creates quite a bit of work getting the cmd_test.go
// working.  It seems odd to use this template to generate the output within the
// test itself, though that would certainly make it "self-updating".
// I guess it could be that we define the PR Description markdown in the Mock, and
// then define Changelog{...} data inside the test suite and pass that to the template
// render and then compare.
// TODO: make this so.
const changelogTemplate = `## {{ .Version }}
{{- if or .Additions .Changes .Removals .Deprecations .Bugfixes .Breaking -}}
{{- with .Additions }}

### Additions
{{ range . }}
{{ if .Link }}#### {{ .Link }}{{ end }}

{{ .Description }}
{{- end }}{{- end }}
{{- with .Changes }}

### Changes
{{ range . }}
{{ if .Link }}#### {{ .Link }}{{ end }}

{{ .Description }}
{{- end }}{{- end }}
{{- with .Removals }}

### Removals
{{ range . }}
{{ if .Link }}#### {{ .Link }}{{ end }}

{{ .Description }}
{{- end }}{{- end }}
{{- with .Deprecations }}

### Deprecations
{{ range . }}
{{ if .Link }}#### {{ .Link }}{{ end }}

{{ .Description }}
{{- end }}{{- end }}
{{- with .Bugfixes }}

### Bug Fixes
{{ range . }}
{{ if .Link }}#### {{ .Link }}{{ end }}

{{ .Description }}
{{- end }}{{- end }}
{{- with .Breaking }}

### Breaking Changes
{{ range . }}
{{ if .Link }}#### {{ .Link }}{{ end }}

{{ .Description }}
{{- end }}{{- end }}{{- else }}

No changes for this release!{{ end }}
`

var changelogTmpl = template.Must(template.New("changelog").Parse(changelogTemplate))

func (c *Changelog) Template() ([]byte, error) {
	w := &bytes.Buffer{}
	if err := changelogTmpl.Execute(w, c); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (c *Changelog) WriteFile(path string) error {
	data, err := c.Template()
	if err != nil {
		return err
	}
	existingFile, err := ioutil.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if errors.Is(err, os.ErrNotExist) || len(existingFile) == 0 {
		return ioutil.WriteFile(path, data, 0644)
	}

	data = append(data, '\n')
	data = append(data, existingFile...)
	return ioutil.WriteFile(path, data, 0644)
}
