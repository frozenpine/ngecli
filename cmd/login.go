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

	"github.com/frozenpine/ngerest"

	"github.com/frozenpine/ngecli/models"
	"github.com/frozenpine/pkcs8"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// APIAuthCache api auth cache
type APIAuthCache struct {
	savedAuths   *viper.Viper
	apiKeyCache  map[string]*models.APIKey
	apiAuthCache map[string]context.Context
	authList     []*models.Authentication
	nextKeyIDX   int
}

// SetConfigFile set auth config file path
func (auth *APIAuthCache) SetConfigFile(path string) {
	auths.savedAuths.SetConfigFile(path)
}

// WriteConfig write login info to auth config file
func (auth *APIAuthCache) WriteConfig() error {
	return auths.savedAuths.WriteConfig()
}

// SaveLoginInfo save host's login info in viper config
func (auth *APIAuthCache) SaveLoginInfo(host, identity string, password *models.Password) {
	auths.savedAuths.Set(host+".identity", identity)
	auths.savedAuths.Set(host+".password", password.String())
}

func (auth *APIAuthCache) retriveAuth() error {
	baseHost := getBaseHost()

	if identity == "" || password.IsSet() {
		if !auth.savedAuths.IsSet(baseHost) {
			return models.ErrAuthMissing
		}

		login := auth.savedAuths.Sub(baseHost)

		identity = login.GetString("identity")

		password.ShadowSet(login.GetString("password"))
	}

	if key := GetUserDefaultKey(rootCtx); key != nil {
		authInfo := models.Authentication{
			Identity: identity,
			Password: password,
			APIKey:   *key,
		}
		auth.authList = append(auth.authList, &authInfo)
	} else {
		return fmt.Errorf("retrive %s's api key from %s failed", identity, baseHost)
	}

	return nil
}

// NextAuth get next auth context
func (auth *APIAuthCache) NextAuth(parent context.Context) context.Context {
	if parent == nil {
		parent = rootCtx
	}

	defer func() {
		auth.nextKeyIDX++

		if auth.nextKeyIDX >= len(auth.authList) {
			auth.nextKeyIDX = 0
		}
	}()

	if len(auth.authList) < 1 {
		err := auth.retriveAuth()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	authInfo := auth.authList[auth.nextKeyIDX]

	var ctx context.Context

	if ctx, exist := auth.apiAuthCache[authInfo.Identity]; !exist {
		ctx = context.WithValue(
			parent, ngerest.ContextAPIKey, ngerest.APIKey{
				Key:    authInfo.Key,
				Secret: authInfo.Secret,
			})

		auth.apiAuthCache[authInfo.Identity] = ctx
	}

	if authInfo.Identity == "" {
		return ctx
	}

	if _, exist := auth.apiKeyCache[authInfo.Identity]; !exist {
		auth.apiKeyCache[authInfo.Identity] = &authInfo.APIKey
	}

	return ctx
}

// NewAPIAuthCache create new api auth cache
func NewAPIAuthCache() *APIAuthCache {
	cache := APIAuthCache{
		savedAuths: viper.New(),
	}

	return &cache
}

var (
	auths *APIAuthCache
)

// Login login with identity & password to get auth Context
func Login(ctx context.Context, identity string, password *models.Password) context.Context {
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

	login["password"] = pubKey.Encrypt(password.Show())

	auth, _, err := client.User.UserLogin(ctx, login)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return auth
}

// GetUserDefaultKey get user's default sys api key
func GetUserDefaultKey(loginAuth context.Context) *models.APIKey {
	if _, ok := loginAuth.Value(ngerest.ContextQuantToken).(ngerest.QuantToken); !ok {
		fmt.Println("invalid login auth")
		return nil
	}

	priKey := pkcs8.GeneratePriveKey(2048)

	userDefault, _, err := client.User.UserGetDefaultAPIKey(loginAuth, priKey)
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

	initClient()

	return true
}

func loginAndSave(host string) {
	identity := ReadLine("Identity: ", nil)
	password := models.NewPassword()
	password.Set(ReadLine("Password: ", nil))

	if auth := Login(rootCtx, identity, password); auth == nil {
		fmt.Println("Login failed.")
		os.Exit(1)
	}

	auths.SaveLoginInfo(host, identity, password)
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
		auths = NewAPIAuthCache()
	}

	confDIR := filepath.Join(home, ".ngecli")
	if _, err := os.Stat(confDIR); os.IsNotExist(err) {
		os.Mkdir(confDIR, os.ModePerm)
	}

	auths.SetConfigFile(filepath.Join(confDIR, "auths.yaml"))
}
