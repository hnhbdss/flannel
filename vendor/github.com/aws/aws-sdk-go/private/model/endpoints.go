// +build codegen

package model

import (
	"bytes"
	"go/format"
	"io"
	"text/template"
)

// GenerateEndpoints writes a Go file to the given writer.
func GenerateEndpoints(endpoints interface{}, w io.Writer) error {
	tmpl, err := template.New("endpoints").Parse(t)
	if err != nil {
		return err
	}

	out := bytes.NewBuffer(nil)
	if err = tmpl.Execute(out, endpoints); err != nil {
		return err
	}

	b, err := format.Source(bytes.TrimSpace(out.Bytes()))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, bytes.NewReader(b))
	return err
}

const t = `
package endpoints

// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

type endpointStruct struct {
	Version   int
	Endpoints map[string]endpointEntry
}

type endpointEntry struct {
	Endpoint      string
	SigningRegion string
}

var endpointsMap = endpointStruct{
	Version: {{ .Version }},
	Endpoints: map[string]endpointEntry{
		{{ range $key, $entry := .Endpoints }}"{{ $key }}": endpointEntry{
			Endpoint:      "{{ $entry.Endpoint }}",
			{{ if ne $entry.SigningRegion "" }}SigningRegion: "{{ $entry.SigningRegion }}",
			{{ end }}
		},
		{{ end }}
	},
}
`