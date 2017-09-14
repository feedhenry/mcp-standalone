package metrics

import "time"

type Keycloak struct{}

// Gather will retrieve varous metrics from keycloak
func (kc *Keycloak) Gather() ([]*metric, error) {
	m := &metric{Type: "logins", XValue: time.Now().String(), YValue: 10}
	return []*metric{m}, nil
}
