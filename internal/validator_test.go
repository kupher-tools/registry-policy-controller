package internal

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestExtractAndCheck_CommonRegistryAllowed(t *testing.T) {
	policyLoader = func() (*RegistryPolicy, error) {
		return &RegistryPolicy{
			Common: []string{"registry.kupher.io/", "ghcr.io/your-org/"},
			PerNS:  map[string][]string{},
		}, nil
	}

	containers := []corev1.Container{
		{Name: "app", Image: "ghcr.io/your-org/app:v1"},
	}

	if !extractAndCheck("any-namespace", containers) {
		t.Error("expected image to be allowed via common registry")
	}
}

func TestExtractAndCheck_NamespaceRegistryAllowed(t *testing.T) {
	policyLoader = func() (*RegistryPolicy, error) {
		return &RegistryPolicy{
			Common: []string{},
			PerNS: map[string][]string{
				"devops": {"ghcr.io/devops-team/"},
			},
		}, nil
	}

	containers := []corev1.Container{
		{Name: "svc", Image: "ghcr.io/devops-team/service:v2"},
	}

	if !extractAndCheck("devops", containers) {
		t.Error("expected image to be allowed via namespace registry")
	}
}

func TestExtractAndCheck_DisallowedImage(t *testing.T) {
	policyLoader = func() (*RegistryPolicy, error) {
		return &RegistryPolicy{
			Common: []string{"registry.kupher.io/"},
			PerNS:  map[string][]string{},
		}, nil
	}

	containers := []corev1.Container{
		{Name: "hacker", Image: "docker.io/library/nginx:latest"},
	}

	if extractAndCheck("default", containers) {
		t.Error("expected image to be denied")
	}
}

func TestExtractAndCheck_EmptyPolicyDeniesAll(t *testing.T) {
	policyLoader = func() (*RegistryPolicy, error) {
		return &RegistryPolicy{}, nil
	}

	containers := []corev1.Container{
		{Name: "nginx", Image: "docker.io/library/nginx:latest"},
	}

	if extractAndCheck("default", containers) {
		t.Error("expected image to be denied due to empty policy")
	}
}

func TestExtractAndCheck_MultipleContainers_OneAllowed(t *testing.T) {
	policyLoader = func() (*RegistryPolicy, error) {
		return &RegistryPolicy{
			Common: []string{"ghcr.io/your-org/"},
		}, nil
	}

	containers := []corev1.Container{
		{Name: "safe", Image: "ghcr.io/your-org/app:v1"},
		{Name: "unsafe", Image: "docker.io/library/nginx"},
	}

	if !extractAndCheck("any", containers) {
		t.Error("expected image to be allowed because one container is valid")
	}
}
