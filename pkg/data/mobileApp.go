package data

import (
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
)

// MobileAppValidator defines what a validator should do
type MobileAppValidator interface {
	PreCreate(a *mobile.App) error
	PreUpdate(old *mobile.App, new *mobile.App) error
}

// MobileAppRepo interacts with the data store that backs the mobile objects
type MobileAppRepo struct {
	client    corev1.ConfigMapInterface
	validator MobileAppValidator
}

// NewMobileAppRepo instansiates a new MobileAppRepo
func NewMobileAppRepo(c corev1.ConfigMapInterface, v MobileAppValidator) *MobileAppRepo {
	rep := &MobileAppRepo{
		client:    c,
		validator: v,
	}
	if rep.validator == nil {
		rep.validator = &DefaultMobileAppValidator{}
	}
	return rep
}

// ReadByName attempts to read a mobile app by its unique name
func (mar *MobileAppRepo) ReadByName(name string) (*mobile.App, error) {
	_, cm, err := mar.readMobileAppAndConfigMap(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve mobile app ")
	}
	return convertConfigMapToMobileApp(cm), nil
}

// Create creates a mobile app object. Fails on duplicates
func (mar *MobileAppRepo) Create(app *mobile.App) error {
	if err := mar.validator.PreCreate(app); err != nil {
		return errors.Wrap(err, "validation failed during create")
	}
	cm := v1.ConfigMap{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: app.Name,
			Labels: map[string]string{
				"group": "mobileapp",
			},
		},

		Data: map[string]string{
			"name":       app.Name,
			"clientType": app.ClientType,
			"apiKey":     app.APIKey,
		},
	}
	if _, err := mar.client.Create(&cm); err != nil {
		return errors.Wrap(err, "failed to create underlying configmap for mobile app")
	}
	return nil
}

//DeleteByName will delte the underlying configmap
func (mar *MobileAppRepo) DeleteByName(name string) error {
	return mar.client.Delete(name, &meta_v1.DeleteOptions{})
}

//List will list the configmaps and convert them to mobileapps
func (mar *MobileAppRepo) List() ([]*mobile.App, error) {
	list, err := mar.client.List(meta_v1.ListOptions{LabelSelector: "group=mobileapp"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list mobileapp configmaps")
	}

	var apps = []*mobile.App{}
	for _, a := range list.Items {
		apps = append(apps, convertConfigMapToMobileApp(&a))
	}
	return apps, nil
}

// Update will update the underlying configmap with the new details
func (mar *MobileAppRepo) Update(app *mobile.App) (*mobile.App, error) {
	old, cm, err := mar.readMobileAppAndConfigMap(app.Name)
	if err != nil {
		return nil, err
	}
	if err := mar.validator.PreUpdate(old, app); err != nil {
		return nil, errors.Wrap(err, "validation failed before update")
	}
	cm.Data["name"] = app.Name
	cm.Data["clientType"] = app.ClientType
	var cmap *v1.ConfigMap
	cmap, err = mar.client.Update(cm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update mobile app configmap")
	}
	return convertConfigMapToMobileApp(cmap), nil
}

func convertConfigMapToMobileApp(m *v1.ConfigMap) *mobile.App {
	return &mobile.App{
		Name:       m.Data["name"],
		ClientType: m.Data["clientType"],
		APIKey:     m.Data["apiKey"],
		Labels:     m.Labels,
	}
}

func convertMobileAppToConfigMap(app *mobile.App) *v1.ConfigMap {
	return &v1.ConfigMap{
		Data: map[string]string{
			"name":       app.Name,
			"clientType": app.ClientType,
			"apiKey":     app.APIKey,
		},
	}
}

func (mar *MobileAppRepo) readMobileAppAndConfigMap(name string) (*mobile.App, *v1.ConfigMap, error) {
	cm, err := mar.readUnderlyingConfigMap(&mobile.App{Name: name})
	if err != nil {
		return nil, nil, err
	}
	app := convertConfigMapToMobileApp(cm)
	return app, cm, err
}

func (mar *MobileAppRepo) readUnderlyingConfigMap(a *mobile.App) (*v1.ConfigMap, error) {
	cm, err := mar.client.Get(a.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read underlying configmap for app "+a.Name)
	}
	return cm, nil
}

func NewMobileAppRepoBuilder() mobile.AppRepoBuilder {
	return &MobileAppRepoBuilder{}
}

// MobileAppRepoBuilder builds a MobileAppRepo
type MobileAppRepoBuilder struct {
	client corev1.ConfigMapInterface
}

// WithClient sets the client to use
func (marb *MobileAppRepoBuilder) WithClient(c corev1.ConfigMapInterface) mobile.AppRepoBuilder {
	return &MobileAppRepoBuilder{
		client: c,
	}
}

// Build builds the final repo
func (marb *MobileAppRepoBuilder) Build() mobile.AppCruder {
	return NewMobileAppRepo(marb.client, DefaultMobileAppValidator{})
}
