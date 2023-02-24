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
	"github.com/spf13/viper"
)

var getAccountNameCmd = &cobra.Command{
	Use:   "name",
	Short: "Get your account name",
	Long:  "Gets the name on your account.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		type output struct {
			Request  resultRequest `json:"request"`
			Response struct {
				Message string `json:"message"`
				Name    string `json:"name"`
			} `json:"response"`
		}
		var result output
		body := callAPIWithParams(
			http.MethodGet,
			"/account/"+viper.GetString("email")+"/name",
			nil,
			true,
		)
		err := json.Unmarshal(body, &result)
		cobra.CheckErr(err)
		if result.Request.Success {
			fmt.Println(result.Response.Message)
		} else {
			cobra.CheckErr(fmt.Errorf(result.Response.Message))
		}
	},
}

func init() {
	getAccountCmd.AddCommand(getAccountNameCmd)
}
