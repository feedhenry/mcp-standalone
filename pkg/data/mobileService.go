package data

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

// MobileServiceValidator defines what a validator should do
type MobileServiceValidator interface {
	PreCreate(a *mobile.Service) error
	PreUpdate(old *mobile.Service, new *mobile.Service) error
}

// defaultSecretConvertor will provide a default secret to config conversion
type defaultSecretConvertor struct{}

func (dsc defaultSecretConvertor) Convert(s v1.Secret) (*mobile.ServiceConfig, error) {
	conf := map[string]string{}
	for k, v := range s.Data {
		conf[k] = string(v)
	}
	return &mobile.ServiceConfig{
		Config: conf,
		Name:   string(s.Data["name"]),
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

func (sa *secretAttributer) GetType() string {
	return strings.TrimSpace(string(sa.Secret.Data["type"]))
}

// MobileServiceRepo implments the mobile.ServiceCruder interface. it backed by the secret resource in kubernetes
type MobileServiceRepo struct {
	client     corev1.SecretInterface
	convertors map[string]SecretConvertor
	logger     *logrus.Logger
	validator  MobileServiceValidator
}

// NewMobileServiceRepo returns a new MobileServiceRepo
func NewMobileServiceRepo(client corev1.SecretInterface) *MobileServiceRepo {
	return &MobileServiceRepo{
		client: client,
		// if a secret needs a special convertor it is added here otherwise the default convertor will be used
		convertors: map[string]SecretConvertor{
			"keycloak": keycloakSecretConvertor{},
		},
		logger:    logrus.StandardLogger(),
		validator: DefaultMobileServiceValidator{},
	}
}

// Create will take a mobile service and create a secret to represent it
func (msr *MobileServiceRepo) Create(ms *mobile.Service) error {
	if err := msr.validator.PreCreate(ms); err != nil {
		return errors.Wrap(err, "create failed validation")
	}
	ms.ID = ms.Name + "-" + fmt.Sprintf("%v", time.Now().Unix())
	sct := convertMobileAppToSecret(*ms)
	if _, err := msr.client.Create(sct); err != nil {
		return errors.Wrap(err, "failed to create backing secret for mobile service")
	}
	return nil
}

func convertMobileAppToSecret(ms mobile.Service) *v1.Secret {
	data := map[string][]byte{}
	labels := map[string]string{
		"group": "mobile",
	}
	for k, v := range ms.Labels {
		labels[k] = v
	}
	data["uri"] = []byte(ms.Host)
	data["name"] = []byte(ms.Name)
	data["type"] = []byte(ms.Type)
	for k, v := range ms.Params {
		data[k] = []byte(v)
	}
	return &v1.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels: labels,
			Name:   ms.ID,
		},
		Data: data,
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
				msr.logger.Info("failed to find converter for ", svc.Name, " using default convertor")
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

// UpdateEnabledIntegrations will set labels on the underlying secret to indicate if an integration is enabled (it is really used as a que to the ui)
func (msr *MobileServiceRepo) UpdateEnabledIntegrations(svcName string, integrations map[string]string) error {
	secret, err := msr.client.Get(svcName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to read secret while updating enabled integrations")
	}
	if secret.Labels == nil {
		secret.Labels = map[string]string{}
	}
	for k, v := range integrations {
		secret.Labels[k] = v
	}
	_, err = msr.client.Update(secret)
	if err != nil {
		return errors.Wrap(err, "failed to update enabled integrations")
	}
	return nil
}

func convertSecretToMobileService(s v1.Secret) *mobile.Service {
	params := map[string]string{}
	for key, value := range s.Data {
		if key != "uri" && key != "name" {
			params[key] = string(value)
		}
	}
	external := s.Labels["external"] == "true"
	return &mobile.Service{
		ID:           s.Name,
		External:     external,
		Labels:       s.Labels,
		Name:         strings.TrimSpace(string(s.Data["name"])),
		Type:         strings.TrimSpace(string(s.Data["type"])),
		Host:         string(s.Data["uri"]),
		Params:       params,
		Integrations: map[string]*mobile.ServiceIntegration{},
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
