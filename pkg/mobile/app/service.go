package app

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Service struct {
}

func (s *Service) Create(appCrudder mobile.AppCruder, apiKeyMapCruder mobile.ServiceCruder, app *mobile.App) error {
	uid := uuid.NewV4()
	app.APIKey = uid.String()
	switch app.ClientType {
	case mobile.AndroidApp:
		app.MetaData["icon"] = "fa-android"
		break
	case mobile.IOSApp:
		app.MetaData["icon"] = "fa-apple"
		break
	case mobile.CordovaApp:
		app.MetaData["icon"] = "icon-cordova"
		break
	}

	if err := appCrudder.Create(app); err != nil {
		err = errors.Wrap(err, "mobile app create: Attempted to create app via app repo")
		return err
	}

	if err := apiKeyMapCruder.UpdateAPIKeyMap(app.ID, app.APIKey); err != nil {
		err = errors.Wrap(err, "app create, could not add api key")
		return err
	}

	return nil
}

func (s *Service) Delete(appCruder mobile.AppCruder, apiKeyMapCruder mobile.ServiceCruder, appID string) error {
	if err := appCruder.DeleteByName(appID); err != nil {
		return err
	}
	if err := apiKeyMapCruder.RemoveAPIKeyByID(appID); err != nil {
		return err
	}
	return nil
}
