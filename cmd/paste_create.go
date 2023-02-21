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
	pasteCreateTitle    string
	pasteCreateFilename string
	pasteCreateListed   bool
	pasteCreateCmd      = &cobra.Command{
		Use:   "create",
		Short: "create or update a paste",
		Long: `Create or update a paste in your pastebin.

Specify a title with the --title flag. If the title is already in use,
that paste will be updated. If the title is not in use, a new paste will
be created.

If you specify a filename with the --filename flag, the content of the file
will be used. If you do not specify a filename, the content will be read
from stdin.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			type Input struct {
				Title   string `json:"title"`
				Content string `json:"content"`
				Listed  int    `json:"listed,omitempty"`
			}
			type Result struct {
				Request struct {
					StatusCode int  `json:"status_code"`
					Success    bool `json:"success"`
				} `json:"request"`
				Response struct {
					Message string `json:"message"`
					Title   string `json:"title"`
				} `json:"response"`
			}
			var result Result
			var listed int
			var content string
			if pasteCreateFilename != "" {
				input, err := os.ReadFile(pasteCreateFilename)
				cobra.CheckErr(err)
				content = string(input)
			} else {
				stdin, err := io.ReadAll(os.Stdin)
				cobra.CheckErr(err)
				content = string(stdin)
			}
			if pasteCreateListed {
				listed = 1
			} else {
				listed = 0
			}
			paste := Input{pasteCreateTitle, content, listed}
			body := callAPI(
				http.MethodPost,
				"/address/"+viper.GetString("address")+"/pastebin",
				paste,
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
	pasteCreateCmd.Flags().StringVarP(
		&pasteCreateTitle,
		"title",
		"t",
		"",
		"title of the paste to create",
	)
	pasteCreateCmd.Flags().StringVarP(
		&pasteCreateFilename,
		"filename",
		"f",
		"",
		"file to read paste from (default stdin)",
	)
	pasteCreateCmd.Flags().BoolVarP(
		&pasteCreateListed,
		"listed",
		"l",
		false,
		"create paste as listed (default false)",
	)
	pasteCmd.AddCommand(pasteCreateCmd)
}
