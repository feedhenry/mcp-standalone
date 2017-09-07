package mobile

import "strings"

// App represents a mobile app
type App struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ClientType  string            `json:"clientType"`
	Labels      map[string]string `json:"labels"`
	APIKey      string            `json:"apiKey"`
	MetaData    map[string]string `json:"metadata"`
}

type StatusError struct {
	Message string
	Code    int
}

type Service struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	Host         string                         `json:"host"`
	Description  string                         `json:"description"`
	Capabilities map[string][]string            `json:"capabilities"`
	Params       map[string]string              `json:"params"`
	Labels       map[string]string              `json:"labels"`
	Integrations map[string]*ServiceIntegration `json:"integrations"`
	External     bool                           `json:"external"`
}

func NewMobileService() *Service {
	return &Service{
		Capabilities: map[string][]string{},
		Params:       map[string]string{},
		Labels:       map[string]string{"group": "mobile"},
		Integrations: map[string]*ServiceIntegration{},
	}
}

type ServiceIntegration struct {
	Enabled         bool   `json:"enabled"`
	Component       string `json:"component"`
	Service         string `json:"service"`
	Namespace       string `json:"namespace"`
	ComponentSecret string `json:"componentSecret"`
}

type ServiceConfig struct {
	Config interface{} `json:"config"`
	Name   string      `json:"name"`
}

type KeycloakConfig struct {
	SSLRequired   string `json:"ssl-required"`
	AuthServerURL string `json:"auth-server-url"`
	Realm         string `json:"realm"`
	Resource      string `json:"resource"`
	ClientID      string `json:"clientId"`
	URL           string `json:"url"`
	Credentials   struct {
		Secret string `json:"secret"`
	} `json:"credentials"`
}

type AttrFilterFunc func(attrs Attributer) bool

func (se *StatusError) Error() string {
	return se.Message
}

func (se *StatusError) StatusCode() int {
	return se.Code
}

//AppTypes are the valid app types
type AppTypes []string

func (at AppTypes) Contains(v string) bool {
	for _, val := range at {
		if v == val {
			return true
		}
	}
	return false
}

func (at AppTypes) String() string {
	return strings.Join(at, " : ")
}

//ValidAppTypes is a list of valid app types
var ValidAppTypes = AppTypes{"cordova", "android", "iOS"}

const (
	//AppAPIKeyHeader is the header sent by mobile clients when they want to interact with mcp
	AppAPIKeyHeader = "x-app-api-key"
	//SkipSARoleBindingHeader is the head the admin api key is sent with
	SkipSARoleBindingHeader = "x-skip-role-binding"
)
