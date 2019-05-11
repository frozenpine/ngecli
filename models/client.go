package models

import (
	"fmt"
	"sync"

	"github.com/frozenpine/ngecli/common"

	"github.com/frozenpine/ngerest"

	"github.com/spf13/viper"
)

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
		return nil, common.ErrHost
	}

	if client, exist := hub.clientsMap[host]; exist {
		return client, nil
	}

	cfg := ngerest.NewConfiguration()
	client := ngerest.NewAPIClient(cfg)

	hostURL := viper.GetString("scheme") + "://" + host
	client.ChangeBasePath(hostURL + viper.GetString("base-uri"))
	if hostURL != common.GetBaseURL() {
		fmt.Println("Change host to:", hostURL)
	}

	hub.clientsMap[host] = client

	return client, nil
}
