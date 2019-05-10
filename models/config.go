package models

import (
	"fmt"
	"strconv"
	"sync"

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
	initFlag   sync.Once
}

func (hub *ClientHub) init() {
	hub.initFlag.Do(func() {
		hub.clientsMap = make(map[string]*ngerest.APIClient)
	})
}

// GetClient to get client instance by host string
func (hub *ClientHub) GetClient(host string) (*ngerest.APIClient, error) {
	hub.init()

	if host == "" {
		return nil, ErrHost
	}

	if client, exist := hub.clientsMap[host]; exist {
		return client, nil
	}

	cfg := ngerest.NewConfiguration()
	client := ngerest.NewAPIClient(cfg)

	hostURL := viper.GetString("scheme") + "://" + host
	client.ChangeBasePath(hostURL + viper.GetString("base-uri"))
	fmt.Println("Change host to:", hostURL)

	hub.clientsMap[host] = client

	return client, nil
}
