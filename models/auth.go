package models

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// IdentityMap identity pattern map
type IdentityMap map[string]*regexp.Regexp

// CheckIdentity check & modify login map
func (idMap *IdentityMap) CheckIdentity(id string, login map[string]string) error {
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
func (idMap *IdentityMap) AddPattern(name string, pattern *regexp.Regexp) error {
	if _, exist := (*idMap)[name]; exist {
		return fmt.Errorf("named[%s] pattern already exists", name)
	}

	(*idMap)[name] = pattern

	return nil
}

// NewIdentityMap generate identity pattern map
func NewIdentityMap() IdentityMap {
	mp := make(IdentityMap)

	mp["email"] = regexp.MustCompile(`[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)*`)
	mp["mobile"] = regexp.MustCompile(`(\+?[0-9]{2,3})?[0-9-]{6,13}`)

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
	return p.ShadowSet(value)
}

// MarshalCSV marshal password to csv
func (p *Password) MarshalCSV() string {
	return p.String()
}

// UnmarshalJSON unmarshal from json string
func (p *Password) UnmarshalJSON(data []byte) error {
	return p.ShadowSet(strings.Trim(string(data), "\""))
}

// MarshalJSON marshal to json string
func (p *Password) MarshalJSON() ([]byte, error) {
	var buff bytes.Buffer
	buff.WriteString((*p).String())

	return buff.Bytes(), nil
}

// Shadow shadow & store a password
func (p *Password) Shadow(value string) string {
	encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, &p.key.PublicKey, []byte(value))
	p.shadowed = base64.StdEncoding.EncodeToString(encrypted)
	return p.shadowed
}

// Show show unshadowed password
func (p *Password) Show() string {
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

// NewPassword get new password struct
func NewPassword() *Password {
	var defaultKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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

	block, rest := pem.Decode(defaultKey)
	if block == nil || len(rest) > 0 {
		panic("pem decode key string failed.")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("parse pkcs8 private key failed: " + err.Error())
	}

	pass := Password{
		key: key,
	}

	return &pass
}

// Authentication auth info
type Authentication struct {
	Identity string   `csv:"identity" json:"identity"`
	Password Password `csv:"password" json:"password"`
	APIKey
}

// APIKey key & secret for api
type APIKey struct {
	Key    string `csv:"api_key" json:"api_key"`
	Secret string `csv:"api_secret" json:"api_secret"`
}
