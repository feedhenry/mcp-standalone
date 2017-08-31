package integration

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

// MobileService holds the business logic for dealing with the mobile services and integrations with those services
type MobileService struct {
}

// DiscoverMobileServices will discover mobile services configured in the current namespace
func (ms *MobileService) DiscoverMobileServices(serviceCruder mobile.ServiceCruder) ([]*mobile.Service, error) {
	//todo move to config
	svc, err := serviceCruder.List(ms.filterServices)
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	return svc, nil
}

func (ms *MobileService) filterServices(att mobile.Attributer) bool {
	var serviceNames = []string{"fh-sync-server", "keycloak"}
	for _, sn := range serviceNames {
		if sn == att.GetName() {
			return true
		}
	}
	return false
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
