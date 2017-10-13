package k8s

import (
	"fmt"
	"testing"

	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
)

var token = "j6hCIHz8O1yk19C-WJi7J4mR1Aq6X98Cb7KBqATYWck"
var ns = "mobileapp"
var host = "https://192.168.37.1:8443"

func TestSCBindKeycloak(t *testing.T) {
	cases := []struct {
		Name string
	}{
		{
			Name: "test generate binding",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			externalReq := httpclient.NewClientBuilder().Insecure(true).Build()
			k8Builder := NewClientBuilder(ns, host, false)
			kc, err := k8Builder.WithToken(token).BuildClient()
			if err != nil {
				t.Fatal("failed to get k8 client", err)
			}
			sc := &serviceCatalogClient{namespace: ns, token: token, externalRequester: externalReq, k8host: "https://192.168.37.1:8443", k8client: kc}
			if err := sc.BindServiceToKeyCloak("fh-sync-server", ns); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestScGetInstances(t *testing.T) {
	cases := []struct {
		Name string
	}{
		{
			Name: "test get instances",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			sc := &serviceCatalogClient{k8host: "https://192.168.37.1:8443"}
			sc.getInstances("sVLIQ2J4iCokFBOS3PGVWdLkxOC-5o5iKwJ6xrZuG0k", "myproject")
		})
	}
}

func TestScGetClasses(t *testing.T) {
	cases := []struct {
		Name string
	}{
		{
			Name: "test get instances",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			sc := &serviceCatalogClient{k8host: "https://192.168.37.1:8443"}
			sc.serviceClasses(token, "myproject")
		})
	}
}

func TestSC_GetServiceClassByServiceName(t *testing.T) {
	cases := []struct {
		Name string
	}{
		{
			Name: "test get service class",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			sc := &serviceCatalogClient{k8host: "https://192.168.37.1:8443"}
			class, _ := sc.serviceClassByServiceName("keycloak", token)
			fmt.Print("classs", class)
		})
	}
}
