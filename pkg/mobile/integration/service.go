package integration

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
)

// MobileService holds the business logic for dealing with the mobile services and integrations with those services
type MobileService struct {
}

// DiscoverMobileServices will discover mobile services configured in the current namespace
func (ms *MobileService) DiscoverMobileServices(serviceCruder mobile.ServiceCruder) ([]*mobile.Service, error) {
	//todo move to config
	return ms.FindByNames([]string{"fh-sync-server", "keycloak"}, serviceCruder)
}

//FindByNames will return all services with a name that matches the provided name
func (ms *MobileService) FindByNames(names []string, serviceCruder mobile.ServiceCruder) ([]*mobile.Service, error) {
	filter := func(att mobile.Attributer) bool {
		for _, sn := range names {
			if sn == att.GetName() {
				return true
			}
		}
	}
	return nil, false
}

// GenerateMobileServiceConfigs will return a map of services and their mobile configs
func (ms *MobileService) GenerateMobileServiceConfigs(serviceCruder mobile.ServiceCruder) (map[string]*mobile.ServiceConfig, error) {
	svcConfigs, err := serviceCruder.ListConfigs(ms.filterServices)
	if err != nil {
		return nil, errors.Wrap(err, "GenerateMobileServiceConfigs failed during a list of configs")
	}
	configs := map[string]*mobile.ServiceConfig{}
	for _, sc := range svcConfigs {
		configs[sc.Name] = sc
	}
	return configs, nil
}

//MountSecretForComponent will work within namespace and mount secretName into componentName, so it can be configured to use serviceName, returning the modified deployment
func (ms *MobileService) MountSecretForComponent(k8s kubernetes.Interface, secretName, componentName, serviceName, namespace string) (*v1beta1.Deployment, error) {
	deploy, err := k8s.AppsV1beta1().Deployments(namespace).Get(componentName, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	id, err := findContainerID(componentName, deploy.Spec.Template.Spec.Containers)
	if err != nil {
		return nil, err
	}

	//only create the volume if it doesn't exist yet
	if vol := findVolumeByName(secretName, deploy.Spec.Template.Spec.Volumes); vol.Name != secretName {
		newVol := v1.Volume{
			Name: secretName,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		}
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, newVol)
	}

	if mount := findMountByName(secretName, deploy.Spec.Template.Spec.Containers[id].VolumeMounts); mount.Name != secretName {
		newMount := v1.VolumeMount{Name: secretName, ReadOnly: true, MountPath: "/etc/secrets/" + serviceName}
		deploy.Spec.Template.Spec.Containers[id].VolumeMounts = append(deploy.Spec.Template.Spec.Containers[id].VolumeMounts, newMount)
	}

	k8s.AppsV1beta1().Deployments(namespace).Update(deploy)

	return deploy, nil
}

func findContainerID(name string, containers []v1.Container) (int, error) {
	for id, container := range containers {
		if container.Name == name {
			return id, nil
		}
	}
	return -1, errors.New("could not find container with name: '" + name + "'")
}

func findVolumeByName(name string, volumes []v1.Volume) v1.Volume {
	for _, vol := range volumes {
		if vol.Name == name {
			return vol
		}
	}

	return v1.Volume{}
}

func findMountByName(name string, mounts []v1.VolumeMount) v1.VolumeMount {
	for _, mount := range mounts {
		if mount.Name == name {
			return mount
		}
	}

	return v1.VolumeMount{}
}
