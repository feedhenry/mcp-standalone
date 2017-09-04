package data

import (
	"encoding/json"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type SecretConvertor interface {
	Convert(s v1.Secret) (*mobile.ServiceConfig, error)
}

// defaultSecretConvertor will provide a default secret to config conversion
type defaultSecretConvertor struct{}

func (dsc defaultSecretConvertor) Convert(s v1.Secret) (*mobile.ServiceConfig, error) {
	return &mobile.ServiceConfig{
		Config: map[string]string{
			"uri":  string(s.Data["uri"]),
			"name": string(s.Data["name"]),
		},
		Name: string(s.Data["name"]),
	}, nil
}

type keycloakSecretConvertor struct{}

func (ksc keycloakSecretConvertor) Convert(s v1.Secret) (*mobile.ServiceConfig, error) {
	kc := &mobile.KeycloakConfig{}
	err := json.Unmarshal(s.Data["installation"], kc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshalal keycloak configuration ")
	}
	return &mobile.ServiceConfig{
		Config: kc,
		Name:   string(s.Data["name"]),
	}, nil
}

type secretAttributer struct {
	*v1.Secret
}

func (sa *secretAttributer) GetName() string {
	var name = strings.TrimSpace(string(sa.Secret.Data["name"]))
	if "" == name {
		//remove once we fix keycloak apb
		name = strings.TrimSpace(string(sa.Secret.Data["NAME"]))
	}
	return name
}

// MobileServiceRepo implments the mobile.ServiceCruder interface. it backed by the secret resource in kubernetes
type MobileServiceRepo struct {
	client     corev1.SecretInterface
	convertors map[string]SecretConvertor
	logger     *logrus.Logger
}

// NewMobileServiceRepo returns a new MobileServiceRepo
func NewMobileServiceRepo(client corev1.SecretInterface) *MobileServiceRepo {
	return &MobileServiceRepo{
		client: client,
		// if a secret needs a special convertor it is added here otherwise the default convertor will be used
		convertors: map[string]SecretConvertor{
			"keycloak": keycloakSecretConvertor{},
		},
		logger: logrus.StandardLogger(),
	}
}

// List will read all the secrets in the namespace and filter them based on the passed function
func (msr *MobileServiceRepo) List(filter mobile.AttrFilterFunc) ([]*mobile.Service, error) {
	svs, err := msr.client.List(meta_v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list secrets in namespace")
	}
	ret := []*mobile.Service{}
	for _, item := range svs.Items {
		if filter(&secretAttributer{&item}) {
			ret = append(ret, convertSecretToMobileService(item))
		}
	}
	return ret, nil
}

// Read the mobile service
func (msr *MobileServiceRepo) Read(name string) (*mobile.Service, error) {
	svc, err := msr.client.Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get backing secret in repo Read ")
	}
	return convertSecretToMobileService(*svc), nil
}

// ListConfigs will build a list of configs based on the available services that are represented by secrets in the namespace
func (msr *MobileServiceRepo) ListConfigs(filter mobile.AttrFilterFunc) ([]*mobile.ServiceConfig, error) {
	svs, err := msr.client.List(meta_v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list secrets in namespace")
	}
	ret := []*mobile.ServiceConfig{}
	for _, item := range svs.Items {
		if filter(&secretAttributer{&item}) {
			var svcConifg *mobile.ServiceConfig
			svc := convertSecretToMobileService(item)
			if _, ok := msr.convertors[svc.Name]; !ok {
				msr.logger.Info("failed to find converter for ", svc.Name, "using default convertor")
				convertor := defaultSecretConvertor{}
				svcConifg, err = convertor.Convert(item)
				if err != nil {
					//bail out here as now our config may not be compelete?
					return nil, errors.Wrap(err, "failed to convert config for service: "+svc.Name)
				}
			} else {
				// we can only convert what is available
				convertor := msr.convertors[svc.Name]
				svcConifg, err = convertor.Convert(item)
				if err != nil {
					//bail out here as now our config may not be compelete?
					return nil, errors.Wrap(err, "failed to convert config for service: "+svc.Name)
				}
			}
			ret = append(ret, svcConifg)
		}
	}
	return ret, nil
}

func convertSecretToMobileService(s v1.Secret) *mobile.Service {
	params := map[string]string{}
	for key, value := range s.Data {
		if key != "uri" && key != "name" {
			params[key] = string(value)
		}
	}
	return &mobile.Service{
		ID:                s.Name,
		Name:              strings.TrimSpace(string(s.Data["name"])),
		Host:              string(s.Data["uri"]),
		BindingSecretName: s.GetName(),
		Params:            params,
		Integrations:      map[string]*mobile.ServiceIntegration{},
	}
}

// NewServiceRepoBuilder provides an implementation of mobile.ServiceRepoBuilder
func NewServiceRepoBuilder() mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{}
}

// MobileServiceRepoBuilder builds a ServiceCruder
type MobileServiceRepoBuilder struct {
	client corev1.SecretInterface
}

// WithClient sets the client to use
func (marb *MobileServiceRepoBuilder) WithClient(client corev1.SecretInterface) mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{
		client: client,
	}
}

// Build builds the final repo
func (marb *MobileServiceRepoBuilder) Build() mobile.ServiceCruder {
	return NewMobileServiceRepo(marb.client)
}
