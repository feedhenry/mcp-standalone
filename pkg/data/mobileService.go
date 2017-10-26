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

// SecretConvertor converts a kubernetes secret into a mobile.ServiceConfig
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

//Convert a kubernetes secret to a mobile.ServiceConfig
func (dsc defaultSecretConvertor) Convert(s v1.Secret) (*mobile.ServiceConfig, error) {
	config := map[string]interface{}{}
	headers := map[string]string{}
	for k, v := range s.Data {
		config[k] = string(v)
	}
	config["headers"] = headers
	return &mobile.ServiceConfig{
		Config: config,
		Name:   string(s.Data["name"]),
	}, nil
}

type keycloakSecretConvertor struct{}

//Convert a kubernetes keycloak secret into a keycloak mobile.ServiceConfig
func (ksc keycloakSecretConvertor) Convert(s v1.Secret) (*mobile.ServiceConfig, error) {
	config := map[string]interface{}{}
	headers := map[string]string{}
	err := json.Unmarshal(s.Data["public_installation"], &config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall keycloak configuration ")
	}
	config["headers"] = headers
	return &mobile.ServiceConfig{
		Config: config,
		Name:   string(s.Data["name"]),
	}, nil
}

type syncSecretConvertor struct{}

//Convert a kubernetes Sync Server secret into a keycloak mobile.ServiceConfig
func (scc syncSecretConvertor) Convert(s v1.Secret) (*mobile.ServiceConfig, error) {
	config := map[string]interface{}{}
	headers := map[string]string{}

	acAppID, acAppIDExists := s.Data["apicast_app_id"]
	acAppKey, acAppKeyExists := s.Data["apicast_app_key"]
	acRoute, acRouteExists := s.Data["apicast_route"]
	if acAppIDExists && acAppKeyExists && acRouteExists {
		headers["app_id"] = string(acAppID)
		headers["app_key"] = string(acAppKey)
		config["uri"] = string(acRoute)
	}
	config["headers"] = headers

	return &mobile.ServiceConfig{
		Config: config,
		Name:   string(s.Data["name"]),
	}, nil
}

type secretAttributer struct {
	*v1.Secret
}

//GetName returns the value of the name field in the secret
func (sa *secretAttributer) GetName() string {
	var name = strings.TrimSpace(string(sa.Secret.Data["name"]))
	return name
}

//GetType returns the value of the type field in the secret
func (sa *secretAttributer) GetType() string {
	return strings.TrimSpace(string(sa.Secret.Data["type"]))
}

// MobileServiceRepo implements the mobile.ServiceCruder interface. it backed by the secret resource in kubernetes
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
			"keycloak":       keycloakSecretConvertor{},
			"fh-sync-server": syncSecretConvertor{},
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
	if ms.DisplayName == "" {
		ms.DisplayName = ms.Name
	}
	sct := convertMobileAppToSecret(*ms)
	if _, err := msr.client.Create(sct); err != nil {
		return errors.Wrap(err, "failed to create backing secret for mobile service")
	}
	return nil
}

func convertMobileAppToSecret(ms mobile.Service) *v1.Secret {
	data := map[string][]byte{}
	labels := map[string]string{
		"group":     "mobile",
		"namespace": ms.Namespace,
	}
	for k, v := range ms.Labels {
		labels[k] = v
	}
	data["uri"] = []byte(ms.Host)
	data["name"] = []byte(ms.Name)
	data["displayName"] = []byte(ms.DisplayName)
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
	if nil == svs {
		return nil, errors.New("no secrets returned for list")
	}
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

func (msr *MobileServiceRepo) Delete(serviceID string) error {
	if err := msr.client.Delete(serviceID, &meta_v1.DeleteOptions{}); err != nil {
		return errors.Wrap(err, "mobile serive repo failed to delete underlying secret")
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
		Namespace:    s.Labels["namespace"],
		ID:           s.Name,
		External:     external,
		Labels:       s.Labels,
		Name:         strings.TrimSpace(string(s.Data["name"])),
		DisplayName:  strings.TrimSpace(retrieveDisplayNameFromSecret(s)),
		Type:         strings.TrimSpace(string(s.Data["type"])),
		Host:         string(s.Data["uri"]),
		Params:       params,
		Integrations: map[string]*mobile.ServiceIntegration{},
	}
}

// NewServiceRepoBuilder provides an implementation of mobile.ServiceRepoBuilder
func NewServiceRepoBuilder(clientBuilder mobile.K8ClientBuilder, namespace, saToken string) mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{
		clientBuilder: clientBuilder,
		saToken:       saToken,
		namespace:     namespace,
	}
}

// MobileServiceRepoBuilder builds a ServiceCruder
type MobileServiceRepoBuilder struct {
	clientBuilder mobile.K8ClientBuilder
	token         string
	namespace     string
	saToken       string
}

func (marb *MobileServiceRepoBuilder) WithToken(token string) mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{
		clientBuilder: marb.clientBuilder,
		token:         token,
		saToken:       marb.saToken,
		namespace:     marb.namespace,
	}
}

//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
func (marb *MobileServiceRepoBuilder) UseDefaultSAToken() mobile.ServiceRepoBuilder {
	return &MobileServiceRepoBuilder{
		clientBuilder: marb.clientBuilder,
		token:         marb.saToken,
		saToken:       marb.saToken,
		namespace:     marb.namespace,
	}

}

// Build builds the final repo
func (marb *MobileServiceRepoBuilder) Build() (mobile.ServiceCruder, error) {
	k8client, err := marb.clientBuilder.WithToken(marb.token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "MobileAppRepoBuilder failed to build a configmap client")
	}
	return NewMobileServiceRepo(k8client.CoreV1().Secrets(marb.namespace)), nil
}

// If there is no display name in the secret then we will use the service name
func retrieveDisplayNameFromSecret(sec v1.Secret) string {
	if string(sec.Data["displayName"]) == "" {
		return string(sec.Data["name"])
	}
	return string(sec.Data["displayName"])
}
