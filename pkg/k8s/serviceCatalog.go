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

	"encoding/json"

	"bytes"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	v1alpha1 "k8s.io/client-go/pkg/apis/settings/v1alpha1"
)

var (
	bindURL         = "%s/apis/servicecatalog.k8s.io/v1alpha1/namespaces/%s/servicebindings"
	instanceURL     = "%s/apis/servicecatalog.k8s.io/v1alpha1/namespaces/%s/serviceinstances"
	serviceClassURL = "%s/apis/servicecatalog.k8s.io/v1alpha1/clusterserviceclasses"
	bindingURL      = "%s/apis/servicecatalog.k8s.io/v1alpha1/namespaces/%s/servicebinding/%s"
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
type Binding struct {
	meta_v1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	meta_v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec               struct {
		InstanceRef struct {
			Name string `json:"name"`
		} `json:"instanceRef"`
		Parameters struct {
			Test string `json:"test"`
		} `json:"parameters"`
		SecretName string `json:"secretName"`
	} `json:"spec"`
}

//TODO this is fragile and should be changed to use the real types and client https://github.com/kubernetes-incubator/service-catalog/issues/1367
func createBindingObject(instance string, params map[string]interface{}, secretName string) (string, error) {
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

func (sc *serviceCatalogClient) podPreset(objectName, secretName, svcName, targetSvcName, namespace string) error {
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
					"run":   targetSvcName,
					svcName: "enabled",
				},
			},
			Volumes: []v1.Volume{
				v1.Volume{
					Name: svcName,
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: secretName,
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

// BindToService will create a binding and pod preset
// finds the service class based on the service name
// finds the first service instances
// creates a binding via service catalog which kicks of the bind apb for the service
// finally creates a pod preset for sync pods to pick up as a volume mount
// TODO perhaps the pod preset could be created as part of the bind API in the apb (would need to pass parameters)
func (sc *serviceCatalogClient) BindToService(bindableService, targetSvcName string, params map[string]interface{}, bindableSvcNamespace, targetSvcNamespace string) error {
	objectName := bindableService + "-" + targetSvcName
	bindableServiceClass, err := sc.serviceClassByServiceName(bindableService, sc.token)
	if err != nil {
		return err
	}
	if nil == bindableServiceClass {
		return errors.New("failed to find service class for service " + bindableService)
	}

	svcInstList, err := sc.serviceInstancesForServiceClass(sc.token, bindableServiceClass.Name, targetSvcNamespace)
	if err != nil {
		return err
	}
	if len(svcInstList.Items) == 0 {
		return errors.New("no service instance of " + bindableService + " found in ns " + targetSvcNamespace)
	}

	// only care about the first one as there only should ever be one.
	svcInst := svcInstList.Items[0]
	pbody, _ := createBindingObject(svcInst.Name, params, objectName)
	req, err := http.NewRequest("POST", fmt.Sprintf(bindURL, sc.k8host, targetSvcNamespace), strings.NewReader(pbody))
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
	bindingResp := &Binding{}
	if err := json.NewDecoder(res.Body).Decode(bindingResp); err != nil {
		return errors.WithStack(err)
	}
	if err := sc.podPreset(objectName, objectName, bindableService, targetSvcName, targetSvcNamespace); err != nil {
		return errors.WithStack(err)
	}
	//update the deployment with an annotation
	dep, err := sc.k8client.AppsV1beta1().Deployments(targetSvcNamespace).Get(targetSvcName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get deployment for service "+targetSvcName)
	}

	dep.Spec.Template.Labels[bindableService] = "enabled"
	dep.Spec.Template.Labels[bindableService+"-binding"] = bindingResp.Name
	if _, err := sc.k8client.AppsV1beta1().Deployments(targetSvcNamespace).Update(dep); err != nil {
		return errors.Wrap(err, "failed up update deployment for "+targetSvcName)
	}
	return nil
}

//UnBindFromService will Delete the binding, the pod preset and the update the deployment
// TODO again deleting the pod preset may be better done in the asb ubind handler
func (sc *serviceCatalogClient) UnBindFromService(bindableService, targetSvcName, targetSvcNamespace string) error {
	objectName := bindableService + "-" + targetSvcName
	dep, err := sc.k8client.AppsV1beta1().Deployments(targetSvcNamespace).Get(targetSvcName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get deployment for service "+targetSvcName)
	}
	bindingID, ok := dep.Spec.Template.Labels[bindableService+"-binding"]
	if !ok {
		return errors.New("no binding id found for service " + targetSvcName)
	}
	delete(dep.Spec.Template.Labels, bindableService+"-binding")
	delete(dep.Spec.Template.Labels, bindableService)

	unbindURL := fmt.Sprintf(bindingURL, sc.k8host, targetSvcNamespace, bindingID)
	body := meta_v1.DeleteOptions{
		TypeMeta: meta_v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "DeleteOptions",
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "failed to decode delete options")
	}
	req, err := http.NewRequest("DELETE", unbindURL, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "failed to create delete request for binding ")
	}
	req.Header.Set("Authorization", "Bearer "+sc.token)
	res, err := sc.externalRequester.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do delete request against the service catalog to delete binding "+bindingID)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("unexpected response code from service catalog " + res.Status)
	}
	// binding deleted we will remove the pod preset and update deployment
	if err := sc.k8client.SettingsV1alpha1().PodPresets(targetSvcNamespace).Delete(objectName, meta_v1.NewDeleteOptions(0)); err != nil {
		return errors.Wrap(err, "unbinding "+bindableService+" and "+targetSvcName+" failed to delete pod preset")
	}
	if _, err := sc.k8client.AppsV1beta1().Deployments(targetSvcNamespace).Update(dep); err != nil {
		return errors.Wrap(err, "failed to update the deployment for "+targetSvcName+" after unbinding "+bindableService)
	}
	return nil
}

