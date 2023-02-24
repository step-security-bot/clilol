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

	"github.com/spf13/cobra"
)

type getServiceOutput struct {
	Request  resultRequest `json:"request"`
	Response struct {
		Message   string `json:"message"`
		Members   int    `json:"members"`
		Addresses int    `json:"addresses"`
		Profiles  int    `json:"profiles"`
	} `json:"response"`
}

var getServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Get service stats",
	Long:  "Gets statistics for omg.lol services.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getService()
		cobra.CheckErr(err)
		if result.Request.Success {
			fmt.Println(result.Response.Message)
		} else {
			cobra.CheckErr(fmt.Errorf(result.Response.Message))
		}
	},
}

func getService() (getServiceOutput, error) {
	var result getServiceOutput
	body := callAPIWithParams(http.MethodGet, "/service/info", nil, false)
	err := json.Unmarshal(body, &result)
	return result, err
}

func init() {
	getCmd.AddCommand(getServiceCmd)
}
