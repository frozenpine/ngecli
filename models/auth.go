package models

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/frozenpine/ngerest"
	"github.com/frozenpine/pkcs8"
	"github.com/gocarina/gocsv"
	"github.com/spf13/viper"
)

const (
	viperHostnameKeyDelim = "_"
)

// IdentityMap identity pattern map
type IdentityMap map[string]*regexp.Regexp

// CheckIdentity check & modify login map
func (idMap *IdentityMap) CheckIdentity(
	id string, login map[string]string) error {
	for name, pattern := range *idMap {
		if !pattern.MatchString(id) {
			continue
		}

		login[name] = id
		login["type"] = "account"
		login["verifyCode"] = ""

		return nil
	}

	return errors.New("identity should either be email or mobile")
}

// AddPattern add new pattern to IdentityMap
func (idMap *IdentityMap) AddPattern(
	name string, pattern *regexp.Regexp) error {
	if _, exist := (*idMap)[name]; exist {
		return fmt.Errorf("named[%s] pattern already exists", name)
	}

	(*idMap)[name] = pattern

	return nil
}

// NewIdentityMap generate identity pattern map
func NewIdentityMap() IdentityMap {
	mp := make(IdentityMap)

	mp.AddPattern("email", regexp.MustCompile(`[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)*`))
	mp.AddPattern("mobile", regexp.MustCompile(`(\+?[0-9]{2,3})?[0-9-]{6,13}`))

	return mp
}

// Password shadowed password store
type Password struct {
	shadowed string
	key      *rsa.PrivateKey
}

// IsSet verify if password set
func (p *Password) IsSet() bool {
	return p.shadowed != ""
}

func (p *Password) String() string {
	return string(p.shadowed)
}

// Set set password
func (p *Password) Set(value string) error {
	p.Shadow(value)

	return nil
}

// Type get password type
func (p *Password) Type() string {
	return "Password"
}

// ShadowSet set shadowed password
func (p *Password) ShadowSet(value string) (err error) {
	origin := p.shadowed
	p.shadowed = value

	defer func() {
		if recErr := recover(); recErr != nil {
			p.shadowed = origin
			err = errors.New("invalid shadowed password")
		}
	}()

	p.Show()

	return
}

// UnmarshalCSV unmarshal password from csv
func (p *Password) UnmarshalCSV(value string) error {
	if value == "" {
		return nil
	}

	return p.ShadowSet(value)
}

// MarshalCSV marshal password to csv
func (p *Password) MarshalCSV() string {
	return p.String()
}

// UnmarshalJSON unmarshal from json string
func (p *Password) UnmarshalJSON(data []byte) error {
	strValue := strings.Trim(string(data), "\"")

	if strValue == "" {
		return nil
	}

	return p.ShadowSet(strValue)
}

// MarshalJSON marshal to json string
func (p *Password) MarshalJSON() ([]byte, error) {
	var buff bytes.Buffer
	buff.WriteString("\"" + (*p).String() + "\"")

	return buff.Bytes(), nil
}

// Shadow shadow & store a password
func (p *Password) Shadow(value string) string {
	if p.key == nil {
		p.defaultKey(nil)
	}

	encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, &p.key.PublicKey, []byte(value))
	p.shadowed = base64.StdEncoding.EncodeToString(encrypted)
	return p.shadowed
}

