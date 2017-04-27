package ibmcloud

import (
	"time"

	"github.com/IBM-Bluemix/bluemix-go/endpoints"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ibmid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IBM ID.",
				DefaultFunc: schema.EnvDefaultFunc("IBMID", ""),
			},
			"ibmid_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for the IBM ID.",
				DefaultFunc: schema.EnvDefaultFunc("IBMID_PASSWORD", ""),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Bluemix Region (for example 'us-south').",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_REGION", "BLUEMIX_REGION"}, "us-south"),
			},
			"bluemix_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any Bluemix API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_TIMEOUT", "BLUEMIX_TIMEOUT"}, 60),
			},
			"softlayer_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any SoftLayer API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_TIMEOUT", "SOFTLAYER_TIMEOUT"}, 60),
			},
			"softlayer_account_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SoftLayer IMS account number linked with IBM ID.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_ACCOUNT_NUMBER", "SOFTLAYER_ACCOUNT_NUMBER"}, ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ibmcloud_cf_account":           dataSourceIBMCloudCfAccount(),
			"ibmcloud_cf_org":               dataSourceIBMCloudCfOrg(),
			"ibmcloud_cf_service_instance":  dataSourceIBMCloudCfServiceInstance(),
			"ibmcloud_cf_service_key":       dataSourceIBMCloudCfServiceKey(),
			"ibmcloud_cf_service_plan":      dataSourceIBMCloudCfServicePlan(),
			"ibmcloud_cf_space":             dataSourceIBMCloudCfSpace(),
			"ibmcloud_cs_cluster":           dataSourceIBMCloudCsCluster(),
			"ibmcloud_cs_cluster_config":    dataSourceIBMCloudArmadaClusterConfig(),
			"ibmcloud_cs_worker":            dataSourceIBMCloudCsWorker(),
			"ibmcloud_infra_dns_domain":     dataSourceIBMCloudInfraDNSDomain(),
			"ibmcloud_infra_image_template": dataSourceIBMCloudInfraImageTemplate(),
			"ibmcloud_infra_ssh_key":        dataSourceIBMCloudInfraSSHKey(),
			"ibmcloud_infra_virtual_guest":  dataSourceIBMCloudInfraVirtualGuest(),
			"ibmcloud_infra_vlan":           dataSourceIBMCloudInfraVlan(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ibmcloud_cf_service_instance":               resourceIBMCloudCfServiceInstance(),
			"ibmcloud_cf_service_key":                    resourceIBMCloudCfServiceKey(),
			"ibmcloud_cs_cluster":                        resourceIBMCloudArmadaCluster(),
			"ibmcloud_cs_cluster_service_bind":           resourceIBMCloudClusterBindService(),
			"ibmcloud_infra_bare_metal":                  resourceIBMCloudInfraBareMetal(),
			"ibmcloud_infra_basic_monitor":               resourceIBMCloudInfraBasicMonitor(),
			"ibmcloud_infra_block_storage":               resourceIBMCloudInfraBlockStorage(),
			"ibmcloud_infra_dns_domain":                  resourceIBMCloudInfraDNSDomain(),
			"ibmcloud_infra_dns_domain_record":           resourceIBMCloudInfraDNSDomainRecord(),
			"ibmcloud_infra_file_storage":                resourceIBMCloudInfraFileStorage(),
			"ibmcloud_infra_fw_hardware_dedicated":       resourceIBMCloudInfraFwHardwareDedicated(),
			"ibmcloud_infra_fw_hardware_dedicated_rules": resourceIBMCloudInfraFwHardwareDedicatedRules(),
			"ibmcloud_infra_global_ip":                   resourceIBMCloudInfraGlobalIp(),
			"ibmcloud_infra_lb_local":                    resourceIBMCloudInfraLbLocal(),
			"ibmcloud_infra_lb_local_service":            resourceIBMCloudInfraLbLocalService(),
			"ibmcloud_infra_lb_local_service_group":      resourceIBMCloudInfraLbLocalServiceGroup(),
			"ibmcloud_infra_lb_vpx":                      resourceIBMCloudInfraLbVpx(),
			"ibmcloud_infra_lb_vpx_ha":                   resourceIBMCloudInfraLbVpxHa(),
			"ibmcloud_infra_lb_vpx_service":              resourceIBMCloudInfraLbVpxService(),
			"ibmcloud_infra_lb_vpx_vip":                  resourceIBMCloudInfraLbVpxVip(),
			"ibmcloud_infra_objectstorage_account":       resourceIBMCloudInfraObjectStorageAccount(),
			"ibmcloud_infra_provisioning_hook":           resourceIBMCloudInfraProvisioningHook(),
			"ibmcloud_infra_scale_group":                 resourceIBMCloudInfraScaleGroup(),
			"ibmcloud_infra_scale_policy":                resourceIBMCloudInfraScalePolicy(),
			"ibmcloud_infra_security_certificate":        resourceIBMCloudInfraSecurityCertificate(),
			"ibmcloud_infra_ssh_key":                     resourceIBMCloudInfraSSHKey(),
			"ibmcloud_infra_user":                        resourceIBMCloudInfraUser(),
			"ibmcloud_infra_virtual_guest":               resourceIBMCloudInfraVirtualGuest(),
			"ibmcloud_infra_vlan":                        resourceIBMCloudInfraVlan(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	softlayerTimeout := d.Get("softlayer_timeout").(int)
	bluemixTimeout := d.Get("bluemix_timeout").(int)

	region := d.Get("region").(string)
	endpointLocator := endpoints.NewEndpointLocator(region)
	iamEndpoint, err := endpointLocator.IAMEndpoint()
	if err != nil {
		return nil, err
	}
	config := Config{
		IBMID:                   d.Get("ibmid").(string),
		IBMIDPassword:           d.Get("ibmid_password").(string),
		Region:                  d.Get("region").(string),
		BluemixTimeout:          time.Duration(bluemixTimeout) * time.Second,
		SoftLayerTimeout:        time.Duration(softlayerTimeout) * time.Second,
		SoftLayerAccountNumber:  d.Get("softlayer_account_number").(string),
		IAMEndpoint:             iamEndpoint,
		RetryCount:              3,
		RetryDelay:              30 * time.Millisecond,
		SoftLayerEndpointURL:    SoftlayerRestEndpoint,
		SoftlayerXMLRPCEndpoint: SoftlayerXMLRPCEndpoint,
	}
	return config.ClientSession()
}
