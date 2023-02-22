// SPDX-License-Identifier: MPL-2.0
//
// Copyright (c) 2023 Mark Cornick
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	updateWebPFPFilename string
	updateWebPFPCmd      = &cobra.Command{
		Use:   "pfp",
		Short: "set your profile picture",
		Long: `Sets your profile picture.

Specify an image file with the --filename flag.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			type Result struct {
				Request struct {
					StatusCode int  `json:"status_code"`
					Success    bool `json:"success"`
				} `json:"request"`
				Response struct {
					Message string `json:"message"`
				} `json:"response"`
			}
			var result Result
			content, err := os.ReadFile(updateWebPFPFilename)
			cobra.CheckErr(err)
			encoded := "data:text/plain;base64," + base64.StdEncoding.EncodeToString(content)
			body := callAPIWithRawData(
				http.MethodPost,
				"/address/"+viper.GetString("address")+"/pfp",
				encoded,
				true,
			)
			err = json.Unmarshal(body, &result)
			cobra.CheckErr(err)
			if !silentFlag {
				if !jsonFlag {
					if result.Request.Success {
						fmt.Println(result.Response.Message)
					} else {
						cobra.CheckErr(fmt.Errorf(result.Response.Message))
					}
				} else {
					fmt.Println(string(body))
				}
			}
		},
	}
)

func init() {
	updateWebPFPCmd.Flags().StringVarP(
		&updateWebPFPFilename,
		"filename",
		"f",
		"",
		"file to read PFP from (default stdin)",
	)
	updateWebCmd.AddCommand(updateWebPFPCmd)
}