// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package template

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gbr-ufs/tubenoisseur/cli/state"
)

func TestSortDomainCounts(t *testing.T) {
	tests := map[string]struct {
		domainCounts map[string]uint
		want         []source
	}{
		"ArjanCodes": {
			domainCounts: map[string]uint{"git.arjan.codes": 1},
			want:         []source{{Count: 1, Domain: "git.arjan.codes"}},
		},
		"bashbunni": {
			want: []source{},
		},
		"ChrisTitusTech": {
			want: []source{},
		},
		"TheLinuxEXP": {
			domainCounts: map[string]uint{
				"thecybersecguru.com":    1,
				"nerds.xyz":              1,
				"www.pcworld.com":        1,
				"www.phoronix.com":       1,
				"discourse.ubuntu.com":   1,
				"www.techradar.com":      2,
				"cybernews.com":          1,
				"store.steampowered.com": 1,
				"www.gamingonlinux.com":  2,
				"www.blender.org":        1,
				"www.anthropic.com":      1,
				"fedoramagazine.org":     3,
				"keepandroidopen.org":    1,
			},
			want: []source{
				{Count: 3, Domain: "fedoramagazine.org"},
				{Count: 2, Domain: "www.gamingonlinux.com"},
				{Count: 2, Domain: "www.techradar.com"},
				{Count: 1, Domain: "cybernews.com"},
				{Count: 1, Domain: "discourse.ubuntu.com"},
				{Count: 1, Domain: "keepandroidopen.org"},
				{Count: 1, Domain: "nerds.xyz"},
				{Count: 1, Domain: "store.steampowered.com"},
				{Count: 1, Domain: "thecybersecguru.com"},
				{Count: 1, Domain: "www.anthropic.com"},
				{Count: 1, Domain: "www.blender.org"},
				{Count: 1, Domain: "www.pcworld.com"},
				{Count: 1, Domain: "www.phoronix.com"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := SortDomainCounts(test.domainCounts)

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got:\n%+v\nwant:\n%+v", got, test.want)
			}
		})
	}
}

type diskFailureWriter struct{}

func (e diskFailureWriter) Write(p []byte) (int, error) {
	return 0, errors.New("disk failure")
}

