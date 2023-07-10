package main

import (
	"log"
	"os"
	"path"
	"text/template"

	"github.com/go-git/go-git/v5"
)

func main() {
	dir, err := os.MkdirTemp("", "aws-sdk-go-v2")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: "https://github.com/aws/aws-sdk-go-v2.git",
	})

	if err != nil {
		log.Fatal(err)
	}

	entries, err := os.ReadDir(path.Join(dir, "service"))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("cmd/services.go")
	defer f.Close()

	packageTemplate.Execute(f, entries)
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.

package cmd

import (
{{- range . }}
    {{ if ne .Name "internal" }}
		"github.com/aws/aws-sdk-go-v2/service/{{ .Name }}"
	{{ end }}
{{- end }}
)

func initClients () []interface{} {
	clients := []interface{}{}
{{- range . }}
    {{ if ne .Name "internal" }}
		clients = append(clients,{{ .Name }}.Client{})
	{{ end }}
{{- end }}	
	return clients
}
`))
