package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	pkgApi "k8s.io/kubernetes/pkg/api"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesServiceCreate,
		Read:   resourceKubernetesServiceRead,
		Update: resourceKubernetesServiceUpdate,
		Delete: resourceKubernetesServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("service", true),
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: serviceSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	serviceSpec, err := expandServiceSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}

	service := &api.Service{
		ObjectMeta: metadata,
		Spec:       serviceSpec,
	}

	out, err := conn.Services(metadata.Namespace).Create(service)
	if err != nil {
		return err
	}

	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesServiceRead(d, meta)
}

func resourceKubernetesServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Reading service %s", name)
	service, err := conn.Services(namespace).Get(name)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Received service: %#v", service)
	err = d.Set("metadata", flattenMetadata(service.ObjectMeta))
	if err != nil {
		return err
	}

	d.Set("spec", flattenServiceSpec(service.Spec))
	return nil
}

func resourceKubernetesServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("spec") {
		specOps, err := patchServiceSpec("/spec", "spec", d)
		if err != nil {
			return err
		}
		ops = append(ops, specOps...)
	}

	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating service %q: %v", name, data)
	out, err := conn.Services(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update service: %s", err)
	}

	log.Printf("[INFO] Submitting updated service: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesServiceRead(d, meta)

}

func resourceKubernetesServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Deleting service: %#v", name)
	err := conn.Services(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Service %s deleted", name)

	d.SetId("")
	return nil

}