func TestWriteReport(t *testing.T) {
	tests := map[string]struct {
		writer       io.Writer
		channelState state.ChannelState
		wantError    bool
	}{
		"WriteFailure": {
			writer:       diskFailureWriter{},
			channelState: state.ChannelState{ChannelHandle: "Error"},
			wantError:    true,
		},
		"ArjanCodes": {
			writer: &bytes.Buffer{},
			channelState: state.ChannelState{
				ChannelHandle: "ArjanCodes",
				DomainCounts:  map[string]uint{"git.arjan.codes": 1},
				ExcludedDomains: []string{
					"arjan.codes",
					"www.arjancodes.com",
					"discord.arjan.codes",
					"amzn.to",
				},
				LastVideoID: "rZpwFN_n2-g",
				Videos: []state.ProcessedVideo{
					{
						ID:        "rZpwFN_n2-g",
						Title:     "Your API Can’t Handle Real-World Integrations",
						Published: "2026-05-01T15:01:33+00:00",
						URLs:      []string{"https://git.arjan.codes/2026/apidata"},
					},
				},
			},
			wantError: false,
		},
		"bashbunni": {
			writer: &bytes.Buffer{},
			channelState: state.ChannelState{
				ChannelHandle: "bashbunni",
				DomainCounts:  map[string]uint{},
				LastVideoID:   "E0tbuDmIXxg",
				Videos: []state.ProcessedVideo{
					{
						ID:        "E0tbuDmIXxg",
						Title:     "The BEST intro to a new programming language (interactive, in-IDE) #learnprogramming #rustlang #bash #short",
						Published: "2026-05-03T13:01:56+00:00",
						URLs:      []string{},
					},
				},
			},
			wantError: false,
		},
		"ChrisTitusTech": {
			writer: &bytes.Buffer{},
			channelState: state.ChannelState{
				ChannelHandle: "ChrisTitusTech",
				ExcludedDomains: []string{
					"www.cttstore.com",
					"www.patreon.com",
					"www.twitch.tv",
					"christitus.com",
				},
				LastVideoID: "7VTWbFqySxg",
				Videos: []state.ProcessedVideo{
					{
						ID:        "7VTWbFqySxg",
						Title:     "I Ditched Chrome for This $60 Browser, but without paying",
						Published: "2026-05-02T11:15:00+00:00",
						URLs:      []string{},
					},
				},
			},
		},
		"TheLinuxEXP": {
			writer: &bytes.Buffer{},
			channelState: state.ChannelState{
				ChannelHandle: "TheLinuxEXP",
				DomainCounts: map[string]uint{
					"thecybersecguru.com":    1,
					"nerds.xyz":              1,
					"xint.io":                1,
					"www.pcworld.com":        1,
					"www.phoronix.com":       1,
					"discourse.ubuntu.com":   2,
					"www.techradar.com":      2,
					"cybernews.com":          1,
					"store.steampowered.com": 1,
					"www.gamingonlinux.com":  2,
					"fedoramagazine.org":     3,
					"keepandroidopen.org":    1,
				},
				ExcludedDomains: []string{
					"squarespace.com",
					"www.tuxedocomputers.com",
					"www.youtube.com",
					"www.patreon.com",
					"paypal.me",
					"liberapay.com",
					"the-linux-experiment.creator-spring.com",
				},
				LastVideoID: "upOaSBo5ygU",
				Videos: []state.ProcessedVideo{
					{
						ID:        "upOaSBo5ygU",
						Published: "2026-05-02T09:46:06+00:00",
						Title:     "Ubuntu under attack, Big flaw affects all Linux distros, Linux beats Windows - Linux Weekly News",
						URLs: []string{
							"https://thecybersecguru.com/news/massive-attack-ubuntu-canonical-313-team-extortion/",
							"https://nerds.xyz/2026/04/copy-fail-linux-root-exploit/",
							"https://xint.io/blog/copy-fail-linux-distributions",
							"https://www.pcworld.com/article/3123900/framework-new-linux-laptop-is-selling-faster-than-its-windows-one.html",
							"https://www.phoronix.com/review/ubuntu-2604-windows-11/7",
							"https://discourse.ubuntu.com/t/the-future-of-ai-in-ubuntu/81130",
							"https://discourse.ubuntu.com/t/the-future-of-ai-in-ubuntu/81130/41",
							"https://www.techradar.com/vpn/vpn-privacy-security/the-eus-age-verification-app-has-a-privacy-problem-and-it-may-be-more-than-just-a-bug-in-an-app",
							"https://cybernews.com/security/eu-age-verification-app-hack/",
							"https://store.steampowered.com/sale/steamcontroller",
							"https://www.techradar.com/computing/peripherals-accessories/valve-steam-controller-2026",
							"https://www.gamingonlinux.com/2026/04/valve-have-plans-for-the-steam-deck-2-plus-a-brief-steam-machine-steam-frame-update/",
							"https://www.blender.org/press/anthropic-joins-the-blender-development-fund-as-corporate-patron/",
							"https://www.anthropic.com/news/claude-for-creative-work",
							"https://www.gamingonlinux.com/2026/05/blender-change-the-anthropic-ai-funding-deal-with-discussions-planned-for-ai-policies/",
							"https://fedoramagazine.org/announcing-fedora-linux-44/",
							"https://fedoramagazine.org/whats-new-in-fedora-kde-plasma-desktop-44/",
							"https://fedoramagazine.org/fedora-asahi-remix-44-is-now-available/",
							"https://keepandroidopen.org/en/",
						},
					},
				},
			},
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := WriteReport(test.writer, test.channelState)

			if (err != nil) != test.wantError {
				t.Fatalf("unexpected error: %v", err)
			}

			if !test.wantError {
				got := test.writer.(*bytes.Buffer).String()
				want, err := os.ReadFile(filepath.Join("..", "..", "fixtures", name+".md"))

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got != string(want) {
					t.Fatalf("got:\n%q\nwant:\n%q", got, want)
				}
			}
		})
	}
}
