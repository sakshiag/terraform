package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TheWeatherCompany/softlayer-go/common"
	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
	"strconv"
	"strings"
	"time"
)

const (
	PACKAGE_TYPE_APPLICATION_DELIVERY_CONTROLLER = "ADDITIONAL_SERVICES_APPLICATION_DELIVERY_APPLIANCE"
	ORDER_TYPE_APPLICATION_DELIVERY_CONTROLLER   = "SoftLayer_Container_Product_Order_Network_Application_Delivery_Controller"
	PACKAGE_ID_APPLICATION_DELIVERY_CONTROLLER   = 192
	DELIMITER                                    = "_"
)

type softLayer_Network_Application_Delivery_Controller_Service struct {
	client softlayer.Client
}

func NewSoftLayer_Network_Application_Delivery_Controller_Service(client softlayer.Client) *softLayer_Network_Application_Delivery_Controller_Service {
	return &softLayer_Network_Application_Delivery_Controller_Service{
		client: client,
	}
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) GetName() string {
	return "SoftLayer_Network_Application_Delivery_Controller"
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) CreateNetscalerVPX(createOptions *softlayer.NetworkApplicationDeliveryControllerCreateOptions) (datatypes.SoftLayer_Network_Application_Delivery_Controller, error) {
	err := slnadcs.checkCreateVpxRequiredValues(createOptions)
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	orderService, err := slnadcs.client.GetSoftLayer_Product_Order_Service()
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	items, err := slnadcs.FindCreatePriceItems(createOptions)
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	order := datatypes.SoftLayer_Container_Product_Order_Network_Application_Delivery_Controller{
		PackageId:   PACKAGE_ID_APPLICATION_DELIVERY_CONTROLLER,
		ComplexType: ORDER_TYPE_APPLICATION_DELIVERY_CONTROLLER,
		Location:    createOptions.Location,
		Prices:      items,
		Quantity:    1,
	}

	receipt, err := orderService.PlaceContainerOrderApplicationDeliveryController(order)
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	vpx, err := slnadcs.findVPXByOrderId(receipt.OrderId)
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	return vpx, nil
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) CreateVirtualIpAddress(nadcId int, template datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template) (bool, error) {
	nadc, err := slnadcs.GetObject(nadcId)
	if err != nil {
		return false, err
	}
	if nadc.Id != nadcId {
		return false, fmt.Errorf("Network application delivery controller with id '%d' is not found", nadcId)
	}

	parameters := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template_Parameters{
		Parameters: []datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template{
			template,
		},
	}

	requestBody, err := json.Marshal(parameters)
	if err != nil {
		return false, fmt.Errorf("Network application delivery controller with id '%d' is not found: %s", nadcId, err)
	}

	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequest(fmt.Sprintf("%s/%d/%s.json", slnadcs.GetName(), nadcId, "createLiveLoadBalancer"), "POST", bytes.NewBuffer(requestBody))
	if err != nil {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#createVirtualIpAddress, error message '%s'", err.Error())
		return false, errors.New(errorMessage)
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#createVirtualIpAddress, HTTP error code: '%d'", errorCode)
		return false, errors.New(errorMessage)
	}

	if response_value := string(response[:]); response_value != "true" {
		return false, fmt.Errorf("Failed to create Virtual IP Address with '%s' name from network application delivery controller with '%d' id. Got '%s' as response from the API", template.Name, nadcId, response_value)
	}

	return true, nil
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) DeleteVirtualIpAddress(nadcId int, name string) (bool, error) {
	nadc, err := slnadcs.GetObject(nadcId)
	if err != nil {
		return false, err
	}
	if nadc.Id != nadcId {
		return false, fmt.Errorf("Network application delivery controller with id '%d' is not found", nadcId)
	}

	parameters := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template_Parameters{
		Parameters: []datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template{
			datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template{
				Name: name,
			},
		},
	}

	requestBody, err := json.Marshal(parameters)
	if err != nil {
		return false, err
	}

	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequest(fmt.Sprintf("%s/%d/%s.json", slnadcs.GetName(), nadcId, "deleteLiveLoadBalancer"), "POST", bytes.NewBuffer(requestBody))
	if err != nil {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#deleteVirtualIpAddress, error message '%s'", err.Error())
		return false, errors.New(errorMessage)
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#deleteVirtualIpAddress, HTTP error code: '%d'", errorCode)
		return false, errors.New(errorMessage)
	}

	if response_value := string(response[:]); response_value != "true" {
		return false, fmt.Errorf("Failed to delete Virtual IP Address with name '%s' from network application delivery controller with id '%d'. Got '%s' as response from the API", name, nadcId, response_value)
	}

	return true, err
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) EditVirtualIpAddress(nadcId int, template datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template) (bool, error) {
	nadc, err := slnadcs.GetObject(nadcId)
	if err != nil {
		return false, err
	}
	if nadc.Id != nadcId {
		return false, fmt.Errorf("Network application delivery controller with id '%d' is not found", nadcId)
	}

	parameters := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template_Parameters{
		Parameters: []datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Template{
			template,
		},
	}

	requestBody, err := json.Marshal(parameters)
	if err != nil {
		return false, err
	}

	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequest(fmt.Sprintf("%s/%d/%s.json", slnadcs.GetName(), nadcId, "updateLiveLoadBalancer"), "POST", bytes.NewBuffer(requestBody))
	if err != nil {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#editVirtualIpAddress, error message '%s'", err.Error())
		return false, errors.New(errorMessage)
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#editVirtualIpAddress, HTTP error code: '%d'", errorCode)
		return false, errors.New(errorMessage)
	}

	if response_value := string(response[:]); response_value != "true" {
		return false, fmt.Errorf("Failed to update Virtual IP Address with id '%d' from network application delivery controller with id '%d'. Got '%s' as response from the API", template.Id, nadcId, response_value)
	}

	return true, err
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) GetVirtualIpAddress(nadcId int, vipName string) (datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress, error) {
	nadc, err := slnadcs.GetObject(nadcId)
	if err != nil {
		return datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress{}, err
	}
	if nadc.Id != nadcId {
		return datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress{}, fmt.Errorf("Network application delivery controller with id '%d' is not found", nadcId)
	}

	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequest(fmt.Sprintf("%s/%d/%s.json", slnadcs.GetName(), nadcId, "getLoadBalancers"), "GET", new(bytes.Buffer))
	if err != nil {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#getVirtualIpAddress, error message '%s'", err.Error())
		return datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress{}, errors.New(errorMessage)
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#getVirtualIpAddress, HTTP error code: '%d'", errorCode)
		return datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress{}, errors.New(errorMessage)
	}

	addresses := datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress_Array{}
	err = json.Unmarshal(response, &addresses)
	if err != nil {
		return datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress{}, err
	}

	var result datatypes.SoftLayer_Network_LoadBalancer_VirtualIpAddress
	for _, address := range addresses {
		if address.Name == vipName {
			result = address
			break
		}
	}

	return result, err
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) GetObject(id int) (datatypes.SoftLayer_Network_Application_Delivery_Controller, error) {

	objectMask := []string{
		"id",
		"createDate",
		"name",
		"typeId",
		"modifyDate",
		"description",
		"managedResourceFlag",
		"managementIpAddress",
		"primaryIpAddress",
		"password",
		"notes",
		"datacenter",
		"averageDailyPublicBandwidthUsage",
		"licenseExpirationDate",
		"networkVlan",
		"networkVlanCount",
		"networkVlans",
		"subnetCount",
		"subnets",
		"tagReferenceCount",
		"tagReferences",
		"type",
		"virtualIpAddressCount",
		"virtualIpAddresses",
	}

	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequestWithObjectMask(fmt.Sprintf("%s/%d/getObject.json", slnadcs.GetName(), id), objectMask, "GET", new(bytes.Buffer))
	if err != nil {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#getObject, error message '%s'", err.Error())
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, errors.New(errorMessage)
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not perform SoftLayer_Network_Application_Delivery_Controller#getObject, HTTP error code: '%d'", errorCode)
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, errors.New(errorMessage)
	}

	nadc := datatypes.SoftLayer_Network_Application_Delivery_Controller{}
	err = json.Unmarshal(response, &nadc)
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	return nadc, nil
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) GetBillingItem(volumeId int) (datatypes.SoftLayer_Billing_Item, error) {

	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequest(fmt.Sprintf("%s/%d/getBillingItem.json", slnadcs.GetName(), volumeId), "GET", new(bytes.Buffer))
	if err != nil {
		return datatypes.SoftLayer_Billing_Item{}, err
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not SoftLayer_NetWork_Storage#getBillingItem, HTTP error code: '%d'", errorCode)
		return datatypes.SoftLayer_Billing_Item{}, errors.New(errorMessage)
	}

	billingItem := datatypes.SoftLayer_Billing_Item{}
	err = json.Unmarshal(response, &billingItem)
	if err != nil {
		return datatypes.SoftLayer_Billing_Item{}, err
	}

	return billingItem, nil
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) DeleteObject(id int) (bool, error) {
	billingItem, err := slnadcs.GetBillingItem(id)
	if err != nil {
		return false, err
	}
	fmt.Println(billingItem.Id)
	if billingItem.Id > 0 {
		billingItemService, err := slnadcs.client.GetSoftLayer_Billing_Item_Service()
		if err != nil {
			return false, err
		}

		deleted, err := billingItemService.CancelService(billingItem.Id)
		if err != nil {
			return false, err
		}

		if deleted {
			return false, nil
		}
	}

	fmt.Errorf("softlayer-go: could not SoftLayer_Network_Storage_Service#deleteIscsiVolume with id: '%d'", id)

	return true, err
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) CancelService(billingId int) (bool, error) {
	response, errorCode, err := slnadcs.client.GetHttpClient().DoRawHttpRequest(fmt.Sprintf("%s/%d/cancelService.json", slnadcs.GetName(), billingId), "GET", new(bytes.Buffer))
	if err != nil {
		return false, err
	}

	if res := string(response[:]); res != "true" {
		return false, nil
	}

	if common.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("softlayer-go: could not SoftLayer_Billing_Item#CancelService, HTTP error code: '%d'", errorCode)
		return false, errors.New(errorMessage)
	}

	return true, err
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) FindCreatePriceItems(createOptions *softlayer.NetworkApplicationDeliveryControllerCreateOptions) ([]datatypes.SoftLayer_Product_Item_Price, error) {
	items, err := slnadcs.getApplicationDeliveryControllerItems()
	if err != nil {
		return []datatypes.SoftLayer_Product_Item_Price{}, err
	}

	nadcKey := slnadcs.getVPXPriceItemKeyName(createOptions.Version, createOptions.Speed, createOptions.Plan)
	ipKey := slnadcs.getPublicIpItemKeyName(createOptions.IpCount)

	var nadcItemPrice, ipItemPrice datatypes.SoftLayer_Product_Item_Price

	for _, item := range items {
		itemKey := item.Key
		if itemKey == nadcKey {
			nadcItemPrice = item.Prices[0]
		}
		if itemKey == ipKey {
			ipItemPrice = item.Prices[0]
		}
	}

	var errorMessages []string

	if nadcItemPrice.Id == 0 {
		errorMessages = append(errorMessages, fmt.Sprintf("VPX version, speed or plan have incorrect values"))
	}

	if ipItemPrice.Id == 0 {
		errorMessages = append(errorMessages, fmt.Sprintf("Ip quantity value is incorrect"))
	}

	if len(errorMessages) > 0 {
		err = errors.New(strings.Join(errorMessages, "\n"))
		return []datatypes.SoftLayer_Product_Item_Price{}, err
	}

	return []datatypes.SoftLayer_Product_Item_Price{nadcItemPrice, ipItemPrice}, nil
}

