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
	serviceNames := []string{"fh-sync-server", "keycloak"}
	filter := func(att mobile.Attributer) bool {
		for _, sn := range serviceNames {
			if sn == att.GetName() {
				return true
			}
		}
		return false
	}

	svc, err := serviceCruder.List(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Attempting to discover mobile services.")
	}
	return svc, nil
}
