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

// host string: host:port w/o scheme(http | https)
func parseArgHost(host string) bool {
	hostParts := strings.Split(host, ":")

	hostname := hostParts[0]
	if hostname != viper.GetString("host") {
		viper.Set("host", hostname)
	}

	if len(hostParts) > 1 {
		hostPort, err := strconv.Atoi(hostParts[1])
		if err != nil {
			fmt.Println("Invalid host:", host)
			return false
		}

		if hostPort != viper.GetInt("port") {
			viper.Set("port", hostPort)
		}
	}

	return true
}

// CollectLoginInfo from stdin
func CollectLoginInfo() (identity string, password *models.Password) {
	password = models.NewPassword()

	if debugLevel > 0 {
		identity = "sonny.frozenpine@gmail.com"
		password.Set("yuanyang")
	} else {
		identity = ReadLine("Identity: ", nil)
		password.Set(ReadLine("Password: ", nil))
	}

	return
}

// host string: host:port w/o scheme(http | https)
func loginAndSave(host string) {
	identity, password := CollectLoginInfo()

	fmt.Println("Try to login into:", models.GetBaseURL())

	if auth := auths.Login(identity, password); auth == nil {
		fmt.Println("Login failed.")
		os.Exit(1)
	} else {
		fmt.Println("Login success.")
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
				parts := strings.Split(host, "://")

				if strings.Contains(parts[0], "http") {
					viper.Set("scheme", parts[0])
				}

				hostString := parts[len(parts)-1]

				if !parseArgHost(hostString) {
					continue
				}

				loginAndSave(hostString)
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
