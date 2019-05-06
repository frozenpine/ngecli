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
	"path"
	"strconv"

	"github.com/frozenpine/ngecli/models"

	"github.com/frozenpine/ngerest"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultScheme  = "https"
	defaultHost    = "trade"
	defaultPort    = 80
	defaultBaseURI = "/api/v1"

	defaultSymbol = "XBTUSD"
)

var cfgFile string

var (
	client            *ngerest.APIClient
	rootCtx, stopFunc = context.WithCancel(context.Background())

	identity string
	password models.Password
	authFile string
	symbol   string
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("main run.")
	},
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

	viper.SetDefault("scheme", defaultScheme)
	rootCmd.PersistentFlags().String("scheme", defaultScheme, "Host scheme for NGE.")
	viper.BindPFlag("scheme", rootCmd.PersistentFlags().Lookup("scheme"))

	viper.SetDefault("host", defaultHost)
	rootCmd.PersistentFlags().StringP("host", "H", defaultHost, "Host address for NGE.")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))

	viper.SetDefault("port", defaultPort)
	rootCmd.PersistentFlags().IntP("port", "P", defaultPort, "Host port for NGE.")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	viper.SetDefault("base-uri", defaultBaseURI)
	rootCmd.PersistentFlags().String("uri", defaultBaseURI, "Base URI for NGE.")
	viper.BindPFlag("base-uri", rootCmd.PersistentFlags().Lookup("uri"))

	rootCmd.PersistentFlags().StringVarP(&identity, "id", "u", "", "Identity used for login.")
	rootCmd.PersistentFlags().VarP(&password, "pass", "p", "Password used for login.")

	rootCmd.PersistentFlags().StringVar(&authFile, "auth", "", "Auth info for NGE.")

	rootCmd.PersistentFlags().StringVar(&symbol, "symbol", defaultSymbol, "Symbol name.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".ngecli" (without extension).
		viper.AddConfigPath(path.Join(home, ".ngecli"))
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("ngecli")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
READ_CONFIG:
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found, creating one...")

		confDir := path.Join(home, ".ngecli")
		if _, err := os.Stat(confDir); os.IsNotExist(err) {
			os.Mkdir(confDir, os.ModePerm)
		}

		err = viper.WriteConfigAs(path.Join(confDir, "config.yaml"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		goto READ_CONFIG
	}
}

func initClient() {
	cfg := ngerest.NewConfiguration()
	client = ngerest.NewAPIClient(cfg)

	basePath := getBasePath()

	client.ChangeBasePath(basePath)

	fmt.Println("Change host to:", basePath)
}

func getBasePath() string {
	baseURI := viper.GetString("base-uri")

	return getBaseURL() + baseURI
}

func getBaseHost() string {
	port := viper.GetInt("port")
	host := viper.GetString("host")

	if port != defaultPort {
		return host + ":" + strconv.Itoa(port)
	}

	return host
}

func getBaseURL() string {
	scheme := viper.GetString("scheme")

	return scheme + "://" + getBaseHost()
}