// create pod preset with apikeys secret, update deployment with label
func (sc *serviceCatalogClient) AddMobileApiKeys(targetSvcName, namespace string) error {
	objectName := mobile.IntegrationAPIKeys + "-" + targetSvcName
	if err := sc.podPreset(objectName, mobile.IntegrationAPIKeys, mobile.IntegrationAPIKeys, targetSvcName, namespace); err != nil {
		return errors.WithStack(err)
	}
	dep, err := sc.k8client.AppsV1beta1().Deployments(namespace).Get(targetSvcName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get deployment for service "+targetSvcName+" cannot redeploy.")
	}
	dep.Spec.Template.Labels[mobile.IntegrationAPIKeys] = "enabled"
	if _, err := sc.k8client.AppsV1beta1().Deployments(namespace).Update(dep); err != nil {
		return errors.Wrap(err, "failed up update deployment for "+targetSvcName)
	}
	return nil
}

// create pod preset with apikeys secret, update deployment with label
func (sc *serviceCatalogClient) RemoveMobileApiKeys(targetSvcName, namespace string) error {
	objectName := mobile.IntegrationAPIKeys + "-" + targetSvcName
	if err := sc.k8client.SettingsV1alpha1().PodPresets(namespace).Delete(objectName, meta_v1.NewDeleteOptions(0)); err != nil {
		return errors.Wrap(err, "removing api keys failed to delete pod preset")
	}
	dep, err := sc.k8client.AppsV1beta1().Deployments(namespace).Get(targetSvcName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get deployment for service "+targetSvcName+" cannot redeploy.")
	}

	delete(dep.Spec.Template.Labels, mobile.IntegrationAPIKeys)
	if _, err := sc.k8client.AppsV1beta1().Deployments(namespace).Update(dep); err != nil {
		return errors.Wrap(err, "failed up update deployment for "+targetSvcName)
	}
	return nil
}

func (sc *serviceCatalogClient) getInstances(token, ns string) (*ServiceInstanceList, error) {
	u := fmt.Sprintf(instanceURL, sc.k8host, ns)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := sc.externalRequester.Do(req)
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
	res, err := sc.externalRequester.Do(req)
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
