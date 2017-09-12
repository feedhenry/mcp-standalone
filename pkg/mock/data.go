package mock

import (
	v1 "k8s.io/client-go/pkg/api/v1"
)

func KeycloakSecret() v1.Secret {
	return v1.Secret{
		Data: map[string][]byte{
			"name": []byte("keycloak"),
			"type": []byte("keycloak"),
			"uri":  []byte("http://keycloak.com"),
			"public_installation": []byte("{\"ssl-required\": \"external\", \"auth-server-url\": \"http://keycloak-authmobile.192.168.37.1.nip.io/auth\", \"realm\": \"authmobile\", \"resource\": \"zlFwvnNrkRnBOfeuUVss\", \"credentials\": {\"secret\": \"m3onpFzg2xSm0K4Hze83\"}}"),
		},
	}
}
