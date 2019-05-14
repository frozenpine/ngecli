// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/frozenpine/viper"

	"github.com/spf13/cobra"
)

// tradeCmd represents the trade command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "ngecli default config",
	Long:  `show ngecli default configs & save to config.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonBytes, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")

		fmt.Println(string(jsonBytes))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
