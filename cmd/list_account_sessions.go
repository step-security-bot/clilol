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

var listAccountSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "list your sessions",
	Long:  `Lists the active sessions on your account.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		type Result struct {
			Request struct {
				StatusCode int  `json:"status_code"`
				Success    bool `json:"success"`
			} `json:"request"`
			Response []struct {
				SessionID string `json:"session_id"`
				UserAgent string `json:"user_agent"`
				CreatedIP string `json:"created_ip"`
				CreatedOn int64  `json:"created_on"`
				ExpiresOn int64  `json:"expires_on"`
			} `json:"response"`
		}
		var result Result
		body := callAPIWithParams(
			http.MethodGet,
			"/account/"+viper.GetString("email")+"/sessions",
			nil,
			true,
		)
		err := json.Unmarshal(body, &result)
		cobra.CheckErr(err)
		if !silentFlag {
			if !jsonFlag {
				if result.Request.Success {
					for _, session := range result.Response {
						fmt.Printf("\n%s\n", session.SessionID)
						fmt.Println(session.UserAgent)
						fmt.Println(session.CreatedIP)
						fmt.Println(time.Unix(session.CreatedOn, 0))
					}
				} else {
					cobra.CheckErr(fmt.Errorf("%d", result.Request.StatusCode))
				}
			} else {
				fmt.Println(string(body))
			}
		}
	},
}

func init() {
	listAccountCmd.AddCommand(listAccountSessionsCmd)
}