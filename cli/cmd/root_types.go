// SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gbr-ufs/tubenoisseur/cli/state"
	"github.com/gbr-ufs/tubenoisseur/cli/template"
)

type httpGetter interface {
	Get(url string) (*http.Response, error)
}

type storage interface {
	LoadState(channelHandle string) (state.ChannelState, error)
	SaveState(channelHandle string, channelState state.ChannelState) error
	WriteReport(channelHandle string, channelState state.ChannelState) error
}

type tubenoisseur struct {
	HTMLClient httpGetter
	Storage    storage
	XMLClient  httpGetter
}

type diskStorage struct{}

func (d diskStorage) LoadState(channelHandle string) (state.ChannelState, error) {
	channelDirPath := filepath.Join("tubenoisseur", channelHandle)

	if err := os.MkdirAll(channelDirPath, 0755); err != nil {
		return state.ChannelState{}, err
	}

	channelJSONPath := filepath.Join(channelDirPath, channelHandle+".json")
	file, err := os.Open(channelJSONPath)

	if errors.Is(err, os.ErrNotExist) {
		return state.ChannelState{ChannelHandle: channelHandle}, nil
	} else if err != nil {
		return state.ChannelState{}, err
	}

	defer file.Close()

	return state.LoadState(file)
}

func (d diskStorage) SaveState(channelHandle string, channelState state.ChannelState) error {
	channelJSONPath := filepath.Join("tubenoisseur", channelHandle, channelHandle+".json")
	updatedJSON, err := json.MarshalIndent(channelState, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(channelJSONPath, updatedJSON, 0644)
}

func (d diskStorage) WriteReport(channelHandle string, channelState state.ChannelState) error {
	channelMarkdownPath := filepath.Join("tubenoisseur", channelHandle, channelHandle+".md")
	file, err := os.OpenFile(channelMarkdownPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	return template.WriteReport(file, channelState)
}
