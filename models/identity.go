package models

import (
	"errors"
	"fmt"
	"regexp"
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

// APIKey key & secret for api
type APIKey struct {
	Key    string `csv:"api_key" json:"api_key"`
	Secret string `csv:"api_secret" json:"api_secret"`
}
