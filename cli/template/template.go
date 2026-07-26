// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Package template implements routines for writing out the Markdown file used
// by the application to present the JSON state in a more human-readable shape.

package template

import (
	"cmp"
	"io"
	"slices"
	"text/template"

	"github.com/gbr-ufs/tubenoisseur/cli/state"
)

const markdownTemplate = `# {{.ChannelHandle}}

Top Sources:
{{range .Sources}}
- <{{.Domain}}> = {{.Count}}
{{- end}}

{{range .Videos}}
## {{.Title}}

- ID: {{.ID}}
- Published: {{.Published}}

### References
{{ range .URLs}}
- <{{.}}>
{{- end}}
{{end}}`

var parsedTemplate = template.Must(template.New("Channel Report").Parse(markdownTemplate))

// SortDomainCounts sorts a domainCounts map into a list of sources. This is
// necessary because Go maps are unordered.
//
// It is used to show the top sources for a channel in order of most to least
// cited.
func SortDomainCounts(domainCounts map[string]uint) []source {
	sources := make([]source, 0, len(domainCounts))

	for domain, count := range domainCounts {
		sources = append(sources, source{Count: count, Domain: domain})
	}

	slices.SortFunc(sources, func(a, b source) int {
		if n := cmp.Compare(b.Count, a.Count); n != 0 {
			return n
		}

		return cmp.Compare(a.Domain, b.Domain)
	})

	return sources
}

// WriteReport writes out the Markdown file representing a YouTube channel and its sources.
func WriteReport(
	w io.Writer,
	channelState state.ChannelState,
) error {

	channelStateForTemplate := struct {
		ChannelHandle string
		Sources       []source
		Videos        []state.ProcessedVideo
	}{
		ChannelHandle: channelState.ChannelHandle,
		Sources:       SortDomainCounts(channelState.DomainCounts),
		Videos:        channelState.Videos,
	}

	return parsedTemplate.Execute(w, channelStateForTemplate)
}
