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
	"context"
	"fmt"
	"os"

	"github.com/frozenpine/ngerest"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultHost    = "http://trade"
	defaultBaseURI = "/api/v1"

	defaultSymbol = "XBTUSD"
)

var cfgFile string

var (
	client            *ngerest.APIClient
	rootCtx, stopFunc = context.WithCancel(context.Background())

	host, baseURI      string
	identity, password string
	symbol             string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ngecli",
	Short: "CLI for NGE with optional shell interface.",
	Long: `A CLI tool for interactive with NGE trade engine.
Supported:
	1. order
	2. trade
	3. execution
	4. position
	5. all websocket interface`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ngecli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", defaultHost, "Host address for NGE.")
	rootCmd.PersistentFlags().StringVarP(&baseURI, "uri", "R", defaultBaseURI, "Base URI for NGE.")

	rootCmd.PersistentFlags().StringVarP(&identity, "id", "u", "", "Identity used for login.")
	rootCmd.PersistentFlags().StringVarP(&password, "pass", "p", "", "Password used for login.")

	rootCmd.PersistentFlags().StringVar(&symbol, "symbol", defaultSymbol, "Symbol name.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ngecli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ngecli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
