package metrics

import "time"

type Keycloak struct{}

// Gather will retrieve varous metrics from keycloak
func (kc *Keycloak) Gather() ([]*metric, error) {
	now := time.Now()

	m := &metric{Type: "logins", XValue: now.Format("2006-01-02 15:04:05"), YValue: 10}
	return []*metric{m}, nil
}
