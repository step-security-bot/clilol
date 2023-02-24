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

var updateAccountNameCmd = &cobra.Command{
	Use:   "name [name]",
	Short: "set the name on your account",
	Long:  "Sets the name on your account.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		type Input struct {
			Name string `json:"name"`
		}
		type Result struct {
			Request  responseRequest `json:"request"`
			Response struct {
				Message string `json:"message"`
				Name    string `json:"name"`
			} `json:"response"`
		}
		var result Result
		account := Input{args[0]}
		body := callAPIWithParams(
			http.MethodPost,
			"/account/"+viper.GetString("email")+"/name",
			account,
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
	updateAccountCmd.AddCommand(updateAccountNameCmd)
}
