// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package scraping

// VideoEntry defines a video entry in an RSS feed.
type VideoEntry struct {
	// Description defines the video description.
	Description string `xml:"group>description"`
	// ID defines the video ID.
	ID string `xml:"videoId"`
	// Published defines the video publication date.
	Published string `xml:"published"`
	// Title defines the video title.
	Title string `xml:"title"`
}

type youTubeFeed struct {
	Entries []VideoEntry `xml:"entry"`
}
