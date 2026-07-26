// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Package scraping implements routines for HTML scraping, XML parsing and URL
// extraction for working with YouTube RSS feeds.

package scraping

import (
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"regexp"

	"mvdan.cc/xurls/v2"
)

var rssURL = regexp.MustCompile(`"rssUrl":"([^"]+)"`)
var rxStrict = xurls.Strict()

// ExtractFeed searches the HTML of a YouTube channel for the channel's RSS feed.
func ExtractFeed(channelHTML io.Reader) (string, error) {
	buffer, err := io.ReadAll(channelHTML)

	if err != nil {
		return "", err
	}

	subMatch := rssURL.FindSubmatch(buffer)

	if len(subMatch) < 2 {
		return "", errors.New("could not find RSS URL in channel HTML")
	}

	return string(subMatch[1]), nil
}

// ExtractURLs gets all the URLs with schema from text, ignoring those present
// in the excludedDomains map.
func ExtractURLs(text string, excludedDomains map[string]bool) []string {
	urls := rxStrict.FindAllString(text, -1)
	result := []string{}

	for _, item := range urls {
		parsed, err := url.Parse(item)

		if err != nil {
			continue
		}

		if !excludedDomains[parsed.Hostname()] {
			result = append(result, item)
		}
	}

	return result
}

// ExtractVideos parses the XML of a YouTube RSS feed from an [io.Reader],
// returning a slice of video entries.
func ExtractVideos(reader io.Reader) ([]VideoEntry, error) {
	feed := youTubeFeed{}
	decoder := xml.NewDecoder(reader)

	if err := decoder.Decode(&feed); err != nil {
		return []VideoEntry{}, err
	}

	return feed.Entries, nil
}
