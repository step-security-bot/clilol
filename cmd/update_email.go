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

type updateEmailInput struct {
	Destination string `json:"destination"`
}
type updateEmailOutput struct {
	Request  resultRequest `json:"request"`
	Response struct {
		Message           string   `json:"message"`
		DestinationString string   `json:"destination_string"`
		DestinationArray  []string `json:"destination_array"`
		Address           string   `json:"address"`
		EmailAddress      string   `json:"email_address"`
	} `json:"response"`
}

var updateEmailCmd = &cobra.Command{
	Use:   "email [address]",
	Short: "set email forwarding address(es)",
	Long: `Sets your email forwarding address(es).
	
To specify multiple addresses, separate them with commas.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var result updateEmailOutput
		email := updateEmailInput{args[0]}
		result, err := updateEmail(email)
		cobra.CheckErr(err)
		if result.Request.Success {
			fmt.Println(result.Response.Message)
		} else {
			cobra.CheckErr(fmt.Errorf(result.Response.Message))
		}
	},
}

func init() {
	updateCmd.AddCommand(updateEmailCmd)
}

func updateEmail(params updateEmailInput) (updateEmailOutput, error) {
	err := checkConfig("address")
	cobra.CheckErr(err)
	var result updateEmailOutput
	body := callAPIWithParams(
		http.MethodPost,
		"/address/"+viper.GetString("address")+"/email",
		params,
		true,
	)
	err = json.Unmarshal(body, &result)
	return result, err
}
