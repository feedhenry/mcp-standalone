/**
Doing this manually for now, but should probably look at bring in the service catalog client at some point or just calling the api
from the ui client https://github.com/kubernetes-incubator/service-catalog/issues/1367 .
Longer term this is likely to be removed once OSCP 3.7 is released and supports binding parameters
*/
package k8s

import (
	"fmt"

	"net/http"
	"strings"

	"io/ioutil"

	"crypto/tls"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	v1alpha1 "k8s.io/client-go/pkg/apis/settings/v1alpha1"
)

var (
	bindURL         = "%s/apis/servicecatalog.k8s.io/v1alpha1/namespaces/%s/bindings"
	instanceURL     = "%s/apis/servicecatalog.k8s.io/v1alpha1/namespaces/%s/instances"
	serviceClassURL = "%s/apis/servicecatalog.k8s.io/v1alpha1/serviceclasses"
)

type ServiceCatalogClientBuilder struct {
	httpClient mobile.ExternalHTTPRequester
	k8host     string
	token      string
	namespace  string
	saToken    string
	k8builder  mobile.K8ClientBuilder
}

func (scb *ServiceCatalogClientBuilder) WithToken(token string) mobile.SCClientBuilder {
	s := *scb
	s.token = token
	return &s
}

func (scb *ServiceCatalogClientBuilder) UseDefaultSAToken() mobile.SCClientBuilder {
	s := *scb
	s.token = scb.saToken
	return &s
}

func (scb *ServiceCatalogClientBuilder) WithHost(host string) mobile.SCClientBuilder {
	s := *scb
	s.k8host = host
	return &s
}
func (scb *ServiceCatalogClientBuilder) Build() (mobile.SCCInterface, error) {
	k8, err := scb.k8builder.WithToken(scb.token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build k8 client required by service catalog client")
	}
	return &serviceCatalogClient{
		k8host:            scb.k8host,
		token:             scb.token,
		namespace:         scb.namespace,
		externalRequester: scb.httpClient,
		k8client:          k8,
	}, nil
}

func NewServiceCatalogClientBuilder(builder mobile.K8ClientBuilder, httpClient mobile.ExternalHTTPRequester, saToken, namespace, k8host string) *ServiceCatalogClientBuilder {
	return &ServiceCatalogClientBuilder{httpClient: httpClient, saToken: saToken, namespace: namespace, k8builder: builder, k8host: k8host}
}

type serviceCatalogClient struct {
	k8host            string
	token             string
	namespace         string
	externalRequester mobile.ExternalHTTPRequester
	k8client          kubernetes.Interface
}
type ServiceClass struct {
	meta_v1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	meta_v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	ExternalMetadata map[string]string `json:"externalMetadata"`
}

type ServiceInstance struct {
	meta_v1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	meta_v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec               struct {
		ServiceClassName string `json:"serviceClassName"`
	} `json:"spec"`
}

type ServiceClassList struct {
	Items []ServiceClass `json:"items"`
}

type ServiceInstanceList struct {
	Items []ServiceInstance `json:"items"`
}

//TODO this is fragile and should be changed to use the real types and client https://github.com/kubernetes-incubator/service-catalog/issues/1367
func createBindingObject(instance string, params map[string]string, secretName string) (string, error) {
	pdata, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return `{"kind":"Binding","apiVersion":"servicecatalog.k8s.io/v1alpha1","metadata":{"generateName":"` + instance + `-"},
	 "spec":{
	     "instanceRef":{"name":"` + instance + `"},
	     "secretName":"` + secretName + `",
	     "parameters":` + string(pdata) + `
     }
	}`, nil
}

func (sc *serviceCatalogClient) podPreset(objectName, svcName, targetSvcName, namespace string) error {
	podPreset := v1alpha1.PodPreset{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: objectName,
			Labels: map[string]string{
				"group":   "mobile",
				"service": svcName,
			},
		},
		Spec: v1alpha1.PodPresetSpec{
			Selector: meta_v1.LabelSelector{
				MatchLabels: map[string]string{
					"run": targetSvcName,
				},
			},
			Volumes: []v1.Volume{
				v1.Volume{
					Name: svcName,
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: objectName,
						},
					},
				},
			},
			VolumeMounts: []v1.VolumeMount{
				v1.VolumeMount{
					Name:      svcName,
					MountPath: "/etc/secrets/" + svcName,
				},
			},
		},
	}
	if _, err := sc.k8client.SettingsV1alpha1().PodPresets(namespace).Create(&podPreset); err != nil {
		return errors.Wrap(err, "failed to create pod preset for service ")
	}
	return nil
}

