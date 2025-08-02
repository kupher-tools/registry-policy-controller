package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"gopkg.in/yaml.v2"
)

type RegistryPolicy struct {
	Common []string            `yaml:"commonRegistries"`
	PerNS  map[string][]string `yaml:"namespaceRegistryMap"`
}

func ValidateRegistry(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		http.Error(w, "cannot parse admission review", http.StatusBadRequest)
		return
	}

	kind := admissionReview.Request.Kind.Kind
	namespace := admissionReview.Request.Namespace
	var containers []corev1.Container
	switch kind {
	case "Pod":
		var pod corev1.Pod
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
			writeResponse(w, admissionReview.Request.UID, false, "Failed to parse Pod object")
			return
		}
		containers = pod.Spec.Containers
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &deployment); err != nil {
			writeResponse(w, admissionReview.Request.UID, false, "Failed to parse Deployment object")
			return
		}
		containers = deployment.Spec.Template.Spec.Containers

	case "DaemonSet":
		var daemonSet appsv1.DaemonSet
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &daemonSet); err != nil {
			writeResponse(w, admissionReview.Request.UID, false, "Failed to parse DaemonSet object")
			return
		}
		containers = daemonSet.Spec.Template.Spec.Containers
	case "StatefulSet":
		var statefulSet appsv1.StatefulSet
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &statefulSet); err != nil {
			writeResponse(w, admissionReview.Request.UID, false, "Failed to parse StatefulSet object")
			return
		}
		containers = statefulSet.Spec.Template.Spec.Containers
	default:
		log.Printf("Unsupported kind: %s", kind)
		writeResponse(w, admissionReview.Request.UID, false, "Unsupported kind: "+kind)
		return
	}

	if extractAndCheck(namespace, containers) {
		writeResponse(w, admissionReview.Request.UID, true, "")
	} else {
		writeResponse(w, admissionReview.Request.UID, false, "Container image is not allowed")
	}
}

func loadPolicyFromFiles() (*RegistryPolicy, error) {
	commonData, err := os.ReadFile("/etc/policies/commonRegistries.yaml")
	if err != nil {
		return nil, err
	}

	var common []string
	if err := yaml.Unmarshal(commonData, &common); err != nil {
		return nil, err
	}

	nsData, err := os.ReadFile("/etc/policies/namespaceRegistryMap.yaml")
	if err != nil {
		return nil, err
	}

	var perNS map[string][]string
	if err := yaml.Unmarshal(nsData, &perNS); err != nil {
		return nil, err
	}

	return &RegistryPolicy{
		Common: common,
		PerNS:  perNS,
	}, nil
}

var policyLoader = loadPolicyFromFiles

func extractAndCheck(namespace string, containers []corev1.Container) bool {
	policy, err := policyLoader()
	if err != nil {
		log.Printf("Error loading policy: %v", err)
		return false
	}
	if len(policy.Common) == 0 && len(policy.PerNS) == 0 {
		log.Println("No policies defined")
		return false
	}
	log.Println("Loaded policy:", policy)
	for _, container := range containers {
		// 1. Check common registries
		log.Printf("Checking container image: %s", container.Image)
		for _, allowedPrefix := range policy.Common {
			if strings.HasPrefix(container.Image, allowedPrefix) {
				return true
			}
		}
		// 2. Check namespace-specific registries
		if allowedNSPrefixes, ok := policy.PerNS[namespace]; ok {
			for _, nsPrefix := range allowedNSPrefixes {
				if strings.HasPrefix(container.Image, nsPrefix) {
					return true
				}
			}
		}
	}
	return false
}

func writeResponse(w http.ResponseWriter, uid types.UID, allowed bool, msg string) {
	response := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
		},
	}
	if !allowed {
		response.Response.Result = &metav1.Status{
			Message: msg,
		}
	}
	respBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
