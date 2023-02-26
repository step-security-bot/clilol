// SPDX-License-Identifier: MPL-2.0
//
// Copyright (c) 2023 Mark Cornick
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getWeblogLatestOutput struct {
	Request  resultRequest `json:"request"`
	Response struct {
		Message string `json:"message"`
		Post    struct {
			Address  string `json:"address"`
			Location string `json:"location"`
			Title    string `json:"title"`
			Date     int64  `json:"date"`
			Type     string `json:"type"`
			Status   string `json:"status"`
			Source   string `json:"source"`
			Body     string `json:"body"`
			Output   string `json:"output"`
			Metadata string `json:"metadata"`
			Entry    string `json:"entry"`
		} `json:"post"`
	} `json:"response"`
}

var getWeblogLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Get the latest weblog entry",
	Long:  "Gets your weblog's latest entry",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getWeblogLatest()
		cobra.CheckErr(err)
		if result.Request.Success {
			fmt.Printf(
				"%s (%s) modified on %s\n\n%s\n",
				result.Response.Post.Entry,
				fmt.Sprintf(
					"https://%s.weblog.lol%s",
					result.Response.Post.Address,
					result.Response.Post.Location,
				),
				time.Unix(result.Response.Post.Date, 0),
				result.Response.Post.Body,
			)
		} else {
			cobra.CheckErr(fmt.Errorf(result.Response.Message))
		}
	},
}

func init() {
	getWeblogCmd.AddCommand(getWeblogLatestCmd)
}

func getWeblogLatest() (getWeblogLatestOutput, error) {
	err := checkConfig("address")
	cobra.CheckErr(err)
	var result getWeblogLatestOutput
	body := callAPIWithParams(
		http.MethodGet,
		"/address/"+viper.GetString("address")+"/weblog/post/latest",
		nil,
		true,
	)
	err = json.Unmarshal(body, &result)
	return result, err
}
