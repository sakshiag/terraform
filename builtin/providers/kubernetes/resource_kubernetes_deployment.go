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
	v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesDeploymentCreate,
		Read:   resourceKubernetesDeploymentRead,
		Update: resourceKubernetesDeploymentUpdate,
		Exists: resourceKubernetesDeploymentExists,
		Delete: resourceKubernetesDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("deployment", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the pod owned by the cluster",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: deploymentSpecFileds(),
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
							Default:     false,
							Description: "Should the dependent objects be orphaned. If true/false, the orphan finalizer will be added to/removed from the object's finalizers list.",
						},
					},
				},
			},
		},
	}
}

func resourceKubernetesDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandDeploymentSpec(d.Get("spec").([]interface{}))
	replicaCount := int(*spec.Replicas)
	if err != nil {
		return err
	}
	deployment := &v1beta1.Deployment{
		ObjectMeta: metadata,
		Spec:       spec,
	}
	log.Printf("[INFO] Creating new deployment: %#v", deployment)
	out, err := conn.ExtensionsV1beta1().Deployments(metadata.Namespace).Create(deployment)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new deployment: %#v", out)
	stateConf := &resource.StateChangeConf{

		Target:  []string{"Running"},
		Pending: []string{"Pending"},
		Timeout: 5 * time.Minute,
		Refresh: func() (interface{}, string, error) {
			var statusPhase string
			out, err := conn.ExtensionsV1beta1().Deployments(metadata.Namespace).Get(metadata.Name)
			if err != nil {
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "Error", err
			}

			avaliableReplicas := out.Status.AvailableReplicas
			if int(avaliableReplicas) == replicaCount {
				statusPhase = "Running"
			} else {
				statusPhase = "Pending"
			}
			log.Printf("[DEBUG] Deployment %s status received: %#v", out.Name, statusPhase)
			return out, statusPhase, nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deployment %s created", out.Name)
	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesDeploymentRead(d, meta)
}

func resourceKubernetesDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Reading deployment %s", name)
	deployment, err := conn.ExtensionsV1beta1().Deployments(namespace).Get(name)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Received deployment: %#v", deployment)
	err = d.Set("metadata", flattenMetadata(deployment.ObjectMeta))
	if err != nil {
		return err
	}

	userSpec, err := expandDeploymentSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}

	d.Set("spec", flattenDeploymentSpec(deployment.Spec, userSpec))
	return nil
}

func resourceKubernetesDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("spec") {
		specOps, err := patchDeploymentSpec("/spec", "spec", d)
		if err != nil {
			return err
		}
		ops = append(ops, specOps...)
	}
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating  deployment%s: %s", d.Id(), ops)
	out, err := conn.ExtensionsV1beta1().Deployments(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted updated deployment: %#v", out)

	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesDeploymentRead(d, meta)
}

func resourceKubernetesDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Deleting Deployment: %q", name)
	m := make(map[string]string)
	deployment, _ := conn.ExtensionsV1beta1().Deployments(namespace).Get(name)
	m = deployment.Spec.Selector.MatchLabels
	delete_options := expandDeleteOptions(d.Get("delete_options").([]interface{}))
	err := conn.ExtensionsV1beta1().Deployments(namespace).Delete(name, delete_options)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deployment %s deleted", name)
	for k, v := range m {
		x := k + "=" + v
		fmt.Println(x)
		if !(*delete_options.OrphanDependents) {

			options := api.ListOptions{LabelSelector: x}
			rsList, err1 := conn.ExtensionsV1beta1().ReplicaSets(namespace).List(options)
			if err1 != nil {
				return err1
			}
			if rsList != nil {
				for i, _ := range rsList.Items {
					err2 := conn.ExtensionsV1beta1().ReplicaSets(namespace).Delete(rsList.Items[i].ObjectMeta.Name, delete_options)
					if err2 != nil {
						return err2
					}
					log.Printf("[INFO] Replica Set %s deleted", rsList.Items[i].ObjectMeta.Name)
				}
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceKubernetesDeploymentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Checking deployment %s", name)
	_, err := conn.ExtensionsV1beta1().Deployments(namespace).Get(name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	return true, err
}
