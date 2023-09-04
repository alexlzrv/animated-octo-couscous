package greetings

import (
	"os"
	"text/template"
)

const greetingsTemplate = `
Build version: <{{.BuildVersion}}>
Build date: <{{.BuildDate}}>
Build commit: <{{.BuildCommit}}>
`

type Greetings struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

var tmpl = template.Must(template.New("greetings").Parse(greetingsTemplate))

func Hello(buildVersion, buildDate, buildCommit string) error {
	greetings := Greetings{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	if greetings.BuildVersion == "" {
		greetings.BuildVersion = "N/A"
	}
	if greetings.BuildDate == "" {
		greetings.BuildDate = "N/A"
	}
	if greetings.BuildCommit == "" {
		greetings.BuildCommit = "N/A"
	}

	return tmpl.Execute(os.Stdout, greetings)
}
