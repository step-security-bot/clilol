// SPDX-License-Identifier: MPL-2.0
//
// Copyright (c) 2023 Mark Cornick
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const endpoint = "https://api.omg.lol"

var (
	addressFlag string
	idFlag      string
	jsonFlag    bool
	nameFlag    string
	silentFlag  bool
	version     = "dev"
	rootCmd     = &cobra.Command{
		Version: version,
		Use:     "clilol",
		Short:   "a cli for omg.lol",
		Long: `This is the root command. It does nothing on its own.
See the subcommands for more information.`,
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/clilol/")
	viper.AddConfigPath(configDir + "/clilol")
	viper.SetEnvPrefix("clilol")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			cobra.CheckErr(err)
		}
	}
	rootCmd.PersistentFlags().BoolVarP(
		&silentFlag,
		"silent",
		"s",
		false,
		"be silent",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&jsonFlag,
		"json",
		"j",
		false,
		"output json",
	)
	cobra.CheckErr(
		viper.BindPFlag(
			"silent",
			rootCmd.PersistentFlags().Lookup("silent"),
		),
	)
	cobra.CheckErr(
		viper.BindPFlag(
			"json",
			rootCmd.PersistentFlags().Lookup("json"),
		),
	)
	rootCmd.DisableAutoGenTag = true
}

func callAPI(method string, path string, bodyReader io.Reader, auth bool) ([]byte, error) {
	request, err := http.NewRequest(method, endpoint+path, bodyReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "clilol (https://github.com/mcornick/clilol)")
	request.Header.Set("Content-Type", "application/json")
	if auth {
		request.Header.Set("Authorization", "Bearer "+viper.GetString("apikey"))
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func callAPIWithJSON(method string, path string, params interface{}, auth bool) []byte {
	jsonBody, err := json.Marshal(params)
	cobra.CheckErr(err)
	bodyReader := bytes.NewReader(jsonBody)
	body, err := callAPI(method, path, bodyReader, auth)
	cobra.CheckErr(err)
	return body
}

func callAPIWithRawData(method string, path string, data string, auth bool) []byte {
	bodyReader := strings.NewReader(data)
	body, err := callAPI(method, path, bodyReader, auth)
	cobra.CheckErr(err)
	return body
}
