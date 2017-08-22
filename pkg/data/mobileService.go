package data

import (
	"strings"

	"github.com/feedhenry/mobile-server/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type MobileServiceRepo struct {
	client corev1.SecretInterface
}

type secretAttributer struct {
	*v1.Secret
}

func (sa *secretAttributer) GetName() string {
	return strings.TrimSpace(string(sa.Secret.Data["name"]))
}

func NewMobileServiceRepo(client corev1.SecretInterface) *MobileServiceRepo {
	return &MobileServiceRepo{
		client: client,
	}
}

func (msr *MobileServiceRepo) List(f mobile.AttrFilterFunc) ([]*mobile.Service, error) {
	svs, err := msr.client.List(meta_v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list secrets in namespace")
	}
	ret := []*mobile.Service{}
	for _, item := range svs.Items {
		if f(&secretAttributer{&item}) {
			ret = append(ret, convertSecretToMobileService(item))
		}
	}
	return ret, nil
}

func convertSecretToMobileService(s v1.Secret) *mobile.Service {
	return &mobile.Service{
		Name:   strings.TrimSpace(string(s.Data["name"])),
		Host:   string(s.Data["uri"]),
		Params: map[string]string{},
	}
}

func NewServiceRepoBuilder() mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{}
}

// MobileServiceRepoBuilder builds a ServiceCruder
type MobileServiceRepoBuilder struct {
	client corev1.SecretInterface
}

// WithClient sets the client to use
func (marb *MobileServiceRepoBuilder) WithClient(c corev1.SecretInterface) mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{
		client: c,
	}
}

// Build builds the final repo
func (marb *MobileServiceRepoBuilder) Build() mobile.ServiceCruder {
	return NewMobileServiceRepo(marb.client)
}
