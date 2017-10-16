package integration

import (
	"fmt"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

// MobileService holds the business logic for dealing with the mobile services and integrations with those services
type MobileService struct {
	namespace string
}

//NewMobileSevice reutrns  a new mobile server
func NewMobileSevice(ns string) *MobileService {
	return &MobileService{
		namespace: ns,
	}
}

//FindByNames will return all services with a name that matches the provided name
func (ms *MobileService) FindByNames(names []string, serviceCruder mobile.ServiceCruder) ([]*mobile.Service, error) {
	svc, err := serviceCruder.List(filterServices(names))
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	return svc, nil
}

// TODO move to the secret data read when discovering the services
//TODO need to come up with a better way of representing this
var capabilities = map[string]map[string][]string{
	"fh-sync-server": map[string][]string{
		"capabilities": {"data storage, data syncronisation"},
		"integrations": {"keycloak", "mcp-mobile-keys", "3scale"},
	},
	"keycloak": map[string][]string{
		"capabilities": {"authentication, authorisation"},
		"integrations": {"fh-sync"},
	},
	"mcp-mobile-keys": map[string][]string{
		"capabilities": {"access apps"},
		"integrations": {},
	},
	"3scale": map[string][]string{
		"capabilities": {"authentication, authorization"},
		"integrations": {},
	},
	"custom": map[string][]string{
		"capabilities": {""},
		"integrations": {""},
	},
}

// DiscoverMobileServices will discover mobile services configured in the current namespace
func (ms *MobileService) DiscoverMobileServices(serviceCruder mobile.ServiceCruder, authChecker mobile.AuthChecker, client mobile.ExternalHTTPRequester) ([]*mobile.Service, error) {
	svc, err := serviceCruder.List(filterServices(mobile.ServiceTypes))
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	for _, s := range svc {
		s.Capabilities = capabilities[s.Name]
		//non external services are part of the current namespace //TODO maybe should be added to the apbs
		if s.External == false {
			if s.Namespace == "" {
				s.Namespace = ms.namespace
			}
			s.Writable = true
		}
		if s.External {
			perm, err := authChecker.Check("deployments", s.Namespace, client)
			if err != nil {
				return nil, errors.Wrap(err, "error checking access permissions")
			}
			s.Writable = perm
		}
	}
	return svc, nil
}

// ReadMobileServiceAndIntegrations read service and any available service it can integrate with
func (ms *MobileService) ReadMobileServiceAndIntegrations(serviceCruder mobile.ServiceCruder, authChecker mobile.AuthChecker, name string, client mobile.ExternalHTTPRequester) (*mobile.Service, error) {
	svc, err := serviceCruder.Read(name)
	if err != nil {
		return nil, errors.Wrap(err, "attempting to discover mobile services.")
	}
	svc.Capabilities = capabilities[svc.Type]
	if svc.Capabilities != nil {
		integrations := svc.Capabilities["integrations"]
		for _, v := range integrations {
			isvs, err := serviceCruder.List(filterServices([]string{v}))
			if err != nil {
				return nil, errors.Wrap(err, "failed attempting to discover mobile services.")
			}
			if len(isvs) > 0 {
				is := isvs[0]
				enabled := svc.Labels[is.Name] == "true"
				svc.Integrations[v] = &mobile.ServiceIntegration{
					ComponentSecret: svc.ID,
					Component:       svc.Type,
					DisplayName:     is.DisplayName,
					Namespace:       ms.namespace,
					Service:         is.ID,
					Enabled:         enabled,
				}
			}
		}
	}
	svc.Writable = true
	if svc.External {
		perm, err := authChecker.Check("deployments", svc.Namespace, client)
		if err != nil {
			return nil, errors.Wrap(err, "error checking access permissions")
		}
		svc.Writable = perm
	}
	return svc, nil
}

func filterServices(serviceTypes []string) func(att mobile.Attributer) bool {
	return func(att mobile.Attributer) bool {
		for _, sn := range serviceTypes {
			if sn == att.GetType() {
				return true
			}
		}
		return false
	}
}

func (ms *MobileService) BindService(sccClient mobile.SCCInterface, svcCruder mobile.ServiceCruder, targetServiceName, service string) error {
	if mobile.ServiceNameKeycloak == service {
		if err := sccClient.BindServiceToKeyCloak(targetServiceName, ms.namespace); err != nil {
			return errors.Wrap(err, "Binding Service to keycloak failed")
		}
		targetService, err := svcCruder.List(filterServices([]string{targetServiceName}))
		if err != nil || len(targetService) == 0 {
			return errors.New("failed to find client service: '" + targetServiceName + "'")
		}
		if err := svcCruder.UpdateEnabledIntegrations(targetService[0].ID, map[string]string{service: "true"}); err != nil {
			return errors.Wrap(err, "updating the enabled integrations for service "+targetServiceName+" failed ")
		}
		return nil
	}

	return errors.New("unknown service type " + service)
}

func (ms *MobileService) UnBindService(scClient mobile.SCCInterface, svcCruder mobile.ServiceCruder, targetServiceName, service string) error {
	if mobile.ServiceNameKeycloak == service {
		if err := scClient.UnBindServiceToKeyCloak(targetServiceName, ms.namespace); err != nil {
			return errors.Wrap(err, "UnBinding Service from keycloak failed")
		}
		targetService, err := svcCruder.List(filterServices([]string{targetServiceName}))
		if err != nil || len(targetService) == 0 {
			return errors.New("failed to find client service: '" + targetServiceName + "'")
		}
		if err := svcCruder.UpdateEnabledIntegrations(targetService[0].ID, map[string]string{service: "false"}); err != nil {
			return errors.Wrap(err, "updating the enabled integrations for service "+targetServiceName+" failed ")
		}
		return nil
	}
	return nil
}

//MountSecretForComponent will mount secret into component, returning any errors
func (ms *MobileService) MountSecretForComponent(svcCruder mobile.ServiceCruder, mounter mobile.VolumeMounter, clientServiceType, clientServiceName, serviceSecret string) error {
	//check secret exists and store for later update
	service, err := svcCruder.Read(serviceSecret)
	if err != nil {
		return errors.Wrap(err, "failed to find secret: '"+serviceSecret+"'")
	}

	css, err := svcCruder.List(filterServices([]string{clientServiceType}))
	if err != nil || len(css) == 0 {
		return errors.New("failed to find secret for client service: '" + clientServiceType + "'")
	}
	cService := &mobile.Service{}
	for _, cs := range css {
		fmt.Printf("cservice name: %s", cs.Name)
		if cs.Name == clientServiceName {
			cService = cs
		}
	}
	if cService.Name != clientServiceName {
		return errors.New("integration.ms.MountSecretForComponent -> Could not find service of type '" + clientServiceType + "' with name '" + clientServiceName + "'")
	}

	err = mounter.Mount(service, cService)
	if err != nil {
		return errors.Wrap(err, "failed to mount secret '"+serviceSecret+"' into service '"+clientServiceType+"'")
	}

	clientServiceID := cService.ID

	//update secret with integration enabled
	enabled := map[string]string{service.Type: "true"}
	if err := svcCruder.UpdateEnabledIntegrations(clientServiceID, enabled); err != nil {
		return errors.Wrap(err, "failed to update enabled services after mounting secret")
	}

	return nil
}

//UnmountSecretInComponent will unmount secret from component, so it can be no longer use serviceName, returning any errors
func (ms *MobileService) UnmountSecretInComponent(svcCruder mobile.ServiceCruder, unmounter mobile.VolumeUnmounter, clientServiceType, clientServiceName, serviceSecret string) error {
	//check secret exists and store for later update
	service, err := svcCruder.Read(serviceSecret)
	if err != nil {
		return errors.Wrap(err, "failed to find secret: '"+serviceSecret+"'")
	}

	//find the clientService secret name
	css, err := svcCruder.List(filterServices([]string{clientServiceType}))
	if err != nil || len(css) == 0 {
		return errors.New("failed to find secret for client service: '" + clientServiceType + "'")
	}
	cService := &mobile.Service{}
	for _, cs := range css {
		if cs.Name == clientServiceName {
			cService = cs
		}
	}
	if cService.Name != clientServiceName {
		return errors.New("integration.ms.UnmountSecretForComponent -> Could not find service of type '" + clientServiceType + "' with name '" + clientServiceName + "'")
	}

	err = unmounter.Unmount(service, cService)
	if err != nil {
		return errors.Wrap(err, "failed to unmount secret '"+serviceSecret+"' from component '"+clientServiceType+"'")
	}

	clientServiceId := cService.ID

	//update secret with integration enabled
	disabled := map[string]string{service.Type: "false"}
	if err := svcCruder.UpdateEnabledIntegrations(clientServiceId, disabled); err != nil {
		return errors.Wrap(err, "failed to update enabled services after unmounting secret")
	}

	return nil
}
