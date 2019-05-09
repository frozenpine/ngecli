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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frozenpine/ngecli/models"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	auths *models.AuthCache
)

func parseArgHost(hostString string) bool {
	hosts := strings.Split(hostString, ":")

	host := hosts[0]
	if host != viper.GetString("host") {
		viper.Set("host", host)
	}

	if len(hosts) > 1 {
		port, err := strconv.Atoi(hosts[1])
		if err != nil {
			fmt.Println("Invalid host:", hostString)
			return false
		}

		if port != viper.GetInt("port") {
			viper.Set("port", port)
		}
	}

	auths.ChangeHost("")

	return true
}

func loginAndSave(host string) {
	identity := ReadLine("Identity: ", nil)
	password := models.NewPassword()
	password.Set(ReadLine("Password: ", nil))

	auths.ChangeHost(host)
	if auth := auths.Login(identity, password); auth == nil {
		fmt.Println("Login failed.")
		os.Exit(1)
	}

	auths.SetLoginInfo(host, identity, password)
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login NGE trade engine with user identity.",
	Long:  `Login NGE trade engine and save identity info to auths.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, host := range args {
				if !parseArgHost(host) {
					continue
				}

				loginAndSave(host)
			}
		} else {
			loginAndSave("")
		}

		err := auths.WriteConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	initAuthConfig()
}

func initAuthConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if auths == nil {
		auths = models.NewAuthCache(rootCtx, clientHub)
	}

	confDIR := filepath.Join(home, ".ngecli")
	if _, err := os.Stat(confDIR); os.IsNotExist(err) {
		os.Mkdir(confDIR, os.ModePerm)
	}

	auths.SetConfigFile(filepath.Join(confDIR, "auths.yaml"))
}
