package models

import (
	"strconv"

	"github.com/frozenpine/ngerest"

	"github.com/spf13/viper"
)

// GetBasePath to get base uri path
func GetBasePath() string {
	baseURI := viper.GetString("base-uri")

	return GetBaseURL() + baseURI
}

// GetBaseHost to get base host:port string
func GetBaseHost() string {
	port := viper.GetInt("port")
	host := viper.GetString("host")

	if port != 80 {
		return host + ":" + strconv.Itoa(port)
	}

	return host
}

// GetBaseURL to get base full url path
func GetBaseURL() string {
	scheme := viper.GetString("scheme")

	return scheme + "://" + GetBaseHost()
}

// ClientHub is a hub of clients from hosts
type ClientHub struct {
	clientsMap map[string]*ngerest.APIClient
}

// GetClient to get client instance by host string
func (hub *ClientHub) GetClient(host string) (*ngerest.APIClient, error) {
	if host == "" {
		return nil, ErrHost
	}

	if hub.clientsMap == nil {
		hub.clientsMap = make(map[string]*ngerest.APIClient)
	}

	if client, exist := hub.clientsMap[host]; exist {
		return client, nil
	}

	cfg := ngerest.NewConfiguration()
	client := ngerest.NewAPIClient(cfg)
	client.ChangeBasePath(host)

	hub.clientsMap[host] = client

	return client, nil
}