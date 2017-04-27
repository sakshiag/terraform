package kubernetes

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	pkgApi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesReplicationController() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesReplicationControllerCreate,
		Read:   resourceKubernetesReplicationControllerRead,
		Update: resourceKubernetesReplicationControllerUpdate,
		Delete: resourceKubernetesReplicationControllerDelete,
		Exists: resourceKubernetesReplicationControllerExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("replication controller", true),
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Replicas is the number of desired replicas. This is a pointer to distinguish between explicit zero and unspecified. More info: http://kubernetes.io/docs/user-guide/replication-controller#what-is-a-replication-controller",
						},
						"min_ready_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)",
						},
						"selector": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Selector is a label query over pods that should match the Replicas count. If Selector is empty, it is defaulted to the labels present on the Pod template. Label keys and values that must match in order to be controlled by this replication controller, if empty defaulted to labels on Pod template. More info: http://kubernetes.io/docs/user-guide/labels#label-selectors",
						},
						"template": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected. This takes precedence over a TemplateRef. More info: http://kubernetes.io/docs/user-guide/replication-controller#pod-template",
							Elem: &schema.Resource{
								Schema: podTemplateSpecSchema(),
							},
						},
					},
				},
			},
			"delete_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"orphan_dependents": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Should the dependent objects be orphaned. If true/false, the orphan finalizer will be added to/removed from the object's finalizers list.",
						},
					},
				},
			},
		},
	}
}

func resourceKubernetesReplicationControllerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec := expandReplicationControllerSpec(d.Get("spec").([]interface{}))
	replicaCount := int(*spec.Replicas)
	rc := api.ReplicationController{
		ObjectMeta: metadata,
		Spec:       spec,
	}
	log.Printf("[INFO] Creating new replication controleer: %#v", rc)
	out, err := conn.CoreV1().ReplicationControllers(metadata.Namespace).Create(&rc)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Submitted new replication controller: %#v", out)

	stateConf := &resource.StateChangeConf{

		Target:  []string{"Running"},
		Pending: []string{"Pending"},
		Timeout: 5 * time.Minute,
		Refresh: func() (interface{}, string, error) {
			var statusPhase string
			out, err := conn.CoreV1().ReplicationControllers(metadata.Namespace).Get(metadata.Name)
			if err != nil {
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "Error", err
			}

			readyReplicas := out.Status.ReadyReplicas
			if int(readyReplicas) == replicaCount {
				statusPhase = "Running"
			} else {
				statusPhase = "Pending"
			}
			log.Printf("[DEBUG] Replication Controller %s status received: %#v", out.Name, statusPhase)
			return out, statusPhase, nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Replication Controller %s created", out.Name)

	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesReplicationControllerRead(d, meta)
}

func resourceKubernetesReplicationControllerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	log.Printf("[INFO] Reading replication controller %s", name)
	rc, err := conn.CoreV1().ReplicationControllers(namespace).Get(name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received replication controller: %#v", rc)
	err = d.Set("metadata", flattenMetadata(rc.ObjectMeta))
	if err != nil {
		return err
	}

	userSpec := expandReplicationControllerSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}
	err = d.Set("spec", flattenReplicationControllerSpec(rc.Spec, userSpec))
	if err != nil {
		return err
	}
	return nil
}

func resourceKubernetesReplicationControllerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("spec") {
		specOps, err := patchReplicationControllerSpec("/spec", "spec", d)
		if err != nil {
			return err
		}
		ops = append(ops, specOps...)
	}
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating  replication controller%s: %s", d.Id(), ops)
	out, err := conn.CoreV1().ReplicationControllers(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted updated replication controller: %#v", out)

	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesReplicationControllerRead(d, meta)
}

func resourceKubernetesReplicationControllerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	log.Printf("[INFO] Deleting replication controller: %#v", name)
	delete_options := expandDeleteOptions(d.Get("delete_options").([]interface{}))
	err := conn.CoreV1().ReplicationControllers(namespace).Delete(name, delete_options)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Replication Controller %s deleted", name)

	d.SetId("")
	return nil

}

func resourceKubernetesReplicationControllerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	log.Printf("[INFO] Checking replication controller %s", name)
	_, err := conn.CoreV1().ReplicationControllers(namespace).Get(name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
