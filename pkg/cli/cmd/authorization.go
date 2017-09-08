package cmd

import (
	"net/http"

	"github.com/spf13/viper"
)

// add the authorization header needed using viper to pull it from the cmd line or cfg file
func addAuthorizationHeader(headers http.Header) {
	headers.Set("Authorization", "Bearer "+viper.GetString("token"))
}
