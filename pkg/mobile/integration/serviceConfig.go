package integration

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

//SDKService handles logic to do with the SDK
type SDKService struct{}

// GenerateMobileServiceConfigs will return a map of services and their mobile configs
func (ms *SDKService) GenerateMobileServiceConfigs(serviceCruder mobile.ServiceCruder) (map[string]*mobile.ServiceConfig, error) {
	svcConfigs, err := serviceCruder.ListConfigs(filterServices(mobile.ValidServiceTypes))
	if err != nil {
		return nil, errors.Wrap(err, "GenerateMobileServiceConfigs failed during a list of configs")
	}
	configs := map[string]*mobile.ServiceConfig{}
	for _, sc := range svcConfigs {
		configs[sc.Name] = sc
	}
	return configs, nil
}
