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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frozenpine/ngecli/models"
	"github.com/frozenpine/pkcs8"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	auths *viper.Viper
)

// Login login with identity & password to get auth Context
func Login(ctx context.Context, identity, password string) context.Context {
	idMap := models.NewIdentityMap()
	login := make(map[string]string)

	if err := idMap.CheckIdentity(identity, login); err != nil {
		fmt.Println(err)
		return nil
	}

	pubKey, _, err := client.KeyExchange.GetPublicKey(ctx)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	login["password"] = pubKey.Encrypt(password)

	auth, _, err := client.User.UserLogin(ctx, login)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return auth
}

// GetUserDefaultKey get user's default sys api key
func GetUserDefaultKey(auth context.Context) *models.APIKey {
	priKey := pkcs8.GeneratePriveKey(2048)

	userDefault, _, err := client.User.UserGetDefaultAPIKey(auth, priKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	key := models.APIKey{
		Key:    userDefault.APIKey,
		Secret: userDefault.APISecret,
	}

	return &key
}

func parseHost(hostString string) bool {
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

	initClient()

	return true
}

func loginAndSave(host string) {
	identity := ReadLine("Identity: ", nil)
	password := ReadLine("Password: ", nil)

	if auth := Login(rootCtx, identity, password); auth == nil {
		fmt.Println("Login failed.")
		os.Exit(1)
	}

	auths.Set(host+".identity", identity)
	auths.Set(host+".password", password)
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login NGE trade engin with user identity.",
	Long:  `Login NGE trade engin and save identity info to config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, host := range args {
				if !parseHost(host) {
					continue
				}
				
				loginAndSave(host)
			}
		} else {
			initClient()

			loginAndSave(getBaseHost())
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
		auths = viper.New()
	}

	confDIR := filepath.Join(home, ".ngecli")
	if _, err := os.Stat(confDIR); os.IsNotExist(err) {
		os.Mkdir(confDIR, os.ModePerm)
	}

	auths.SetConfigFile(filepath.Join(confDIR, "auths.yaml"))
}
