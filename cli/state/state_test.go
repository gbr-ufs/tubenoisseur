// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package state

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gbr-ufs/tubenoisseur/scraping"
)

func TestLoadState(t *testing.T) {
	tests := map[string]struct {
		want      ChannelState
		wantError bool
	}{
		"ArjanCodes": {
			want: ChannelState{
				ChannelHandle: "ArjanCodes",
				ExcludedDomains: []string{
					"arjan.codes",
					"www.arjancodes.com",
					"discord.arjan.codes",
					"amzn.to",
				},
				LastVideoID:  "rZpwFN_n2-g",
				DomainCounts: map[string]uint{"git.arjan.codes": 1},
				Videos: []ProcessedVideo{
					{
						ID:        "rZpwFN_n2-g",
						Title:     "Your API Can’t Handle Real-World Integrations",
						Published: "2026-05-01T15:01:33+00:00",
						URLs:      []string{"https://git.arjan.codes/2026/apidata"}},
				},
			},
		},
		"bashbunni": {
			want: ChannelState{
				ChannelHandle:   "bashbunni",
				ExcludedDomains: []string{},
				LastVideoID:     "E0tbuDmIXxg",
				DomainCounts:    map[string]uint{},
				Videos: []ProcessedVideo{
					{
						ID:        "E0tbuDmIXxg",
						Title:     "The BEST intro to a new programming language (interactive, in-IDE) #learnprogramming #rustlang #bash #short",
						Published: "2026-05-03T13:01:56+00:00",
						URLs:      []string{},
					},
				},
			},
		},
		"ChrisTitusTech": {
			want: ChannelState{
				ChannelHandle: "ChrisTitusTech",
				ExcludedDomains: []string{
					"www.cttstore.com",
					"www.patreon.com",
					"www.twitch.tv",
					"christitus.com",
				},
				DomainCounts: map[string]uint{},
				LastVideoID:  "7VTWbFqySxg",
				Videos: []ProcessedVideo{
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
			want: ChannelState{
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
				Videos: []ProcessedVideo{
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
		},
		"Malformed": {
			wantError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			json, err := os.OpenFile(
				filepath.Join("..", "..", "fixtures", name+".json"),
				os.O_RDONLY,
				0644,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, err := LoadState(json)

			if (err != nil) != test.wantError {
				t.Fatalf("Unexpect error: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestGetVideos(t *testing.T) {
	tests := map[string]struct {
		feed        []scraping.VideoEntry
		lastVideoID string
		want        []scraping.VideoEntry
	}{
		"Empty Feed": {
			feed: []scraping.VideoEntry{},
			want: []scraping.VideoEntry{},
		},
		"Up to Date": {
			feed: []scraping.VideoEntry{{
				ID:        "989vCMg3Qv8",
				Published: "2026-04-26T17:30:03+00:00",
				Title:     "One Maintainer Is All It Takes To Break A Project",
			}},
			lastVideoID: "989vCMg3Qv8",
			want:        []scraping.VideoEntry{},
		},
		"Catch Up": {
			feed: []scraping.VideoEntry{{
				ID:        "989vCMg3Qv8",
				Published: "2026-04-26T17:30:03+00:00",
				Title:     "One Maintainer Is All It Takes To Break A Project",
			},
				{
					ID:        "I4sP8K8c6yc",
					Published: "2026-04-23T17:30:04+00:00",
					Title:     "Oxygen Squared The KDE Design We Never Had",
				}},
			lastVideoID: "I4sP8K8c6yc",
			want: []scraping.VideoEntry{{
				ID:        "989vCMg3Qv8",
				Published: "2026-04-26T17:30:03+00:00",
				Title:     "One Maintainer Is All It Takes To Break A Project",
			}},
		},
		"Cold Start, Unlisted": {
			feed: []scraping.VideoEntry{
				{
					ID:        "989vCMg3Qv8",
					Published: "2026-04-26T17:30:03+00:00",
					Title:     "One Maintainer Is All It Takes To Break A Project",
				},
				{
					ID:        "I4sP8K8c6yc",
					Published: "2026-04-23T17:30:04+00:00",
					Title:     "Oxygen Squared The KDE Design We Never Had",
				}},
			lastVideoID: "589vGMg3Qv1",
			want: []scraping.VideoEntry{
				{
					ID:        "989vCMg3Qv8",
					Published: "2026-04-26T17:30:03+00:00",
					Title:     "One Maintainer Is All It Takes To Break A Project",
				},
				{
					ID:        "I4sP8K8c6yc",
					Published: "2026-04-23T17:30:04+00:00",
					Title:     "Oxygen Squared The KDE Design We Never Had",
				}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := GetVideos(test.feed, test.lastVideoID)

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %+v, want %+v", got, test.want)
			}
		})
	}
}
