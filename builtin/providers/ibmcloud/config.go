package ibmcloud

import (
	"log"
	"time"

	slsession "github.com/softlayer/softlayer-go/session"

	"github.com/IBM-Bluemix/bluemix-go/api/account/accountv2"
	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/api/k8scluster/k8sclusterv1"
	bxsession "github.com/IBM-Bluemix/bluemix-go/session"
)

//Config stores user provider input config and the API endpoints
type Config struct {
	//The IBM ID
	IBMID string
	//Password fo the IBM ID
	IBMIDPassword string

	//Bluemix region
	Region string
	//Bluemix API timeout
	BluemixTimeout time.Duration

	//Softlayer end point url
	SoftLayerEndpointURL string
	//SoftlayerXMLRPCEndpoint endpoint
	SoftlayerXMLRPCEndpoint string
	//Softlayer API timeout
	SoftLayerTimeout time.Duration
	// Softlayer Account Number
	SoftLayerAccountNumber string

	//IAM endpoint
	IAMEndpoint string

	//Retry Count for API calls
	//Unexposed in the schema at this point as they are used only during session creation for a few calls
	//When sdk implements it we an expose them for expected behaviour
	//https://github.com/softlayer/softlayer-go/issues/41
	RetryCount int
	//Constant Retry Delay for API calls
	RetryDelay time.Duration
}

// ClientSession  contains  Bluemix/SoftLayer session and clients
type ClientSession interface {
	SoftLayerSession() *slsession.Session
	BluemixSession() *bxsession.Session

	ClusterClient() k8sclusterv1.Clusters
	ClusterWorkerClient() k8sclusterv1.Workers
	ClusterSubnetClient() k8sclusterv1.Subnets
	ClusterWebHooksClient() k8sclusterv1.Webhooks

	CloudFoundryOrgClient() cfv2.Organizations
	CloudFoundryServiceInstanceClient() cfv2.ServiceInstances
	CloudFoundryServicePlanClient() cfv2.ServicePlans
	CloudFoundryServiceKeyClient() cfv2.ServiceKeys
	CloudFoundryServiceOfferingClient() cfv2.ServiceOfferings
	CloudFoundrySpaceClient() cfv2.Spaces

	BluemixAcccountClient() accountv2.Accounts
}

//clientSession implements the ClientSession interface
type clientSession struct {
	session *Session

	csClient  k8sclusterv1.Clusters
	csWorker  k8sclusterv1.Workers
	csSubnet  k8sclusterv1.Subnets
	csWebHook k8sclusterv1.Webhooks

	cfOrgClient              cfv2.Organizations
	cfServiceInstanceClient  cfv2.ServiceInstances
	cfSpaceClient            cfv2.Spaces
	cfServicePlanClient      cfv2.ServicePlans
	cfServiceKeysClient      cfv2.ServiceKeys
	cfServiceOfferingsClient cfv2.ServiceOfferings

	bluemixAccountClient accountv2.Accounts
}

// SoftLayerSession providers SoftLayer Session
func (sess clientSession) SoftLayerSession() *slsession.Session {
	return sess.session.SoftLayerSession
}

// CloudFoundryOrgClient providers Cloud Foundary org APIs
func (sess clientSession) CloudFoundryOrgClient() cfv2.Organizations {
	return sess.cfOrgClient
}

// CloudFoundrySpaceClient providers Cloud Foundary space APIs
func (sess clientSession) CloudFoundrySpaceClient() cfv2.Spaces {
	return sess.cfSpaceClient
}

// CloudFoundryServiceInstanceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceInstanceClient() cfv2.ServiceInstances {
	return sess.cfServiceInstanceClient
}

// CloudFoundryServiceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServicePlanClient() cfv2.ServicePlans {
	return sess.cfServicePlanClient
}

// CloudFoundryServiceKeyClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceKeyClient() cfv2.ServiceKeys {
	return sess.cfServiceKeysClient
}

// CloudFoundryServiceClient providers Cloud Foundary service APIs
func (sess clientSession) CloudFoundryServiceOfferingClient() cfv2.ServiceOfferings {
	return sess.cfServiceOfferingsClient
}

// BluemixAcccountClient providers Bluemix Account APIs
func (sess clientSession) BluemixAcccountClient() accountv2.Accounts {
	return sess.bluemixAccountClient
}

// ClusterClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterClient() k8sclusterv1.Clusters {
	return sess.csClient
}

// ClusterWorkerClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterWorkerClient() k8sclusterv1.Workers {
	return sess.csWorker
}

// ClusterSubnetClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterSubnetClient() k8sclusterv1.Subnets {
	return sess.csSubnet
}

// ClusterWebHooksClient providers Bluemix Kubernetes Cluster APIs
func (sess clientSession) ClusterWebHooksClient() k8sclusterv1.Webhooks {
	return sess.csWebHook
}

// BluemixSession to provide the Bluemix Session
func (sess clientSession) BluemixSession() *bxsession.Session {
	return sess.session.BluemixSession
}

// ClientSession configures and returns a fully initialized ClientSession
func (c *Config) ClientSession() (interface{}, error) {

	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}

	cfClient, err := cfv2.NewClient(sess.BluemixSession)

	if err != nil {
		return nil, err
	}

	orgAPI := cfClient.Organizations()
	spaceAPI := cfClient.Spaces()
	serviceInstanceAPI := cfClient.ServiceInstances()
	servicePlanAPI := cfClient.ServicePlans()
	serviceKeysAPI := cfClient.ServiceKeys()
	serviceOfferringAPI := cfClient.ServiceOfferings()

	accClient, err := accountv2.NewClient(sess.BluemixSession)
	if err != nil {
		log.Fatal(err)
	}
	accountAPI := accClient.Accounts()

	clusterClient, err := k8sclusterv1.NewClient(sess.BluemixSession)
	if err != nil {
		log.Fatal(err)
	}
	clustersAPI := clusterClient.Clusters()
	clusterWorkerAPI := clusterClient.Workers()
	clusterSubnetsAPI := clusterClient.Subnets()
	clusterWebhookAPI := clusterClient.Webhooks()

	session := clientSession{
		session: sess,

		csClient:  clustersAPI,
		csSubnet:  clusterSubnetsAPI,
		csWorker:  clusterWorkerAPI,
		csWebHook: clusterWebhookAPI,

		cfOrgClient:              orgAPI,
		cfServiceInstanceClient:  serviceInstanceAPI,
		cfServiceKeysClient:      serviceKeysAPI,
		cfServicePlanClient:      servicePlanAPI,
		cfServiceOfferingsClient: serviceOfferringAPI,
		cfSpaceClient:            spaceAPI,
		bluemixAccountClient:     accountAPI,
	}

	return session, err
}
