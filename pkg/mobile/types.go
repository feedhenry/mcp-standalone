package mobile

import (
	"strings"
)

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

//Service represents a serverside application that mobile application will interact with
type Service struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	Namespace    string                         `json:"namespace"`
	Host         string                         `json:"host"`
	Description  string                         `json:"description"`
	Type         string                         `json:"type"`
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

//TODO move out to config or env var
//ServiceTypes are the service types that we are aware of and support
var ServiceTypes = []string{"fh-sync-server", "keycloak", "custom"}

const (
	//AppAPIKeyHeader is the header sent by mobile clients when they want to interact with mcp
	AppAPIKeyHeader = "x-app-api-key"
)

//'total': '250',
//'xData': ["dates", Fri Aug 25 2017 10:53:04 GMT+0100 (IST), Sat Aug 26 2017 10:53:04 GMT+0100 (IST)]
//'yData': ['used', '20', '20', '35', '70', '20', '87', '14', '95', '25', '28', '44', '56', '66', '16', '67', '88', '76', '65', '87']

// Metric is a common container for returning metrics to the dashboard
type Metric struct {
	Type  string   `json:"type"`
	Total string   `json:"total"`
	XData []string `json:"xData"`
	YData []string `json:"yData"`
}

type MetricList []*Metric

func (mm MetricList) Get(t string) *Metric {
	for _, m := range mm {
		if m.Type == t {
			return m
		}
	}
	return nil
}
