// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gbr-ufs/tubenoisseur/cli/state"
	"github.com/gbr-ufs/tubenoisseur/cli/template"
)

type mockGetter struct {
	Error    error
	Response *http.Response
}

func (m *mockGetter) Get(url string) (*http.Response, error) {

	return m.Response, m.Error
}

type mockStorage struct {
	json     string
	markdown string
}

func (m *mockStorage) LoadState(channelHandle string) (state.ChannelState, error) {
	if m.json == "" || m.json == "{}" {
		return state.ChannelState{
			ChannelHandle: channelHandle,
		}, nil
	}
	return state.LoadState(strings.NewReader(m.json))
}

func (m *mockStorage) SaveState(channelHandle string, s state.ChannelState) error {
	return nil
}

func (m *mockStorage) WriteReport(channelHandle string, s state.ChannelState) error {
	var buf bytes.Buffer
	err := template.WriteReport(&buf, s)
	m.markdown = buf.String()
	return err
}

func TestNewRootCommand(t *testing.T) {
	tests := map[string]struct {
		args           []string
		channelHandle  string
		html           string
		htmlError      error
		htmlStatusCode int
		json           string
		wantError      bool
		xml            string
		xmlError       error
		xmlStatusCode  int
	}{
		"bashbunni": {
			args:           []string{"--debug"},
			channelHandle:  "bashbunni",
			html:           `"rssUrl":"http://www.youtube.com/feeds/videos.xml?channel_id=UC9H0HzpKf5JlazkADWnW1Jw"`,
			htmlStatusCode: http.StatusOK,
			json:           "{}",
			wantError:      false,
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns:media="http://search.yahoo.com/mrss/" xmlns="http://www.w3.org/2005/Atom">
 <link rel="self" href="http://www.youtube.com/feeds/videos.xml?channel_id=UC9H0HzpKf5JlazkADWnW1Jw"/>
 <id>yt:channel:9H0HzpKf5JlazkADWnW1Jw</id>
 <yt:channelId>9H0HzpKf5JlazkADWnW1Jw</yt:channelId>
 <title>bashbunni</title>
 <link rel="alternate" href="https://www.youtube.com/channel/UC9H0HzpKf5JlazkADWnW1Jw"/>
 <author>
  <name>bashbunni</name>
  <uri>https://www.youtube.com/channel/UC9H0HzpKf5JlazkADWnW1Jw</uri>
 </author>
 <published>2021-02-06T23:31:09+00:00</published>
 <entry>
  <id>yt:video:E0tbuDmIXxg</id>
  <yt:videoId>E0tbuDmIXxg</yt:videoId>
  <yt:channelId>UC9H0HzpKf5JlazkADWnW1Jw</yt:channelId>
  <title>The BEST intro to a new programming language (interactive, in-IDE) #learnprogramming #rustlang #bash</title>
  <link rel="alternate" href="https://www.youtube.com/shorts/E0tbuDmIXxg"/>
  <author>
   <name>bashbunni</name>
   <uri>https://www.youtube.com/channel/UC9H0HzpKf5JlazkADWnW1Jw</uri>
  </author>
  <published>2026-05-03T13:01:56+00:00</published>
  <updated>2026-05-03T13:10:07+00:00</updated>
  <media:group>
   <media:title>The BEST intro to a new programming language (interactive, in-IDE) #learnprogramming #rustlang #bash</media:title>
   <media:content url="https://www.youtube.com/v/E0tbuDmIXxg?version=3" type="application/x-shockwave-flash" width="640" height="390"/>
   <media:thumbnail url="https://i2.ytimg.com/vi/E0tbuDmIXxg/hqdefault.jpg" width="480" height="360"/>
   <media:description></media:description>
   <media:community>
    <media:starRating count="432" average="5.00" min="1" max="5"/>
    <media:statistics views="5944"/>
   </media:community>
  </media:group>
 </entry>
</feed>`,
			xmlStatusCode: http.StatusOK,
		},
		"TheLinuxEXP": {
			args: []string{
				`--exclude="squarespace.com,www.youtube.com,www.patreon.com,paypal.me,liberapay.com,the-linux-experiment.creator-spring.com"`,
				"--save-exclusions",
			},
			channelHandle:  "TheLinuxEXP",
			html:           `"rssUrl":"https://www.youtube.com/feeds/videos.xml?channel_id=UC5UAwBUum7CPN5buc-_N1Fw"`,
			htmlStatusCode: http.StatusOK,
			json:           `{"channelHandle": "TheLinuxEXP", "domainCounts": {"9to5linux.com": 2}, "excludedDomains": ["www.tuxedocomputers.com"]}`,
			wantError:      false,
			xml: `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns:media="http://search.yahoo.com/mrss/" xmlns="http://www.w3.org/2005/Atom">
 <link rel="self" href="http://www.youtube.com/feeds/videos.xml?channel_id=UC5UAwBUum7CPN5buc-_N1Fw"/>
 <id>yt:channel:5UAwBUum7CPN5buc-_N1Fw</id>
 <yt:channelId>5UAwBUum7CPN5buc-_N1Fw</yt:channelId>
 <title>The Linux Experiment</title>
 <link rel="alternate" href="https://www.youtube.com/channel/UC5UAwBUum7CPN5buc-_N1Fw"/>
 <author>
  <name>The Linux Experiment</name>
  <uri>https://www.youtube.com/channel/UC5UAwBUum7CPN5buc-_N1Fw</uri>
 </author>
 <published>2018-02-21T14:20:56+00:00</published>
 <entry>
  <id>yt:video:upOaSBo5ygU</id>
  <yt:videoId>upOaSBo5ygU</yt:videoId>
  <yt:channelId>UC5UAwBUum7CPN5buc-_N1Fw</yt:channelId>
  <title>Ubuntu under attack, Big flaw affects all Linux distros, Linux beats Windows - Linux Weekly News</title>
  <link rel="alternate" href="https://www.youtube.com/watch?v=upOaSBo5ygU"/>
  <author>
   <name>The Linux Experiment</name>
   <uri>https://www.youtube.com/channel/UC5UAwBUum7CPN5buc-_N1Fw</uri>
  </author>
  <published>2026-05-02T09:46:06+00:00</published>
  <updated>2026-05-03T07:46:39+00:00</updated>
  <media:group>
   <media:title>Ubuntu under attack, Big flaw affects all Linux distros, Linux beats Windows - Linux Weekly News</media:title>
   <media:content url="https://www.youtube.com/v/upOaSBo5ygU?version=3" type="application/x-shockwave-flash" width="640" height="390"/>
   <media:thumbnail url="https://i2.ytimg.com/vi/upOaSBo5ygU/hqdefault.jpg" width="480" height="360"/>
   <media:description>Head to https://squarespace.com/thelinuxexperiment to save 10% off your first purchase of a website or domain using code THELINUXEXPERIMENT

Grab a brand new laptop or desktop running Linux: https://www.tuxedocomputers.com/en#


👏 SUPPORT THE CHANNEL:
Get access to:
- a Daily Linux News show
- a weekly patroncast for more thoughts
- your name in the credits

YouTube: https://www.youtube.com/@TheLinuxEXP/join
Patreon: https://www.patreon.com/thelinuxexperiment

Or, you can donate whatever you want:
https://paypal.me/thelinuxexp
Liberapay: https://liberapay.com/TheLinuxExperiment/

👕 GET TLE MERCH
Support the channel AND get cool new gear: https://the-linux-experiment.creator-spring.com/

Timestamps:

00:00 Intro
00:30 Sponsor: SquareSpace
02:09 Ubuntu's infrastructure is under attack
03:59 New nasty vulnerability for all Linux distros
05:36 Framework Pro's Linux version outsells the Windows one
06:50 Ubuntu LTS beats Windows 11 in benchmarks, again
08:23 Ubuntu will add AI features in the future
10:34 Ubuntu clarifies their inclusion of AI features
12:46 EU app for age verification is structurally flawed
15:26 Steam Controller announced at 99USD
17:43 Anthropic now funds Blender
21:00 Fedora 44 released
22:44 New campaign fighting back against Android lockdown
26:20 Sponsor: Tuxedo Computers

Links:

Ubuntu's infrastructure is under attack
https://thecybersecguru.com/news/massive-attack-ubuntu-canonical-313-team-extortion/

New nasty vulnerability for all Linux distros
https://nerds.xyz/2026/04/copy-fail-linux-root-exploit/
https://xint.io/blog/copy-fail-linux-distributions

Framework Pro's Linux version outsells the Windows one
https://www.pcworld.com/article/3123900/framework-new-linux-laptop-is-selling-faster-than-its-windows-one.html

Ubuntu LTS beats Windows 11 in benchmarks, again
https://www.phoronix.com/review/ubuntu-2604-windows-11/7

Ubuntu will add AI features in the future
https://discourse.ubuntu.com/t/the-future-of-ai-in-ubuntu/81130

Ubuntu clarifies their inclusion of AI features
https://discourse.ubuntu.com/t/the-future-of-ai-in-ubuntu/81130/41

EU app for age verification is structurally flawed
https://www.techradar.com/vpn/vpn-privacy-security/the-eus-age-verification-app-has-a-privacy-problem-and-it-may-be-more-than-just-a-bug-in-an-app
https://cybernews.com/security/eu-age-verification-app-hack/

Steam Controller announced at 99USD
https://store.steampowered.com/sale/steamcontroller
https://www.techradar.com/computing/peripherals-accessories/valve-steam-controller-2026
https://www.gamingonlinux.com/2026/04/valve-have-plans-for-the-steam-deck-2-plus-a-brief-steam-machine-steam-frame-update/

Anthropic now funds Blender
https://www.blender.org/press/anthropic-joins-the-blender-development-fund-as-corporate-patron/
https://www.anthropic.com/news/claude-for-creative-work
https://www.gamingonlinux.com/2026/05/blender-change-the-anthropic-ai-funding-deal-with-discussions-planned-for-ai-policies/

Fedora 44 released
https://fedoramagazine.org/announcing-fedora-linux-44/
https://fedoramagazine.org/whats-new-in-fedora-kde-plasma-desktop-44/
https://fedoramagazine.org/fedora-asahi-remix-44-is-now-available/

New campaign fighting back against Android lockdown
https://keepandroidopen.org/en/

#linuxdesktop #linuxdistro #linuxnews</media:description>
   <media:community>
    <media:starRating count="4122" average="5.00" min="1" max="5"/>
    <media:statistics views="75265"/>
   </media:community>
  </media:group>
 </entry>
</feed>`,
			xmlStatusCode: http.StatusOK,
		},
		"Bad HTML Status Code": {
			channelHandle:  "TheLinuxExp",
			htmlStatusCode: http.StatusNotFound,
			json:           "{}",
			wantError:      true,
		},
		"Bad XML Status Code": {
			channelHandle:  "TheLinuxEXP",
			html:           `"rssUrl":"https://www.youtube.com/feeds/videos.xml?channel_id=UC5UAwBUum7CPN5buc-_N1Fw"`,
			htmlStatusCode: http.StatusOK,
			json:           "{}",
			wantError:      true,
			xmlStatusCode:  http.StatusNotFound,
		},
		"HTML Timeout": {
			channelHandle: "TheLinuxEXP",
			htmlError:     http.ErrHandlerTimeout,
			json:          "{}",
			wantError:     true,
		},
		"XML Timeout": {
			channelHandle: "TheLinuxEXP",
			htmlError:     http.ErrHandlerTimeout,
			json:          "{}",
			wantError:     true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tubenoisseur := tubenoisseur{
				HTMLClient: &mockGetter{
					Error: test.htmlError,
					Response: &http.Response{
						Body:       io.NopCloser(bytes.NewBufferString(test.html)),
						StatusCode: test.htmlStatusCode,
					},
				},
				Storage: &mockStorage{json: test.json},
				XMLClient: &mockGetter{
					Error: test.xmlError,
					Response: &http.Response{
						Body:       io.NopCloser(bytes.NewBufferString(test.xml)),
						StatusCode: test.xmlStatusCode,
					},
				},
			}
			rootCommand := newRootCommand(tubenoisseur)
			args := append([]string{test.channelHandle}, test.args...)

			rootCommand.SetArgs(args)

			if err := rootCommand.Execute(); (err != nil) != test.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
