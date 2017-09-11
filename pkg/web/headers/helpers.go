package headers

import (
	"net/http"
	"strings"
)

// DefaultTokenRetriever returns the authorization token
func DefaultTokenRetriever(headers http.Header) string {
	token := headers.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	return strings.TrimSpace(token)
}
