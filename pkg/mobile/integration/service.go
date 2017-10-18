package integration

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pborman/uuid"
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
		"integrations": {mobile.ServiceNameKeycloak, mobile.IntegrationAPIKeys, mobile.ServiceNameThreeScale},
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
	svc, err := serviceCruder.List(filterServices(mobile.ValidServiceTypes))
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

func buildBindParams(from *mobile.Service, to *mobile.Service) (map[string]string, error) {
	var p map[string]string
	if from.Name == mobile.ServiceNameThreeScale {
		p = map[string]string{
			"apicast_route": from.Host,
			"service_route": to.Host,
			"service_name":  to.Name,
			"app_key":       uuid.New(),
		}
	} else if from.Name == mobile.ServiceNameKeycloak {
		p = map[string]string{
			"service_name": to.Name,
		}
	}

	return p, nil
}

// BindService will find the mobile service backed by a secret. It will use the values here to perform the binding
func (ms *MobileService) BindService(sccClient mobile.SCCInterface, svcCruder mobile.ServiceCruder, targetServiceName, service string) error {
	mobileService, err := svcCruder.Read(service)
	if err != nil {
		return errors.Wrap(err, "failed to read mobile service "+service)
	}
	targetService, err := svcCruder.Read(targetServiceName)
	if err != nil {
		return errors.Wrap(err, "failed to read target mobile service "+targetServiceName)
	}
	var namespace = ms.namespace
	if mobileService.Namespace != "" {
		namespace = mobileService.Namespace
	}
	bindParams, err := buildBindParams(mobileService, targetService)
	if err != nil {
		return errors.Wrap(err, "failed to build bind params for "+service)
	}
	if mobile.IntegrationAPIKeys == service {
		if err := sccClient.AddMobileApiKeys(targetServiceName, namespace); err != nil {
			return errors.Wrap(err, "failed to add mobile API Keys to service "+targetServiceName)
		}
	} else if err := sccClient.BindToService(mobileService.Name, targetService.Name, bindParams, namespace); err != nil {
		return errors.Wrap(err, "Binding "+service+" to "+targetServiceName+" failed")
	}
	if err := svcCruder.UpdateEnabledIntegrations(mobileService.ID, map[string]string{service: "true"}); err != nil {
		return errors.Wrap(err, "updating the enabled integrations for service "+targetServiceName+" failed ")
	}
	return nil
}

func (ms *MobileService) UnBindService(scClient mobile.SCCInterface, svcCruder mobile.ServiceCruder, targetServiceName, bindableService string) error {
	mobileService, err := svcCruder.Read(bindableService)
	if err != nil {
		return errors.Wrap(err, "failed to read mobile service "+bindableService)
	}
	var namespace = ms.namespace
	if mobileService.Namespace != "" {
		namespace = mobileService.Namespace
	}
	if mobile.IntegrationAPIKeys == mobileService.Name {
		if err := scClient.RemoveMobileApiKeys(targetServiceName, namespace); err != nil {
			return errors.Wrap(err, "failed to remove mobile API Keys from service "+targetServiceName)
		}
	} else if err := scClient.UnBindFromService(mobileService.Name, targetServiceName, namespace); err != nil {
		return errors.Wrap(err, "UnBinding Service from "+mobileService.Name+" failed")
	}
	if err := svcCruder.UpdateEnabledIntegrations(mobileService.ID, map[string]string{bindableService: "false"}); err != nil {
		return errors.Wrap(err, "updating the enabled integrations for service "+targetServiceName+" failed ")
	}
	return nil
}
