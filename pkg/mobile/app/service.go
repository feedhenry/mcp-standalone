package app

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Service struct {
}

func (s *Service) Create(appCrudder mobile.AppCruder, app *mobile.App) error {
	uid := uuid.NewV4()
	app.APIKey = uid.String()
	switch app.ClientType {
	case "android":
		app.MetaData["icon"] = "fa-android"
		break
	case "iOS":
		app.MetaData["icon"] = "fa-apple"
		break
	case "cordova":
		app.MetaData["icon"] = "icon-cordova"
		break
	}

	if err := appCrudder.Create(app); err != nil {
		err = errors.Wrap(err, "mobile app create: Attempted to create app via app repo")
		return err
	}

	if err := appCrudder.UpdateAppAPIKeys(app); err != nil {
		err = errors.Wrap(err, "app create, could not add api key")
		return err
	}

	return nil
}
