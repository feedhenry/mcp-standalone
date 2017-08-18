package data

import (
	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type MobileAppRepo struct {
	client corev1.ConfigMapInterface
}

func NewMobileAppRepo(c corev1.ConfigMapInterface) *MobileAppRepo {
	return &MobileAppRepo{
		client: c,
	}
}
func (mar *MobileAppRepo) ReadByName(name string) (*mobile.App, error) {
	cm, err := mar.client.Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve mobile app ")
	}
	return convertConfigMapToMobileApp(cm), nil
}

func (mar *MobileAppRepo) Create(app *mobile.App) error {
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
	list, err := mar.client.List(meta_v1.ListOptions{LabelSelector: "mobileapp"})
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
	return nil, nil
}

func convertConfigMapToMobileApp(m *v1.ConfigMap) *mobile.App {
	return &mobile.App{
		Name:       m.Data["name"],
		ClientType: m.Data["clientType"],
	}
}

type MobileAppRepoBuilder struct {
	client corev1.ConfigMapInterface
}

func (marb *MobileAppRepoBuilder) WithClient(c corev1.ConfigMapInterface) *MobileAppRepoBuilder {
	return &MobileAppRepoBuilder{
		client: c,
	}
}

func (marb *MobileAppRepoBuilder) Build() *MobileAppRepo {
	return NewMobileAppRepo(marb.client)
}
