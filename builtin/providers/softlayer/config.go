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
	productOrderService            softlayer.SoftLayer_Product_Order_Service
	dnsDomainResourceRecordService softlayer.SoftLayer_Dns_Domain_Resource_Record_Service
	dnsDomainService               softlayer.SoftLayer_Dns_Domain_Service
}

func (c *Config) Client() (*Client, error) {
	slc := slclient.NewSoftLayerClient(c.Username, c.ApiKey)

	dnsDomainService, err := slc.GetSoftLayer_Dns_Domain_Service()

	if err != nil {
		return nil, err
	}

	dnsDomainResourceRecordService, err := slc.GetSoftLayer_Dns_Domain_Resource_Record_Service()

	client := &Client {
		dnsDomainService :   		    dnsDomainService,
		dnsDomainResourceRecordService: dnsDomainResourceRecordService,
	}

	log.Println("[INFO] Created SoftLayer client")

	return client, nil
}
