package mobile

import "strings"

// App represents a mobile app
type App struct {
	Name       string            `json:"name"`
	ClientType string            `json:"clientType"`
	Labels     map[string]string `json:"labels"`
	APIKey     string            `json:"apiKey"`
}

type StatusError struct {
	Message string
	Code    int
}

type Service struct {
	Name              string            `json:"name"`
	Host              string            `json:"host"`
	Params            map[string]string `json:"params"`
	BindingSecretName string            `json:"binding_secret_name"`
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
	//AuthHeader the header where authorisation token is stored
	AuthHeader = "x-auth"
	//AppAPIKeyHeader is the header sent by mobile clients when they want to interact with mcp
	AppAPIKeyHeader = "x-app-api-key"
	//SkipSARoleBindingHeader is the head the admin api key is sent with
	SkipSARoleBindingHeader = "x-skip-role-binding"
)
