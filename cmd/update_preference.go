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

var (
	updatePreferenceItem  string
	updatePreferenceValue string
	updatePreferenceCmd   = &cobra.Command{
		Use:   "preference",
		Short: "set a preference",
		Long: `Sets omg.lol preferences.

Specify the preference item with the --item flag, and the value with
the --value flag.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			type Input struct {
				Item  string `json:"item"`
				Value string `json:"value"`
			}
			type Result struct {
				Request  responseRequest `json:"request"`
				Response struct {
					Message string `json:"message"`
					Item    string `json:"item"`
					Value   string `json:"value"`
				} `json:"response"`
			}
			var result Result
			pref := Input{updatePreferenceItem, updatePreferenceValue}
			body := callAPIWithParams(
				http.MethodPost,
				"/preferences/"+viper.GetString("address"),
				pref,
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
)

func init() {
	updatePreferenceCmd.Flags().StringVarP(
		&updatePreferenceItem,
		"item",
		"i",
		"",
		"preference item to set",
	)
	updatePreferenceCmd.Flags().StringVarP(
		&updatePreferenceValue,
		"value",
		"v",
		"",
		"value to set it to",
	)
	err := updatePreferenceCmd.MarkFlagRequired("item")
	cobra.CheckErr(err)
	err = updatePreferenceCmd.MarkFlagRequired("value")
	cobra.CheckErr(err)
	updateCmd.AddCommand(updatePreferenceCmd)
}
