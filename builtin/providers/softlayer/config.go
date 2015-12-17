package softlayer

import (
	"log"

	slclient "github.com/maximilien/softlayer-go/client"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type Config struct {
	Username string
	ApiKey string
}

type Client struct {
	networkApplicationDeliveryControllerService softlayer.SoftLayer_Network_Application_Delivery_Controller_Service
}

func (c *Config) Client() (*Client, error) {
	slc := slclient.NewSoftLayerClient(c.Username, c.ApiKey)
	networkApplicationDeliveryControllerService, err := slc.GetSoftLayer_Network_Application_Delivery_Controller_Service()

	if err != nil {
		return nil, err
	}

	client := &Client {
		networkApplicationDeliveryControllerService: networkApplicationDeliveryControllerService,
	}

	log.Println("[INFO] Created SoftLayer client")

	return client, nil
}