// Show show unshadowed password
func (p *Password) Show() string {
	if p.key == nil {
		p.defaultKey(nil)
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(p.shadowed)
	if err != nil {
		panic(err.Error())
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, p.key, cipherBytes)
	if err != nil {
		panic(err.Error())
	}

	return string(decrypted)
}

func (p *Password) defaultKey(keyContent []byte) {
	if keyContent == nil || len(keyContent) < 1 {
		keyContent = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAw5QsrDwSbN5iAd4R2D7ARXo/4x5IGlvcbBx1jSnE8s2y9kn2
8ee/ujc+VWZ7I5SJDxV8VEa1AD73tpKOYVkz88D7mKzL4E6zGVTMRQnqGifUNr+l
KmKo2y13cOCL+hGGV31AJMnAygBKdSaY7ywvVZeiDUuYlb2COBY54EC2BbDwvgyo
o1dh1SQ9Yo7iTWI+nE0as9ugN/ljgBzk8UlPF0vSMBBdeagklzaOynJsTZi+oFzN
FDawArBDn2/vYKTjHsi0hRqAQsphqee9jZ+2P248FvPsPyjNCWeB4cEiaDBhW6C3
bCA3mxWh+RL8PThTAzScmtTPUFc+pGrV9h3twwIDAQABAoIBAE+lix/FAvflBGKg
TgITY/enVlcmoNRjLnu0h0aqiPMcQ9I6wt82soSiNLdQmbsepUZISK6FcPadrgFi
46rPSfHtWEiPriM1yYf6WYmQBQ8Lw0dcemWtcfh4JpkISNjYxC3i6vlQVDuvLKNS
yZs1Ej24F510eLoaR+qRWpZxo+7irnK35t6L6k6M6VIHJ5snXfcPowvKnMr++Hi+
RU+Gnqv1m8rVjpU5ADPJ0TB8giHJzl6/MS8ZbAMzcVGPRM+fa/S+Vdh1/2pNZWz5
0YgWIhATyywav7A2hd11ImUET5wdW/IJ+A1AM3UlaF8B7O9oZ+FhPYGfrGBbQoSo
o4mAk7ECgYEA6OYtjwnAwZ+7NUN4NCzhUk1ufAlWCIlGXovpe+gMN/dKRmBls4yk
nav52Ch4Y3PlZ0oPMvFrGWmDUlpHr+dv2MYYT4fRsLwcqedNakp/5syQ1tjlNQNv
HPynTPQnr0TKO1zB//d9fGJnKCZWkAgJGdxyfpRoPoUrYgOanaleVdsCgYEA1vpZ
62QiaTQuL63mGbKqj4fNP2nsWHwVsm/c6NAIZlbq2CY//0IR2nUVR/K3FhgT5HLm
pvUkK0SR7mlmTe36QVrTQ+Yo4VNuNkZOyxEYhN5/ZsJIk5SQaaeQrhK+CX90E0wJ
xirmWRBAsNuf+6l5ZSlRc32mh6SCdZUcEkG3cDkCgYEAlATWp7YXH/gYzz1WRDLR
8bDsq0BzwXEdnDFn7ywHt/oe5qOVf4u/g8YtQEhYWzzpa8AR8Nqmqrv4jnp6XT3G
RAuCn+k+SAkGXqV2+jrnFxSkaSfoZM0N7WpWGf6Cyk36CchmM/xjcI5J6aaUFW5F
+n209uXzaujQLbcEqXdfUUkCgYB8VAFY/2pfSYxEit/+kLPPmox7VjkX23t43PT3
uAiDl1TueQCeEYndu8T4/UghgP9QKZt3h2LJmziCl3ZRL4aB8ZMpO5z845Fj1jmP
e22gukUYGth6cXsrf3tPEQvS1mE9H8avUvQxIhMntXzKwPKyLLksf8ilveCtO/Um
IdeDEQKBgDomdSAKO26u69qfjTfdTtDI25VJ98YDQVIMAGNHTjDlRei7zGZcVkWS
sDf9BNxVpu0u2tOKf+oigcYnRPlUbZcFk8zUlPbjfz+r/bhS8PoMd9UNFTaO4z+L
Z4Njdti1yOD3gUoJ3DmqWRv0oS+L9iXag3p2GwzTG7El+LaoDUUS
-----END RSA PRIVATE KEY-----`)
	}

	block, rest := pem.Decode(keyContent)
	if block == nil || len(rest) > 0 {
		panic("pem decode key string failed.")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("parse pkcs8 private key failed: " + err.Error())
	}

	p.key = key
}

// NewPassword get new password struct
func NewPassword() *Password {
	pass := Password{}

	pass.defaultKey(nil)

	return &pass
}

// NewPasswordByPEM create password with user specified pem key,
// if pemPath not specified or is invalid,
// will return Password with default key
func NewPasswordByPEM(pemPath string) *Password {
	if pemPath == "" {
		fmt.Println("pem file path missing, using default key.")
		return NewPassword()
	}

	pemFile, err := os.OpenFile(pemPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("open pem file failed: %s, using default key.", err.Error())
		return NewPassword()
	}

	var keyContent []byte

	_, err = pemFile.Read(keyContent)
	if err != nil {
		fmt.Printf(
			"reading pem file failed: %s, using default key.", err.Error())
		return NewPassword()
	}

	pass := Password{}
	pass.defaultKey(keyContent)

	return &pass
}

// Authentication auth info including:
// 1. identity
// 2. Password instance
// 3. embedded APIKey struct
type Authentication struct {
	Identity string   `csv:"identity" json:"identity"`
	Password Password `csv:"password" json:"password"`
	APIKey
}

// Validate completion of auth info
func (auth *Authentication) Validate() bool {
	if auth.Password.IsSet() {
		if auth.Identity == "" {
			return false
		}
	} else {
		if !auth.APIKey.Validate() {
			return false
		}
	}

	return true
}

// APIKey key & secret for api
type APIKey struct {
	Key    string `csv:"api_key" json:"api_key"`
	Secret string `csv:"api_secret" json:"api_secret"`
}

// Validate completion of api key info
func (key *APIKey) Validate() bool {
	if key.Key == "" || key.Secret == "" {
		return false
	}

	// if len(key.Key) != 20 || len(key.Secret) != 99 {
	// 	return false
	// }

	return true
}

// AuthCache api auth cache
type AuthCache struct {
	savedAuths *viper.Viper

	keyCtxCache map[string]context.Context
	authList    []*Authentication

	clientHub   *ClientHub
	rootCtx     context.Context
	CmdAuthFile string
	DefaultID   string
	DefaultPass Password

	retriveOnece sync.Once
	keyIDX       uint32
}

func (cache *AuthCache) nextIDX() int {
	cache.retriveOnece.Do(func() {
		if len(cache.authList) >= 1 {
			return
		}

		var err error

		if cache.CmdAuthFile == "" {
			err = cache.retriveAuth()
		} else {
			err = cache.readAuthFile(cache.CmdAuthFile)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})

	// authList length can not longer
	idCount := atomic.AddUint32(&cache.keyIDX, 1) - 1

	idx := int(idCount) % len(cache.authList)

	return idx
}

// SetConfigFile set auth config file path
func (cache *AuthCache) SetConfigFile(path string) {
	cache.savedAuths.SetConfigFile(path)
	cache.savedAuths.ReadInConfig()
}

// WriteConfig write login info to auth config file
func (cache *AuthCache) WriteConfig() error {
	return cache.savedAuths.WriteConfig()
}

// SetLoginInfo save host's login info in viper config
func (cache *AuthCache) SetLoginInfo(
	host, identity string, password *Password) {
	cache.savedAuths.Set(
		strings.Join([]string{host, "identity"}, viperHostnameKeyDelim),
		identity)
	cache.savedAuths.Set(
		strings.Join([]string{host, "password"}, viperHostnameKeyDelim),
		password.String())
}

// Login login with identity & password to get auth Context
func (cache *AuthCache) Login(
	identity string, password *Password) context.Context {
	idMap := NewIdentityMap()
	loginInfo := make(map[string]string)

	var err error

	if err = idMap.CheckIdentity(identity, loginInfo); err != nil {
		fmt.Println(err)
		return nil
	}

	client, err := cache.clientHub.GetClient(GetBaseHost())
	if err != nil {
		panic(err)
	}

	pubKey, _, err := client.KeyExchange.GetPublicKey(cache.rootCtx)
	if err != nil {
		LogError(err)

		return nil
	}

	loginInfo["password"] = pubKey.Encrypt(password.Show())

	login, _, err := client.User.UserLogin(cache.rootCtx, loginInfo)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return login
}

// GetUserDefaultKey get user's default sys api key
func (cache *AuthCache) GetUserDefaultKey(loginAuth context.Context) *APIKey {
	if _, ok := loginAuth.Value(ngerest.ContextQuantToken).(ngerest.QuantToken); !ok {
		fmt.Println("invalid login auth")
		return nil
	}

	priKey := pkcs8.GeneratePriveKey(2048)

	client, err := cache.clientHub.GetClient(GetBaseHost())

	userDefault, _, err := client.User.UserGetDefaultAPIKey(
		loginAuth, priKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	key := APIKey{
		Key:    userDefault.APIKey,
		Secret: userDefault.APISecret,
	}

	return &key
}

// HasSavedAuth to judge if auth info is already saved in auth.yaml
func (cache *AuthCache) HasSavedAuth(baseHost string) bool {
	return cache.savedAuths.IsSet(baseHost)
}

// HasDefaultAuth to judge if auth info is exist from cmd args
func (cache *AuthCache) HasDefaultAuth() bool {
	return cache.DefaultID != "" && cache.DefaultPass.IsSet()
}

func (cache *AuthCache) retriveAuth() error {
	baseHost := GetBaseHost()

	if cache.DefaultID == "" || !cache.DefaultPass.IsSet() {
		if !cache.HasSavedAuth(baseHost) {
			return ErrAuthMissing
		}

		login := cache.savedAuths.Sub(baseHost)

		cache.DefaultID = login.GetString("identity")

		cache.DefaultPass.ShadowSet(login.GetString("password"))
	}

	fmt.Println("Login with identity:", cache.DefaultID)

	var loginAuth context.Context
	if loginAuth = cache.Login(
		cache.DefaultID, &cache.DefaultPass); loginAuth == nil {
		return fmt.Errorf(
			"login failed with identity: %s", cache.DefaultID)
	}

	var key *APIKey
	if key = cache.GetUserDefaultKey(loginAuth); key == nil {
		return fmt.Errorf(
			"retrive %s's api key from %s failed", cache.DefaultID, baseHost)
	}

	authInfo := Authentication{
		Identity: cache.DefaultID,
		Password: cache.DefaultPass,
		APIKey:   *key,
	}

	cache.authList = append(cache.authList, &authInfo)

	return nil
}

func (cache *AuthCache) readAuthFile(authFile string) error {
	if _, err := os.Stat(authFile); os.IsNotExist(err) {
		return err
	}

	csvFile, err := os.OpenFile(authFile, os.O_RDONLY, os.ModePerm)

	if err != nil {
		return err
	}

	var auths []*Authentication

	if err = gocsv.UnmarshalFile(csvFile, &auths); err != nil {
		return err
	}

	// nextIDX使用 uint32 CAS自增操作以实现go routine 安全的自旋，
	// 故验证信息总量不能超过max uint32
	if len(auths) > math.MaxUint32 {
		fmt.Println("auth info length can't be longer than max uint32.")
		auths = auths[:int(math.MaxUint32)]
	}

	for idx, authInfo := range auths {
		if !authInfo.Validate() {
			jsonBytes, _ := json.Marshal(authInfo)
			fmt.Printf("Record[%d]@line[%d] is invalid: %s",
				idx+1, idx+2, string(jsonBytes))

			continue
		}

		cache.authList = append(cache.authList, authInfo)
	}

	if len(cache.authList) < 1 {
		return fmt.Errorf("no valid auth info in file: %s", authFile)
	}

	return nil
}

// NextAuth get next auth context
func (cache *AuthCache) NextAuth(parent context.Context) context.Context {
	if parent == nil {
		parent = cache.rootCtx
	}

	authInfo := cache.authList[cache.nextIDX()]

	keyCtx, exist := cache.keyCtxCache[authInfo.Identity]
	if !exist {
		keyCtx = context.WithValue(
			parent, ngerest.ContextAPIKey, ngerest.APIKey{
				Key:    authInfo.Key,
				Secret: authInfo.Secret,
			})

		if authInfo.Identity != "" {
			cache.keyCtxCache[authInfo.Identity] = keyCtx
		}
	}

	return keyCtx
}

// NewAuthCache create new api auth cache
func NewAuthCache(ctx context.Context, clientHub *ClientHub) *AuthCache {
	if ctx == nil {
		panic("root context is nil.")
	}
	if clientHub == nil {
		panic("client hub is nil pointer.")
	}

	cache := AuthCache{
		savedAuths:  viper.New(),
		rootCtx:     ctx,
		clientHub:   clientHub,
		keyCtxCache: make(map[string]context.Context),
	}

	cache.savedAuths.SetKeyDelim(viperHostnameKeyDelim)

	// cache.savedAuths.

	return &cache
}
