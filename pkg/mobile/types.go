package mobile

import "strings"

// App represents a mobile app
type App struct {
	Name       string            `json:"name"`
	ClientType string            `json:"clientType"`
	Labels     map[string]string `json:"labels"`
}

type StatusError struct {
	Message string
	Code    int
}

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
