package softlayer

import (
	"fmt"
	"log"
	"strconv"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
	"strings"
	"time"
)

const (
	NETSCALER_VPX_TYPE = "Netscaler VPX"
)

func resourceSoftLayerNetworkApplicationDeliveryController() *schema.Resource {
	return &schema.Resource{
		Create:   resourceSoftLayerNetworkApplicationDeliveryControllerCreate,
		Read:     resourceSoftLayerNetworkApplicationDeliveryControllerRead,
		Delete:   resourceSoftLayerNetworkApplicationDeliveryControllerDelete,
		Exists:   resourceSoftLayerNetworkApplicationDeliveryControllerExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"front_end_vlan": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan_number": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"primary_router_hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"front_end_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"back_end_vlan": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlan_number": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"primary_router_hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"back_end_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vip_pool": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getVlanId(vlanNumber int, primaryRouterHostname string, meta interface{}) (int, error) {
	client := meta.(*Client).accountService

	mask := []string{
		"id",
	}

	filter := fmt.Sprintf(
		`{"networkVlans":{"primaryRouter":{"hostname":{"operation":"%s"}},`+
			`"vlanNumber":{"operation":%d}}}`,
		primaryRouterHostname,
		vlanNumber)
	client.GetNetworkStorage()
	networkVlan, err := client.GetNetworkVlans(mask, filter)

	if err != nil {
		return 0, fmt.Errorf("Error looking up Vlan: %s", err)
	}

	if len(networkVlan) < 1 {
		return 0, fmt.Errorf(
			"Unable to locate a vlan matching the provided router hostname and vlan number: %s/%d",
			primaryRouterHostname,
			vlanNumber)
	}
	return networkVlan[0].Id, nil
}

func getDatacenterId(name string, meta interface{}) (int, error) {
	client := meta.(*Client).locationDatacenterService

	mask := []string{
		"id",
	}

	filter := fmt.Sprintf(
		`{"name":{"operation":"%s"}}`,
		name)

	datacenters, err := client.GetDatacenters(mask, filter)

	if err != nil {
		return 0, fmt.Errorf("Error looking up Vlan: %s", err)
	}

	if len(datacenters) < 1 {
		return 0, fmt.Errorf(
			"Unable to find a datacenter with a name: %s",
			name)
	}
	return datacenters[0].Id, nil
}

func getSubnetId(subnet string, meta interface{}) (int, error) {
	client := meta.(*Client).accountService

	mask := []string{
		"id",
	}

	subnetInfo := strings.Split(subnet, "/")
	if len(subnetInfo) != 2 {
		return 0, fmt.Errorf(
			"Unable to parse the provided subnet: %s", subnet)
	}

	networkIdentifier := subnetInfo[0]
	cidr := subnetInfo[1]

	filter := fmt.Sprintf(
		`{"subnets":{"cidr":{"operation":%s},`+
			`"networkIdentifier":{"operation":"%s"}}}`,
		cidr,
		networkIdentifier)

	subnets, err := client.GetSubnets(mask, filter)

	if err != nil {
		return 0, fmt.Errorf("Error looking up Subnet: %s", err)
	}

	if len(subnets) < 1 {
		return 0, fmt.Errorf(
			"Unable to locate a subnet matching the provided subnet: %s", subnet)
	}
	return subnets[0].Id, nil
}

func resourceSoftLayerNetworkApplicationDeliveryControllerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	nadcType := NETSCALER_VPX_TYPE

	switch nadcType {
	default:
		return fmt.Errorf("[ERROR] Network application delivery controller type %s is not supported", nadcType)
	case NETSCALER_VPX_TYPE:
		// create Netscaler VPX
		opts := softlayer.NetworkApplicationDeliveryControllerCreateOptions{
			Speed:   d.Get("speed").(int),
			Version: d.Get("version").(string),
			Plan:    d.Get("plan").(string),
			IpCount: d.Get("ip_count").(int),
		}

		if len(d.Get("datacenter").(string)) > 0 {
			datacenterId, err := getDatacenterId(d.Get("datacenter").(string), meta)
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			opts.Location = strconv.Itoa(datacenterId)
		}

		opts.Hardware = make([]datatypes.SoftLayer_Hardware_Template, 1)

		if len(d.Get("front_end_vlan.vlan_number").(string)) > 0 || len(d.Get("front_end_subnet").(string)) > 0 {
			opts.Hardware[0].PrimaryNetworkComponent = &datatypes.SoftLayer_Network_Component{}
		}

		if len(d.Get("front_end_vlan.vlan_number").(string)) > 0 {
			vlanNumber, err := strconv.Atoi(d.Get("front_end_vlan.vlan_number").(string))
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			networkVlanId, err := getVlanId(vlanNumber, d.Get("front_end_vlan.primary_router_hostname").(string), meta)
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			opts.Hardware[0].PrimaryNetworkComponent.NetworkVlanId = networkVlanId
		}

		if len(d.Get("front_end_subnet").(string)) > 0 {
			primarySubnetId, err := getSubnetId(d.Get("front_end_subnet").(string), meta)
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			opts.Hardware[0].PrimaryNetworkComponent.NetworkVlan = &datatypes.SoftLayer_Network_Vlan_Template{
				PrimarySubnetId: primarySubnetId,
			}
		}

		if len(d.Get("back_end_vlan.vlan_number").(string)) > 0 || len(d.Get("back_end_subnet").(string)) > 0 {
			opts.Hardware[0].PrimaryBackendNetworkComponent = &datatypes.SoftLayer_Network_Component{}
		}

		if len(d.Get("back_end_vlan.vlan_number").(string)) > 0 {
			vlanNumber, err := strconv.Atoi(d.Get("back_end_vlan.vlan_number").(string))
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			networkVlanId, err := getVlanId(vlanNumber, d.Get("back_end_vlan.primary_router_hostname").(string), meta)
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			opts.Hardware[0].PrimaryBackendNetworkComponent.NetworkVlanId = networkVlanId
		}
		if len(d.Get("back_end_subnet").(string)) > 0 {
			primarySubnetId, err := getSubnetId(d.Get("back_end_subnet").(string), meta)
			if err != nil {
				return fmt.Errorf("Error creating network application delivery controller: %s", err)
			}
			opts.Hardware[0].PrimaryBackendNetworkComponent.NetworkVlan = &datatypes.SoftLayer_Network_Vlan_Template{
				PrimarySubnetId: primarySubnetId,
			}
		}

		log.Printf("[INFO] Creating network application delivery controller")

		netscalerVPX, err := client.CreateNetscalerVPX(&opts)

		if err != nil {
			return fmt.Errorf("Error creating network application delivery controller: %s", err)
		}

		d.SetId(fmt.Sprintf("%d", netscalerVPX.Id))

		log.Printf("[INFO] Netscaler VPX ID: %s", d.Id())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	IsVipReady := false
	// Wait Virtual IP provisioning
	for vipWaitCount := 0; vipWaitCount < 60; vipWaitCount++ {
		getObjectResult, err := client.GetObject(id)
		if err != nil {
			return fmt.Errorf("Error retrieving network application delivery controller: %s", err)
		}

		ipCount := 0
		if getObjectResult.Subnets != nil {
			ipCount = len(getObjectResult.Subnets[0].IpAddresses)
		}
		if ipCount > 0 {
			IsVipReady = true
			break
		}
		log.Printf("[INFO] Wait 10 seconds for Virtual IP provisioning on Netscaler VPX ID: %d", id)
		time.Sleep(time.Second * 10)
	}

	if !IsVipReady {
		return fmt.Errorf("Failed to create VIPs for Netscaler VPX ID: %d", id)
	}
	return resourceSoftLayerNetworkApplicationDeliveryControllerRead(d, meta)
}

func resourceSoftLayerNetworkApplicationDeliveryControllerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	getObjectResult, err := client.GetObject(id)
	if err != nil {
		return fmt.Errorf("Error retrieving network application delivery controller: %s", err)
	}

	d.Set("name", getObjectResult.Name)
	d.Set("type", getObjectResult.Type.Name)
	if getObjectResult.Datacenter != nil {
		d.Set("datacenter", getObjectResult.Datacenter.Name)
	}

	frontEndVlan := d.Get("front_end_vlan").(map[string]interface{})
	backEndVlan := d.Get("back_end_vlan").(map[string]interface{})
	frontEndSubnet := ""
	backEndSubnet := ""

	for _, vlan := range getObjectResult.NetworkVlans {
		if vlan.PrimaryRouter != nil && vlan.PrimaryRouter.Hostname != "" && vlan.VlanNumber > 0 {
			isFcr, _ := regexp.MatchString("fcr", vlan.PrimaryRouter.Hostname)
			isBcr, _ := regexp.MatchString("bcr", vlan.PrimaryRouter.Hostname)
			if isFcr {
				frontEndVlan["primary_router_hostname"] = vlan.PrimaryRouter.Hostname
				vlanNumber := strconv.Itoa(vlan.VlanNumber)
				frontEndVlan["vlan_number"] = vlanNumber
				if vlan.PrimarySubnets != nil && len(vlan.PrimarySubnets) > 0 {
					ipAddress := vlan.PrimarySubnets[0].NetworkIdentifier
					cidr := strconv.Itoa(vlan.PrimarySubnets[0].Cidr)
					frontEndSubnet = ipAddress + "/" + cidr
				}
			}

			if isBcr {
				backEndVlan["primary_router_hostname"] = vlan.PrimaryRouter.Hostname
				vlanNumber := strconv.Itoa(vlan.VlanNumber)
				backEndVlan["vlan_number"] = vlanNumber
				if vlan.PrimarySubnets != nil && len(vlan.PrimarySubnets) > 0 {
					ipAddress := vlan.PrimarySubnets[0].NetworkIdentifier
					cidr := strconv.Itoa(vlan.PrimarySubnets[0].Cidr)
					backEndSubnet = ipAddress + "/" + cidr
				}
			}
		}
	}

	d.Set("front_end_vlan", frontEndVlan)
	d.Set("back_end_vlan", backEndVlan)
	d.Set("front_end_subnet", frontEndSubnet)
	d.Set("back_end_subnet", backEndSubnet)

	vips := make([]string, 0)
	ipCount := 0
	for i, subnet := range getObjectResult.Subnets {
		for _, ipAddressObj := range subnet.IpAddresses {
			vips = append(vips, ipAddressObj.IpAddress)
			if i == 0 {
				ipCount++
			}
		}
	}

	d.Set("vip_pool", vips)
	d.Set("ip_count", ipCount)

	description := getObjectResult.Description
	r, _ := regexp.Compile(" [0-9]+Mbps")
	speedStr := r.FindString(description)
	r, _ = regexp.Compile("[0-9]+")
	speed, err := strconv.Atoi(r.FindString(speedStr))
	if err == nil && speed > 0 {
		d.Set("speed", speed)
	}

	r, _ = regexp.Compile(" VPX [0-9]+\\.[0-9]+ ")
	versionStr := r.FindString(description)
	r, _ = regexp.Compile("[0-9]+\\.[0-9]+")
	version := r.FindString(versionStr)
	if version != "" {
		d.Set("version", version)
	}

	r, _ = regexp.Compile(" [A-Za-z]+$")
	planStr := r.FindString(description)
	r, _ = regexp.Compile("[A-Za-z]+$")
	plan := r.FindString(planStr)
	if plan != "" {
		d.Set("plan", plan)
	}

	return nil
}

func resourceSoftLayerNetworkApplicationDeliveryControllerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = client.DeleteObject(id)

	if err != nil {
		return fmt.Errorf("Error deleting network application delivery controller: %s", err)
	}

	return nil
}

func resourceSoftLayerNetworkApplicationDeliveryControllerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	nadc, err := client.GetObject(id)

	if err != nil {
		return false, fmt.Errorf("Error fetching network application delivery controller: %s", err)
	}

	return nadc.Id == id && err == nil, nil
}
