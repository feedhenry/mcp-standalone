package data

import "github.com/feedhenry/mcp-standalone/pkg/mobile"

// DefaultMobileAppValidator validates mobile objects
type DefaultMobileAppValidator struct {
}

// PreCreate checks an App is valid before creating
func (mv DefaultMobileAppValidator) PreCreate(a *mobile.App) error {
	err := validateClientType(a)
	if err != nil {
		return err
	}
	if a.APIKey == "" {
		return &mobile.StatusError{Message: "apiKey cannot be empty"}
	}
	return nil
}

// PreUpdate checks that an update is valid before it is committed
func (mv DefaultMobileAppValidator) PreUpdate(old *mobile.App, new *mobile.App) error {
	err := validateClientType(new)
	if err != nil {
		return err
	}
	if new.Labels["group"] != "mobileapp" {
		return &mobile.StatusError{Message: "invalid action cannt update the group label"}
	}
	return nil
}

func validateClientType(a *mobile.App) error {
	if !mobile.ValidAppTypes.Contains(a.ClientType) {
		return &mobile.StatusError{
			Message: "invalid clientTypes " + a.ClientType + " valid client types " + mobile.ValidAppTypes.String(),
			Code:    400,
		}
	}
	return nil
}

type DefaultMobileServiceValidator struct{}

//TODO fill these out
func (msv DefaultMobileServiceValidator) PreCreate(ms *mobile.Service) error {
	return nil
}

func (msv DefaultMobileServiceValidator) PreUpdate(old *mobile.Service, new *mobile.Service) error {
	return nil
}

type DefaultMobileBuildValidator struct{}

//TODO fill these out
func (msv DefaultMobileBuildValidator) PreCreate(mb *mobile.Build) error {
	return nil
}

func (msv DefaultMobileBuildValidator) PreUpdate(old *mobile.Build, new *mobile.Build) error {
	return nil
}
