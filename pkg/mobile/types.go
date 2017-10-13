//Defines the core mobile domain objects.
//This package should not depend on any of our other packages, it is ok for it to use std lib packages.

package mobile

import (
	"errors"
	"net/url"
	"strings"
)

const (
	ServiceNameKeycloak   = "keycloak"
	ServiceNameThreeScale = "3scale"
	ServiceNameSync       = "fh-sync-server"
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

//Service represents a serverside application that mobile application will interact with
type Service struct {
	ID           string                         `json:"id"`
	Name         string                         `json:"name"`
	DisplayName  string                         `json:"displayName"`
	Namespace    string                         `json:"namespace"`
	Host         string                         `json:"host"`
	Description  string                         `json:"description"`
	Type         string                         `json:"type"`
	Capabilities map[string][]string            `json:"capabilities"`
	Params       map[string]string              `json:"params"`
	Labels       map[string]string              `json:"labels"`
	Integrations map[string]*ServiceIntegration `json:"integrations"`
	External     bool                           `json:"external"`
	Writable     bool                           `json:"writable"`
}

type Build struct {
	Name     string         `json:"name"`
	Download *BuildDownload `json:"download"`
}

type BuildDownload struct {
	URL     string `json:"url"`
	Expires int64  `json:"expires"`
	Token   string `json:"-"`
}

// BuildConfig represents a build of a mobile client. It is converted to a buildconfig
type BuildConfig struct {
	AppID   string        `json:"appID"`
	Name    string        `json:"name"`
	GitRepo *BuildGitRepo `json:"gitRepo"`
}

type BuildGitRepo struct {
	URI             string `json:"uri"`
	Private         bool   `json:"private"`
	Ref             string `json:"ref"`
	PublicKey       string `json:"public"`
	PublicKeyID     string `json:"publicKeyId"`
	JenkinsFilePath string `json:"jenkinsFilePath"`
}

type BuildStatus struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Artifacts struct {
			Href string `json:"href"`
		} `json:"artifacts"`
	} `json:"_links"`
	Phase string `json:"phase"`
}

func (bs *BuildStatus) Host() (string, error) {
	if bs.Links.Self.Href == "" {
		return "", errors.New("href is missing from build status")
	}
	u, err := url.Parse(bs.Links.Self.Href)
	if err != nil {
		return "", err
	}
	return u.Scheme + "://" + u.Host, nil
}

func (bs *BuildStatus) ArtifactURL() (*url.URL, error) {
	host, err := bs.Host()
	if err != nil {
		return nil, err
	}
	return url.Parse(host + bs.Links.Artifacts.Href)
}

type BuildAsset struct {
	Password  string
	Platform  string
	Path      string
	BuildName string
	AppName   string
	Name      string
	Type      BuildAssetType
	AssetData map[string][]byte
}

func (ba *BuildAsset) Validate(assetType BuildAssetType) error {

	if assetType == BuildAssetTypeBuildSecret {
		if !ValidAppTypes.Contains(ba.Platform) {
			return errors.New("build asset of type build secret needs a valid platform")
		}
		if ba.Platform == "android" && ba.Password == "" {
			return errors.New("build assets for android require a password")
		}
	} else {
		if ba.BuildName == "" {
			return errors.New("build asset needs a valid build name")
		}
	}
	return nil
}

type BuildAssetType string

var (
	BuildAssetTypeSourceCredential BuildAssetType = "mobile-src"
	BuildAssetTypeBuildSecret      BuildAssetType = "mobile-build"
)

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
	DisplayName     string `json:"displayName"`
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

// AndroidApp Type for Android app
const AndroidApp = "android"

// IOSApp Type for iOS app
const IOSApp = "iOS"

// CordovaApp Type for Cordova app
const CordovaApp = "cordova"

//ValidAppTypes is a list of valid app types
var ValidAppTypes = AppTypes{CordovaApp, AndroidApp, IOSApp}

//TODO move out to config or env var
//ServiceTypes are the service types that we are aware of and support
var ServiceTypes = []string{"fh-sync-server", "keycloak", "aerogear-digger", "custom"}

const (
	//AppAPIKeyHeader is the header sent by mobile clients when they want to interact with mcp
	AppAPIKeyHeader = "x-app-api-key"
)

//'total': '250',
//'xData': ["dates", Fri Aug 25 2017 10:53:04 GMT+0100 (IST), Sat Aug 26 2017 10:53:04 GMT+0100 (IST)]
//'yData': ['used', '20', '20', '35', '70', '20', '87', '14', '95', '25', '28', '44', '56', '66', '16', '67', '88', '76', '65', '87']

// GatheredMetric is a common container for returning metrics to the dashboard
type GatheredMetric struct {
	Type string             `json:"type"`
	X    []string           `json:"x"`
	Y    map[string][]int64 `json:"y"`
}

type User struct {
	User   string
	Groups []string
}

func (u *User) InAnyGroup(groups []string) bool {
	for _, group := range groups {
		for _, userGroup := range u.Groups {
			if group == userGroup {
				return true
			}
		}
	}
	return false
}
