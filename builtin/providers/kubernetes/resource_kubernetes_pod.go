package kubernetes

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	pkgApi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesPod() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesPodCreate,
		Read:   resourceKubernetesPodRead,
		Update: resourceKubernetesPodUpdate,
		Delete: resourceKubernetesPodDelete,
		Exists: resourceKubernetesPodExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("Pod", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the pod owned by the cluster",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: podSpecFields(),
				},
			},
		},
	}
}
func resourceKubernetesPodCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandPodSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return err
	}

	pod := api.Pod{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	log.Printf("[INFO] Creating new pod: %#v", pod)
	out, err := conn.CoreV1().Pods(metadata.Namespace).Create(&pod)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new pod: %#v", out)

	d.SetId(buildId(out.ObjectMeta))

	stateConf := &resource.StateChangeConf{
		Target:  []string{"Running"},
		Pending: []string{"Pending"},
		Timeout: 5 * time.Minute,
		Refresh: func() (interface{}, string, error) {
			out, err := conn.CoreV1().Pods(metadata.Namespace).Get(metadata.Name)
			if err != nil {
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "Error", err
			}

			statusPhase := fmt.Sprintf("%v", out.Status.Phase)
			log.Printf("[DEBUG] Pods %s status received: %#v", out.Name, statusPhase)
			return out, statusPhase, nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Pod %s created", out.Name)

	return resourceKubernetesPodRead(d, meta)
}

func resourceKubernetesPodUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("spec") {
		specOps, err := patchPodSpec("/spec", "spec.0.", d)
		if err != nil {
			return err
		}
		ops = append(ops, specOps...)
	}
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating  pod%s: %s", d.Id(), ops)

	out, err := conn.CoreV1().Pods(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted updated pod: %#v", out)

	d.SetId(buildId(out.ObjectMeta))
	return resourceKubernetesPodRead(d, meta)
}

func resourceKubernetesPodRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Reading pod %s", name)
	pod, err := conn.CoreV1().Pods(namespace).Get(name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received pod: %#v", pod)

	err = d.Set("metadata", flattenMetadata(pod.ObjectMeta))
	if err != nil {
		return err
	}

	secretList, err := conn.CoreV1().Secrets(namespace).List(api.ListOptions{})
	if err != nil {
		return err
	}
	userVolumes := pickUserlVolumes(pod.Spec.Volumes, secretList, namespace)

	err = d.Set("spec", flattenPodSpec(pod.Spec, userVolumes))
	if err != nil {
		return err
	}
	return nil

}

func resourceKubernetesPodDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	namespace, name := idParts(d.Id())
	log.Printf("[INFO] Deleting pod: %#v", name)
	err := conn.CoreV1().Pods(namespace).Delete(name, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Pod %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesPodExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())
	log.Printf("[INFO] Checking pod %s", name)
	_, err := conn.CoreV1().Pods(namespace).Get(name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

//Return volumes which were created by user explicitly excluding the volumes created by k8s internally
func pickUserlVolumes(volumes []api.Volume, secretList *api.SecretList, namespace string) []api.Volume {
	internalVolumes := make(map[string]struct{})
	possiblyInternalVolumes := make([]string, 0, len(volumes))
	for _, v := range volumes {
		if v.Secret != nil && strings.HasPrefix(v.Name, "default-token-") {
			possiblyInternalVolumes = append(possiblyInternalVolumes, v.Name)
		}
	}
	for _, v := range possiblyInternalVolumes {
		for _, s := range secretList.Items {
			if s.Name != v {
				continue
			}
			for key, val := range s.Annotations {
				if key == "kubernetes.io/service-account.name" && val == "default" {
					//guarenteed internal volumes
					internalVolumes[v] = struct{}{}
				}
			}
		}
	}
	userVolumes := make([]api.Volume, 0, len(volumes)-len(internalVolumes))
	for _, v := range volumes {
		//Skip the volume which is internal to the k8s
		if _, ok := internalVolumes[v.Name]; ok {
			continue
		}
		userVolumes = append(userVolumes, v)
	}
	return userVolumes
}
