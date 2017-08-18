package data

import "github.com/feedhenry/mobile-server/pkg/mobile"

// DefaultMobileAppValidator validates mobile objects
type DefaultMobileAppValidator struct {
}

// PreCreate checks an App is valid before creating
func (mv DefaultMobileAppValidator) PreCreate(a *mobile.App) error {
	return validateClientType(a)
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
		return &mobile.StatusError{Message: "invalid clientTypes " + a.ClientType + " valid client types " + mobile.ValidAppTypes.String()}
	}
	return nil
}
