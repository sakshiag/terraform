package brooklyn

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/brooklyncentral/brooklyn-cli/net"
	"github.com/brooklyncentral/brooklyn-cli/api/application"
	"os"
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/resource"
	"time"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Delete: resourceApplicationDelete,

		Schema: map[string]*schema.Schema{
			"application_spec": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateYamlConfigFile,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"locations": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:  &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"links": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:  &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	network := meta.(*net.Network)

	applicationSpec := d.Get("application_spec")

	log.Printf("[DEBUG] Submit brooklyn blueprint: %s", applicationSpec)

	// Create the application
	task, err := application.Create(network, applicationSpec.(string))
	if err != nil {
		return fmt.Errorf("Unable to create application: %s", err)
	}

	// Store the resulting application ID
	d.SetId(task.EntityId)

	// Wait for the application to become running
	log.Printf("[DEBUG] Waiting for instance (%s) to become running", task.EntityId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"STARTING"},
		Target:     "RUNNING",
		Refresh:    ApplicationStateRefreshFunc(network, task.EntityId),
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", task.EntityId, err)
	}

	return resourceApplicationRead(d, meta)
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	network := meta.(*net.Network)

	applicationSummary, err := application.Application(network, d.Id())
	if err != nil {
		return fmt.Errorf("Unable to read application: %s", err)
	}

	d.Set("name", applicationSummary.Spec.Name)
	d.Set("locations", applicationSummary.Spec.Locations)
	d.Set("type", applicationSummary.Spec.Type)
	d.Set("status", applicationSummary.Status)
	d.Set("links", applicationSummary.Links)

	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	network := meta.(*net.Network)

	applicationId := d.Id()

	log.Printf("[DEBUG] Delete application (%s)", applicationId)

	if _, err := application.Delete(network, applicationId); err != nil {
		return fmt.Errorf("Unable to delete application: %s", err)
	}

	return nil
}

func validateYamlConfigFile(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if _, err := os.Stat(value); os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("YAML file %q cannot be found: %q", k, value))
	}
	return
}

// ApplicationStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch an application status.
func ApplicationStateRefreshFunc(network *net.Network, applicationID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := application.Application(network, applicationID)

		if err != nil {
			log.Printf("Error on ApplicationStateRefresh: %s", err)
			return nil, "", err
		}

		if &resp == nil {
			return nil, "", nil
		}

		return resp, string(resp.Status), nil
	}
}