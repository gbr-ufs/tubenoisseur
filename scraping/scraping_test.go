// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package scraping

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractVideos(t *testing.T) {
	tests := map[string]struct {
		want      []VideoEntry
		wantError bool
	}{
		"ArjanCodes": {
			want: []VideoEntry{
				{
					ID:        "rZpwFN_n2-g",
					Published: "2026-05-01T15:01:33+00:00",
					Title:     "Your API Can’t Handle Real-World Integrations",
				},
			},
			wantError: false,
		},
		"bashbunni": {
			want: []VideoEntry{
				{
					ID:        "E0tbuDmIXxg",
					Published: "2026-05-03T13:01:56+00:00",
					Title:     "The BEST intro to a new programming language (interactive, in-IDE) #learnprogramming #rustlang #bash",
				},
			},
			wantError: false,
		},
		"ChrisTitusTech": {
			want: []VideoEntry{
				{
					ID:        "7VTWbFqySxg",
					Published: "2026-05-02T11:15:00+00:00",
					Title:     "I Ditched Chrome for This $60 Browser, but without paying",
				},
			},
			wantError: false,
		},
		"TheLinuxEXP": {
			want: []VideoEntry{
				{
					ID:        "upOaSBo5ygU",
					Published: "2026-05-02T09:46:06+00:00",
					Title:     "Ubuntu under attack, Big flaw affects all Linux distros, Linux beats Windows - Linux Weekly News",
				},
			},
			wantError: false,
		},
		"Malformed": {
			want:      []VideoEntry{},
			wantError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			xml, err := os.OpenFile(filepath.Join("..", "fixtures", name+".xml"), os.O_RDONLY, 0644)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, err := ExtractVideos(xml)

			if (err != nil) != test.wantError {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("\ngot %+v\n want %+v", got, test.want)
			}
		})
	}
}

func TestExtractURLs(t *testing.T) {
	tests := map[string]struct {
		excludedDomains map[string]bool
		want            []string
	}{
		"ArjanCodes": {
			excludedDomains: map[string]bool{
				"arjan.codes":         true,
				"www.arjancodes.com":  true,
				"discord.arjan.codes": true,
				"amzn.to":             true,
			},
			want: []string{
				"https://git.arjan.codes/2026/apidata",
			},
		},
		"bashbunni": {
			excludedDomains: map[string]bool{},
			want:            []string{},
		},
		"ChrisTitusTech": {
			excludedDomains: map[string]bool{
				"www.cttstore.com": true,
				"www.patreon.com":  true,
				"www.twitch.tv":    true,
				"christitus.com":   true,
			},
			want: []string{},
		},
		"TheLinuxEXP": {
			excludedDomains: map[string]bool{
				"squarespace.com":                         true,
				"www.tuxedocomputers.com":                 true,
				"www.youtube.com":                         true,
				"www.patreon.com":                         true,
				"paypal.me":                               true,
				"liberapay.com":                           true,
				"the-linux-experiment.creator-spring.com": true,
			},
			want: []string{
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
		"InvalidHex": {
			// "Z" is not a valid hex digit. Despite that, xurls
			// finds this URL.
			want: []string{"https://example.com/path%20/good"},
		},
	}

	for name, test := range tests {
		t.Run(name, (func(t *testing.T) {
			t.Parallel()

			description, err := os.ReadFile(filepath.Join("..", "fixtures", name+".txt"))

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := ExtractURLs(string(description), test.excludedDomains), test.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %q, want %q", got, want)
			}
		}))
	}
}

type diskFailureReader struct{}

func (e diskFailureReader) Read(p []byte) (int, error) {
	return 0, errors.New("disk failure")
}

func TestExtractFeed(t *testing.T) {
	tests := map[string]struct {
		html      io.Reader
		want      string
		wantError bool
	}{
		"bashbunni": {
			html: bytes.NewBufferString(
				`"rssUrl":"http://www.youtube.com/feeds/videos.xml?channel_id=UC9H0HzpKf5JlazkADWnW1Jw"`,
			),
			want:      "http://www.youtube.com/feeds/videos.xml?channel_id=UC9H0HzpKf5JlazkADWnW1Jw",
			wantError: false,
		},
		"Bad Reader": {
			html:      diskFailureReader{},
			wantError: true,
		},
		"Malformed HTML": {
			html:      bytes.NewBufferString("This isn't even HTML"),
			wantError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ExtractFeed(test.html)

			if (err != nil) != test.wantError {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
