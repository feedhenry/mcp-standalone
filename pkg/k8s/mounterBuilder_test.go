package k8s

import (
	"testing"

	ktesting "k8s.io/client-go/testing"

	kerror "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

func listSecrets(a ktesting.Action) (bool, runtime.Object, error) {
	return true, &v1.SecretList{
		Items: []v1.Secret{
			{
				Data: map[string][]byte{
					"name": []byte("test-service"),
				},
			},
			{
				Data: map[string][]byte{
					"now":        []byte("something"),
					"completely": []byte("different"),
				},
			},
		},
	}, nil
}

func getSecrets(action ktesting.Action) (bool, runtime.Object, error) {
	return true, &v1.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "test-secret",
		},
		Data: map[string][]byte{},
	}, nil
}

func missingSecretFactory(secretName string) ktesting.ReactionFunc {
	return func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, &v1.Secret{}, kerror.NewNotFound(schema.GroupResource{Group: "", Resource: "secret"}, secretName)
	}

}

func updateDeployments(action ktesting.Action) (bool, runtime.Object, error) {
	return true, &v1beta1.Deployment{
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{},
					Containers: []v1.Container{
						{
							Name:         "fh-sync-server",
							VolumeMounts: []v1.VolumeMount{},
						},
					},
				},
			},
		},
	}, nil
}

func getUnmountedDeployments(action ktesting.Action) (bool, runtime.Object, error) {
	return true, &v1beta1.Deployment{
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{},
					Containers: []v1.Container{
						{
							Name:         "test-service",
							VolumeMounts: []v1.VolumeMount{},
						},
					},
				},
			},
		},
	}, nil
}

func getMountedDeployments(action ktesting.Action) (bool, runtime.Object, error) {
	return true, &v1beta1.Deployment{
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "test-secret",
						},
					},
					Containers: []v1.Container{
						{
							Name: "test-service",
							VolumeMounts: []v1.VolumeMount{
								{
									Name: "test-secret",
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func TestMount(t *testing.T) {
	cases := []struct {
		Name      string
		K8sClient func() *fake.Clientset
		Namespace string
		Service   string
		Secret    string
		Validate  func(t *testing.T, mountRes error)
	}{
		{
			Name: "Secret should be succesfully mounted",
			K8sClient: func() *fake.Clientset {
				k8sMock := &fake.Clientset{}
				k8sMock.AddReactor("list", "secrets", listSecrets)
				k8sMock.AddReactor("get", "secrets", getSecrets)
				k8sMock.AddReactor("Update", "deployments", updateDeployments)
				k8sMock.AddReactor("get", "deployments", getUnmountedDeployments)
				return k8sMock
			},
			Namespace: "test-namespace",
			Service:   "test-service",
			Secret:    "test-secret",
			Validate: func(t *testing.T, mountRes error) {
				if mountRes != nil {
					t.Fatalf("Did not expect error mounting in testCase: Valid mount request, got: %v", mountRes)
				}
			},
		},
		{
			Name: "Secret should not be mounted because a bad service name has been provided",
			K8sClient: func() *fake.Clientset {
				k8sMock := &fake.Clientset{}
				k8sMock.AddReactor("list", "secrets", listSecrets)
				k8sMock.AddReactor("get", "secrets", getSecrets)
				k8sMock.AddReactor("Update", "deployments", updateDeployments)
				k8sMock.AddReactor("get", "deployments", getUnmountedDeployments)
				return k8sMock
			},
			Namespace: "test-namespace",
			Service:   "test-bad-service",
			Secret:    "test-secret",
			Validate: func(t *testing.T, mountRes error) {
				if mountRes == nil {
					t.Fatalf("expected error when providing bad clientService name, but got none")
				}
			},
		},
		{
			Name: "Secret should not be mounted because a bad secret name has been provided",
			K8sClient: func() *fake.Clientset {
				k8sMock := &fake.Clientset{}
				k8sMock.AddReactor("list", "secrets", listSecrets)
				k8sMock.AddReactor("get", "secrets", missingSecretFactory("test-bad-secret"))
				k8sMock.AddReactor("Update", "deployments", updateDeployments)
				k8sMock.AddReactor("get", "deployments", getUnmountedDeployments)
				return k8sMock
			},
			Namespace: "test-namespace",
			Service:   "test-service",
			Secret:    "test-bad-secret",
			Validate: func(t *testing.T, mountRes error) {
				if mountRes == nil {
					t.Fatalf("expected error when providing bad secret name, but got none")
				}
			},
		},
	}

	for _, testCase := range cases {
		mb := NewMounterBuilder(testCase.Namespace).WithK8s(testCase.K8sClient()).Build()
		err := mb.Mount(testCase.Secret, testCase.Service)
		testCase.Validate(t, err)
	}
}

func TestUnmount(t *testing.T) {
	cases := []struct {
		Name      string
		K8sClient func() *fake.Clientset
		Namespace string
		Service   string
		Secret    string
		Validate  func(t *testing.T, unmountRes error)
	}{
		{
			Name: "Secret should be succesfully unmounted",
			K8sClient: func() *fake.Clientset {
				k8sMock := &fake.Clientset{}
				k8sMock.AddReactor("list", "secrets", listSecrets)
				k8sMock.AddReactor("get", "secrets", getSecrets)
				k8sMock.AddReactor("Update", "deployments", updateDeployments)
				k8sMock.AddReactor("get", "deployments", getMountedDeployments)
				return k8sMock
			},
			Namespace: "test-namespace",
			Service:   "test-service",
			Secret:    "test-secret",
			Validate: func(t *testing.T, unmountRes error) {
				if unmountRes != nil {
					t.Fatalf("Did not expect error mounting in testCase: Valid mount request, got: %v", unmountRes)
				}
			},
		},
		{
			Name: "Secret should not be unmounted because a bad service name has been provided",
			K8sClient: func() *fake.Clientset {
				k8sMock := &fake.Clientset{}
				k8sMock.AddReactor("list", "secrets", listSecrets)
				k8sMock.AddReactor("get", "secrets", getSecrets)
				k8sMock.AddReactor("Update", "deployments", updateDeployments)
				k8sMock.AddReactor("get", "deployments", getMountedDeployments)
				return k8sMock
			},
			Namespace: "test-namespace",
			Service:   "test-bad-service",
			Secret:    "test-secret",
			Validate: func(t *testing.T, unmountRes error) {
				if unmountRes == nil {
					t.Fatalf("expected error when providing bad clientService name, but got none")
				}
			},
		},
		{
			Name: "Secret should not be unmounted because a bad secret name has been provided",
			K8sClient: func() *fake.Clientset {
				k8sMock := &fake.Clientset{}
				k8sMock.AddReactor("list", "secrets", listSecrets)
				k8sMock.AddReactor("get", "secrets", missingSecretFactory("test-bad-secret"))
				k8sMock.AddReactor("Update", "deployments", updateDeployments)
				k8sMock.AddReactor("get", "deployments", getMountedDeployments)
				return k8sMock
			},
			Namespace: "test-namespace",
			Service:   "test-service",
			Secret:    "test-bad-secret",
			Validate: func(t *testing.T, unmountRes error) {
				if unmountRes == nil {
					t.Fatalf("expected error when providing bad secret name, but got none")
				}
			},
		},
	}

	for _, testCase := range cases {
		mb := NewMounterBuilder(testCase.Namespace).WithK8s(testCase.K8sClient()).Build()
		err := mb.Unmount(testCase.Secret, testCase.Service)
		testCase.Validate(t, err)
	}
}