// Private methods

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) checkCreateVpxRequiredValues(createOptions *softlayer.NetworkApplicationDeliveryControllerCreateOptions) error {
	var err error
	var errorMessages []string
	errorTemplate := "* %s is required and cannot be empty\n"

	if createOptions.Plan == "" {
		errorMessages = append(errorMessages, fmt.Sprintf(errorTemplate, "Vpx Plan"))
	}

	if createOptions.Speed <= 0 {
		errorMessages = append(errorMessages, fmt.Sprintf(errorTemplate, "Network speed"))
	}

	if createOptions.Version == "" {
		errorMessages = append(errorMessages, fmt.Sprintf(errorTemplate, "Vpx version"))
	}

	if createOptions.Location == "" {
		errorMessages = append(errorMessages, fmt.Sprintf(errorTemplate, "Location"))
	}

	if len(errorMessages) > 0 {
		err = errors.New(strings.Join(errorMessages, "\n"))
	}

	return err
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) findVPXByOrderId(orderId int) (datatypes.SoftLayer_Network_Application_Delivery_Controller, error) {
	ObjectFilter := string(`{"applicationDeliveryControllers":{"billingItem":{"orderItem":{"order":{"id":{"operation":` + strconv.Itoa(orderId) + `}}}}}}`)

	// TODO: NEED TO ADD LOGIC TO CHECK IF EXISTS...VPX NOT CREATED IMMEDIATELY (Sleep is temporary hack)
	time.Sleep(30 * time.Second)
	accountService, err := slnadcs.client.GetSoftLayer_Account_Service()
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}
	vpxs, err := accountService.GetApplicationDeliveryControllersWithFilter(ObjectFilter)
	if err != nil {
		return datatypes.SoftLayer_Network_Application_Delivery_Controller{}, err
	}

	if len(vpxs) == 1 {
		return vpxs[0], nil
	}

	return datatypes.SoftLayer_Network_Application_Delivery_Controller{},
		fmt.Errorf("Cannot find Application Delivery Controller with order id '%d'", orderId)
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) getApplicationDeliveryControllerItems() ([]datatypes.SoftLayer_Product_Item, error) {
	productPackageService, err := slnadcs.client.GetSoftLayer_Product_Package_Service()
	if err != nil {
		return []datatypes.SoftLayer_Product_Item{}, err
	}

	return productPackageService.GetItemsByType(PACKAGE_TYPE_APPLICATION_DELIVERY_CONTROLLER)
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) getVPXPriceItemKeyName(version string, speed int, plan string) string {
	name := "CITRIX_NETSCALER_VPX"
	speedMeasurements := "MBPS"
	versionReplaced := strings.Replace(version, ".", DELIMITER, -1)
	speedString := strconv.Itoa(speed) + speedMeasurements
	return strings.Join([]string{name, versionReplaced, speedString, strings.ToUpper(plan)}, DELIMITER)
}

func (slnadcs *softLayer_Network_Application_Delivery_Controller_Service) getPublicIpItemKeyName(ipCount int) string {
	name := "STATIC_PUBLIC_IP_ADDRESSES"
	ipCountString := strconv.Itoa(ipCount)

	return strings.Join([]string{ipCountString, name}, DELIMITER)
}
