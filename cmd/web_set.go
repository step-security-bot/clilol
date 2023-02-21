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
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	webSetFilename string
	webSetPublish  bool
	webSetCmd      = &cobra.Command{
		Use:   "set",
		Short: "set webpage content",
		Long: `Sets your webpage content.

If you specify a filename with the --filename flag, the content of the file
will be used. If you do not specify a filename, the content will be read
from stdin.

Set the --publish flag to true publish your webpage. By default, it will not
be published.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			type Input struct {
				Publish bool   `json:"publish,omitempty"`
				Content string `json:"content"`
			}
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
			var content string
			if webSetFilename != "" {
				input, err := os.ReadFile(webSetFilename)
				cobra.CheckErr(err)
				content = string(input)
			} else {
				stdin, err := io.ReadAll(os.Stdin)
				cobra.CheckErr(err)
				content = string(stdin)
			}
			webPage := Input{webSetPublish, content}
			body := callAPI(
				http.MethodPost,
				"/address/"+viper.GetString("address")+"/web",
				webPage,
				true,
			)
			err := json.Unmarshal(body, &result)
			cobra.CheckErr(err)
			if !silent {
				if !wantJson {
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
	webSetCmd.Flags().StringVarP(
		&webSetFilename,
		"filename",
		"f",
		"",
		"file to read webpage from (default stdin)",
	)
	webSetCmd.Flags().BoolVarP(
		&webSetPublish,
		"publish",
		"p",
		false,
		"publish the updated page (default false)",
	)

	webCmd.AddCommand(webSetCmd)
}
