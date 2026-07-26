// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/gbr-ufs/tubenoisseur/cli/state"
	"github.com/gbr-ufs/tubenoisseur/scraping"
	"github.com/spf13/cobra"
)

func getConfigPath(cmd *cobra.Command) string {
	if cmd.Flags().Changed("config-file") {
		path, _ := cmd.Flags().GetString("config-file")

		return path
	}

	if path := os.Getenv("TUBENOISSEUR_CONFIG_FILE"); path != "" {
		return path
	}

	if _, err := os.Stat("tubenoisseur.toml"); err == nil {
		return "tubenoisseur.toml"
	}

	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "tubenoisseur", "config.toml")
	}

	return "tubenoisseur.toml"
}

func getDomainCounts(domainCounts map[string]uint, videoURLs map[string][]string) map[string]uint {
	for _, urls := range videoURLs {
		for _, rawURL := range urls {
			// [scraping.ExtractURLs] already checks for parsing
			// errors.
			parsed, _ := url.Parse(rawURL)
			domainCounts[parsed.Hostname()]++
		}
	}

	return domainCounts
}

func getVideoURLs(
	videos []scraping.VideoEntry,
	excludedDomains map[string]bool,
) map[string][]string {
	videoURLs := make(map[string][]string, len(videos))

	for _, video := range videos {
		extracted := scraping.ExtractURLs(video.Description, excludedDomains)

		videoURLs[video.ID] = extracted
	}

	return videoURLs
}

func updateChannelState(
	channelState state.ChannelState,
	channelHandle string,
	exclude []string,
	saveExclusions bool,
	feedURL string,
	videos []scraping.VideoEntry,
	videoURLs map[string][]string,
) state.ChannelState {
	channelState.ChannelHandle = channelHandle

	if channelState.DomainCounts == nil {
		channelState.DomainCounts = map[string]uint{}
	}

	channelState.DomainCounts = getDomainCounts(channelState.DomainCounts, videoURLs)

	if saveExclusions {
		for _, domain := range exclude {
			if !slices.Contains(channelState.ExcludedDomains, domain) {
				channelState.ExcludedDomains = append(channelState.ExcludedDomains, domain)
			}
		}
	}

	channelState.FeedURL = feedURL

	if len(videos) > 0 {
		channelState.LastVideoID = videos[0].ID

		processedVideos := make([]state.ProcessedVideo, 0, len(videos))

		for _, video := range videos {
			processedVideos = append(processedVideos, state.ProcessedVideo{
				ID:        video.ID,
				Published: video.Published,
				Title:     video.Title,
				URLs:      videoURLs[video.ID],
			})
		}

		channelState.Videos = append(processedVideos, channelState.Videos...)
	}

	return channelState
}

func getExclusionMap(channelStateExcludedDomains []string, exclude []string) map[string]bool {
	exclusionMap := make(map[string]bool, len(channelStateExcludedDomains)+len(exclude))

	for _, domain := range channelStateExcludedDomains {
		exclusionMap[domain] = true
	}

	for _, domain := range exclude {
		exclusionMap[domain] = true
	}

	return exclusionMap
}
