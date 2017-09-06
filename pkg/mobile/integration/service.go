package integration

import (
	"fmt"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
)

// MobileService holds the business logic for dealing with the mobile services and integrations with those services
type MobileService struct {
	namespace string
}

func NewMobileSevice(ns string) *MobileService {
	return &MobileService{
		namespace: ns,
	}
}

//FindByNames will return all services with a name that matches the provided name
func (ms *MobileService) FindByNames(names []string, serviceCruder mobile.ServiceCruder) ([]*mobile.Service, error) {
	svc, err := serviceCruder.List(ms.filterServices(names))
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	return svc, nil
}

// TODO move to the secret data read when discovering the services
var capabilities = map[string]map[string][]string{
	"fh-sync-server": map[string][]string{
		"capabilities": {"data storage, data syncronisation"},
		"integrations": {"keycloak"},
	},
	"keycloak": map[string][]string{
		"capabilities": {"authentication, authorisation"},
		"integrations": {"fh-sync"},
	},
}

var serviceNames = []string{"fh-sync-server", "keycloak"}

// DiscoverMobileServices will discover mobile services configured in the current namespace
func (ms *MobileService) DiscoverMobileServices(serviceCruder mobile.ServiceCruder) ([]*mobile.Service, error) {
	//todo move to config

	svc, err := serviceCruder.List(ms.filterServices(serviceNames))
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	for _, s := range svc {
		s.Capabilities = capabilities[s.Name]
	}
	return svc, nil
}

// ReadMoileServiceAndIntegrations read servuce and any available service it can integrate with
func (ms *MobileService) ReadMoileServiceAndIntegrations(serviceCruder mobile.ServiceCruder, name string) (*mobile.Service, error) {
	//todo move to config
	svc, err := serviceCruder.Read(name)
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	svc.Capabilities = capabilities[svc.Name]
	if svc.Capabilities != nil {
		integrations := svc.Capabilities["integrations"]
		for _, v := range integrations {
			isvs, err := serviceCruder.List(ms.filterServices([]string{v}))
			if err != nil && !kerror.IsNotFound(err) {
				return nil, errors.Wrap(err, "failed attempting to discover mobile services.")
			}
			if len(isvs) != 0 {
				is := isvs[0]
				fmt.Println("svc label is ", is.Name, svc.Labels[is.Name])
				enabled := svc.Labels[is.Name] == "true"
				svc.Integrations[v] = &mobile.ServiceIntegration{
					ComponentSecret: svc.ID,
					Component:       svc.Name,
					Namespace:       ms.namespace,
					Service:         is.ID,
					Enabled:         enabled,
				}
			}
		}
	}
	return svc, nil
}

func (ms *MobileService) filterServices(serviceNames []string) func(att mobile.Attributer) bool {
	return func(att mobile.Attributer) bool {
		for _, sn := range serviceNames {
			if sn == att.GetName() {
				return true
			}
		}
		return false
	}
}

// GenerateMobileServiceConfigs will return a map of services and their mobile configs
func (ms *MobileService) GenerateMobileServiceConfigs(serviceCruder mobile.ServiceCruder) (map[string]*mobile.ServiceConfig, error) {
	svcConfigs, err := serviceCruder.ListConfigs(ms.filterServices(serviceNames))
	if err != nil {
		return nil, errors.Wrap(err, "GenerateMobileServiceConfigs failed during a list of configs")
	}
	configs := map[string]*mobile.ServiceConfig{}
	for _, sc := range svcConfigs {
		configs[sc.Name] = sc
	}
	return configs, nil
}

// TODO REFACTOR!!
//MountSecretForComponent will work within namespace and mount secretName into componentName, so it can be configured to use serviceName, returning the modified deployment
func (ms *MobileService) MountSecretForComponent(svcCruder mobile.ServiceCruder, k8s kubernetes.Interface, secretName, componentName, serviceName, namespace string, componentSecretName string) (*v1beta1.Deployment, error) {
	fmt.Println("mounting secret ", secretName, "into component ", componentName, serviceName)
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

	deploy, err = k8s.AppsV1beta1().Deployments(namespace).Update(deploy)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update deployment when mounting sercret for service integration with "+componentName)
	}
	//update secret with integration enabled
	bs, err := svcCruder.Read(serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read service secret")
	}
	enabled := map[string]string{bs.Name: "true"}
	if err := svcCruder.UpdateEnabledIntegrations(componentSecretName, enabled); err != nil {
		return nil, errors.Wrap(err, "failed to update enabled services after mounting secret")
	}
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
