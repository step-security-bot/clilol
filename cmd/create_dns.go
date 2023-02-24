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

var (
	createDNSPriority int
	createDNSTTL      int
	createDNSCmd      = &cobra.Command{
		Use:   "dns [name] [type] [data]",
		Short: "Create a DNS record",
		Long:  "Creates a DNS record.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			type input struct {
				Type     string `json:"type"`
				Name     string `json:"name"`
				Data     string `json:"data"`
				Priority int    `json:"priority"`
				TTL      int    `json:"ttl"`
			}
			type output struct {
				Request  resultRequest `json:"request"`
				Response struct {
					Message  string `json:"message"`
					DataSent struct {
						Type     string `json:"type"`
						Priority int    `json:"priority"`
						TTL      int    `json:"ttl"`
						Name     string `json:"name"`
						Content  string `json:"content"`
					} `json:"data_sent"`
					ResponseReceived struct {
						Data struct {
							ID        int       `json:"id"`
							Name      string    `json:"name"`
							Content   string    `json:"content"`
							TTL       int       `json:"ttl"`
							Priority  int       `json:"priority"`
							Type      string    `json:"type"`
							CreatedAt time.Time `json:"created_at"`
							UpdatedAt time.Time `json:"updated_at"`
						} `json:"data"`
					} `json:"response_received"`
				} `json:"response"`
			}
			var result output
			dns := input{args[2], args[1], args[3], createDNSPriority, createDNSTTL}
			body := callAPIWithParams(
				http.MethodPost,
				"/address/"+viper.GetString("address")+"/dns",
				dns,
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
	createDNSCmd.Flags().IntVarP(
		&createDNSPriority,
		"priority",
		"p",
		0,
		"priority of the DNS record",
	)
	createDNSCmd.Flags().IntVarP(
		&createDNSTTL,
		"ttl",
		"T",
		3600,
		"time to live of the DNS record",
	)
	createCmd.AddCommand(createDNSCmd)
}
