package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	pkgApi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesIngress() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesIngressCreate,
		Read:   resourceKubernetesIngressRead,
		Update: resourceKubernetesIngressUpdate,
		Exists: resourceKubernetesIngressExists,
		Delete: resourceKubernetesIngressDelete,
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
					Schema: ingressSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesIngressCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec := expandIngressSpec(d.Get("spec").([]interface{}))

	ingress := &v1beta1.Ingress{
		ObjectMeta: metadata,
		Spec:       spec,
	}
	out, err := conn.ExtensionsV1beta1().Ingresses(metadata.Namespace).Create(ingress)
	if err != nil {
		return err
	}

	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesIngressRead(d, meta)
}

func resourceKubernetesIngressRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Reading Ingress %s", name)
	ingress, err := conn.ExtensionsV1beta1().Ingresses(namespace).Get(name)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Received ingress: %#v", ingress)
	err = d.Set("metadata", flattenMetadata(ingress.ObjectMeta))
	if err != nil {
		return err
	}

	d.Set("spec", flattenIngressSpec(ingress.Spec))
	return nil
}

func resourceKubernetesIngressUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("spec") {
		specOps, err := patchIngressSpec("/spec", "spec", d)
		if err != nil {
			return err
		}
		ops = append(ops, specOps...)
	}

	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating ingress %q: %v", name, data)
	out, err := conn.ExtensionsV1beta1().Ingresses(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update ingress: %s", err)
	}

	log.Printf("[INFO] Submitting updated ingress: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesIngressRead(d, meta)

}

func resourceKubernetesIngressDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Deleting ingress: %#v", name)
	err := conn.ExtensionsV1beta1().Ingresses(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Ingress %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesIngressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Checking ingress %s", name)
	_, err := conn.ExtensionsV1beta1().Ingresses(namespace).Get(name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	return true, err
}
