// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Package cmd implements the command-line interface (CLI) of the application.

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gbr-ufs/tubenoisseur/cli/state"
	"github.com/gbr-ufs/tubenoisseur/scraping"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

var k = koanf.New(".")

func newRootCommand(tubenoisseur tubenoisseur) *cobra.Command {
	rootCommand := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Long: `tubenoisseur is a command-line interface (CLI) meant to be coupled
which a scheduling system, such as systemd or a cron job, to gather links in
the description of videos from a YouTube channel, with the goal of mapping out
their sources.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := k.Load(confmap.Provider(map[string]any{
				"debug":           false,
				"exclude":         []string{},
				"save-exclusions": true,
			}, "."), nil); err != nil {
				return err
			}

			configPath := getConfigPath(cmd)

			if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
				log.Debugf("no configuration file loaded from %s", configPath)
			} else {
				log.Debugf("loaded configuration file from %s", configPath)
			}

			if err := k.Load(env.Provider("TUBENOISSEUR_", ".", func(s string) string {
				return strings.ReplaceAll(
					strings.ToLower(strings.TrimPrefix(s, "TUBENOISSEUR_")),
					"_",
					"-",
				)
			}), nil); err != nil {
				return err
			}

			if err := k.Load(posflag.Provider(cmd.Flags(), ".", k), nil); err != nil {
				return err
			}

			if k.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			channelHandle := args[0]
			exclude := k.Strings("exclude")
			saveExclusions := k.Bool("save-exclusions")

			log.Debug("attempting to load channel state")
			channelState, err := tubenoisseur.Storage.LoadState(channelHandle)

			if err != nil {
				return err
			}

			log.Debug("ensuring feed URL")
			if channelState.FeedURL == "" {
				response, err := tubenoisseur.HTMLClient.Get(
					"https://www.youtube.com/@" + channelHandle,
				)

				if err != nil {
					return err
				}

				if response.StatusCode != http.StatusOK {
					response.Body.Close()

					return fmt.Errorf("unexpected status code: %d", response.StatusCode)
				}

				channelState.FeedURL, err = scraping.ExtractFeed(response.Body)

				if err != nil {
					return err
				}

				defer response.Body.Close()
			}

			log.Debug("fetching feed XML")
			response, err := tubenoisseur.XMLClient.Get(channelState.FeedURL)

			if err != nil {
				return err
			}

			if response.StatusCode != http.StatusOK {
				response.Body.Close()

				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			log.Debug("extracting videos from feed")
			videos, err := scraping.ExtractVideos(response.Body)

			if err != nil {
				return err
			}

			videos = state.GetVideos(videos, channelState.LastVideoID)

			exclusionMap := getExclusionMap(channelState.ExcludedDomains, exclude)

			log.Debug("getting sources from videos")
			videoURLs := getVideoURLs(videos, exclusionMap)

			log.Debug("updating channel state")
			channelState = updateChannelState(
				channelState,
				channelHandle,
				exclude,
				saveExclusions,
				channelState.FeedURL,
				videos,
				videoURLs,
			)

			log.Debug("writing report")
			if err := tubenoisseur.Storage.WriteReport(
				channelHandle,
				channelState,
			); err != nil {
				return err
			}

			log.Debug("saving channel state")
			if err := tubenoisseur.Storage.SaveState(channelHandle, channelState); err != nil {
				return err
			}

			return nil
		},
		Short: "YouTuber Reference Scraper",
		Use:   "tubenoisseur [channel-handle]",
	}

	rootCommand.Flags().String("config-file", "", "path to the configuration file")
	rootCommand.Flags().Bool("debug", false, "whether to debug level messages")
	rootCommand.Flags().
		StringSliceP("exclude", "e", nil, "domains (without schema) to be excluded from the analysis")
	rootCommand.Flags().Bool("save-exclusions", true, "whether to save excluded domains")

	return rootCommand
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	rootCommand := newRootCommand(
		tubenoisseur{HTMLClient: httpClient, Storage: diskStorage{}, XMLClient: httpClient},
	)

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
