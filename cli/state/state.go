// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Package state maintains the state of the program.

package state

import (
	"encoding/json"
	"io"

	"github.com/gbr-ufs/tubenoisseur/scraping"
)

// LoadState loads the external JSON state linked to a channel into the
// application.
//
// Can fail in case of malformed JSON.
func LoadState(reader io.Reader) (ChannelState, error) {
	state := ChannelState{}
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&state)

	if err != nil {
		return state, err
	}

	return state, nil
}

// GetVideos gets all the videos that will be processed by the program.
func GetVideos(feed []scraping.VideoEntry, lastVideoID string) []scraping.VideoEntry {
	// Empty feed.
	if len(feed) == 0 {
		return feed
	}

	// Up to date.
	if lastVideoID == feed[0].ID {
		return []scraping.VideoEntry{}
	}

	// Catch up.
	for index, entry := range feed {
		if entry.ID == lastVideoID {
			return feed[0:index]
		}
	}

	// No videos previously processed, or the last processed video is now
	// unlisted.
	return feed
}
