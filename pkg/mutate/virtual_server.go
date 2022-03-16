package mutate

import (
	"encoding/json"
	"fmt"

	nginxv1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/configuration/v1"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutates Nginx VirtualServer
func MutateVirtualServer(body []byte) ([]byte, error) {

	// Unmarshal AdmissionReview
	admissionReview := v1.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("Unmarshalling failed: %s", err)
	}

	var err error
	var vserver *nginxv1.VirtualServer

	responseBody := []byte{}
	ar := admissionReview.Request
	resp := v1.AdmissionResponse{}

	if ar != nil {
		// Unmarshal the VirtualServer object
		if err := json.Unmarshal(ar.Object.Raw, &vserver); err != nil {
			return nil, fmt.Errorf("Unable to unmarshal virtualserver to object %v", err)
		}

		resp.Allowed = true
		resp.UID = ar.UID
		pT := v1.PatchTypeJSONPatch
		resp.PatchType = &pT

		namespace := vserver.Namespace

		resp.AuditAnnotations = map[string]string{
			"mutateme": "Change service name to imported name",
		}

		// Perform JSONPatch
		p := []map[string]string{}
		for i := range vserver.Spec.Upstreams {
			service := vserver.Spec.Upstreams[i].Service
			patch := map[string]string{
				"op":    "replace",
				"path":  fmt.Sprintf("/spec/upstreams/%d/service", i),
				"value": fmt.Sprintf(DerivedName(namespace, service)),
			}
			p = append(p, patch)
		}

		// Marshal the patch to JSON
		resp.Patch, err = json.Marshal(p)
		resp.Result = &metav1.Status{
			Status: "Success",
		}

		// Create response body to finish the AdmissionReview
		admissionReview.Response = &resp
		responseBody, err = json.Marshal(admissionReview)

		if err != nil {
			return nil, err // untested section
		}
	}

	return responseBody, nil
}