// BindServiceToKeyCloak will create a binding and pod preset
// finds the service class based on the service name label (keycloak in this case)
// finds the first service instances of keycloak
// creates a binding via service catalog which kicks of the keycloak apb
// finally creates a pod preset for sync pods to pick up as a volume mount
func (sc *serviceCatalogClient) BindServiceToKeyCloak(targetSvcName, namespace string) error {
	objectName := "keycloak-" + targetSvcName
	keyCloakServiceClass, err := sc.serviceClassByServiceName(mobile.ServiceNameKeycloak, sc.token)
	if err != nil {
		return err
	}
	keycloakInstList, err := sc.serviceInstancesForServiceClass(sc.token, keyCloakServiceClass.Name, namespace)
	if err != nil {
		return err
	}
	if len(keycloakInstList.Items) == 0 {
		return errors.New("no instances of keycloak found")
	}

	// only care about the first one as there only should ever be one.
	keycloakInst := keycloakInstList.Items[0]

	pbody, _ := createBindingObject(keycloakInst.Name, map[string]string{"service": targetSvcName}, objectName)
	req, err := http.NewRequest("POST", fmt.Sprintf(bindURL, sc.k8host, namespace), strings.NewReader(pbody))
	if err != nil {
		fmt.Println("failed to create request ", err)
	}
	req.Header.Set("Authorization", "Bearer "+sc.token)
	res, err := sc.externalRequester.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to create binding. Request error occurred")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		return errors.New("unexpected status code from service catalog: " + res.Status)
	}
	if err := sc.podPreset(objectName, mobile.ServiceNameKeycloak, targetSvcName, namespace); err != nil {
		return errors.WithStack(err)
	}
	//update the deployment with an annotation
	dep, err := sc.k8client.AppsV1beta1().Deployments(namespace).Get(targetSvcName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get deployment for service "+targetSvcName)
	}
	dep.Spec.Template.Labels[mobile.ServiceNameKeycloak] = "enabled"
	if _, err := sc.k8client.AppsV1beta1().Deployments(namespace).Update(dep); err != nil {
		return errors.Wrap(err, "failed up update deployment for "+targetSvcName)
	}
	return nil
}

//UnBindServiceToKeyCloak will Delete the binding, the pod preset and the update the deployment
func (sc *serviceCatalogClient) UnBindServiceToKeyCloak(targetSvcName, namespace string) error {
	return nil
}

func (sc *serviceCatalogClient) getInstances(token, ns string) (*ServiceInstanceList, error) {
	u := fmt.Sprintf(instanceURL, sc.k8host, ns)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to get service instances")
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the response body for service instances")
	}
	sil := &ServiceInstanceList{}
	if err := json.Unmarshal(data, sil); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body into a service instance list")
	}
	return sil, nil
}

func (sc *serviceCatalogClient) serviceClasses(token, ns string) ([]ServiceClass, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(serviceClassURL, sc.k8host), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error making service class request ", err)
		return nil, err
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	scs := &ServiceClassList{}
	if err := json.Unmarshal(data, &scs); err != nil {
		fmt.Println("error making service class request ", err)
		return nil, err
	}
	fmt.Println("service classes ", scs)

	return scs.Items, nil

}

func (sc *serviceCatalogClient) serviceClassByServiceName(name, token string) (*ServiceClass, error) {
	serviceClasses, err := sc.serviceClasses(token, "")
	if err != nil {
		return nil, err
	}

	for _, sc := range serviceClasses {
		if v, ok := sc.ExternalMetadata["serviceName"]; ok && v == name {
			return &sc, nil
		}
	}
	return nil, nil
}

func (sc *serviceCatalogClient) serviceInstancesForServiceClass(token, serviceClass string, ns string) (*ServiceInstanceList, error) {
	si, err := sc.getInstances(token, ns)
	fmt.Println(si)
	if err != nil {
		return nil, err
	}
	sl := &ServiceInstanceList{}
	for _, i := range si.Items {
		if i.Spec.ServiceClassName == serviceClass {
			sl.Items = append(sl.Items, i)
		}
	}
	return sl, nil
}
