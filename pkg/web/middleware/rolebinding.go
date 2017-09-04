package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

// RoleBinding middlware is a HACKY work around for issue
//https://github.com/openshift/ansible-service-broker/issues/399
// it will check for the existence of the edit role for the mobile service account
// if it is not present, it will create it and grant the edit role
type RoleBinding struct {
	clientBuilder mobile.HTTPRequesterBuilder
	namespace     string
	khost         string
	logger        *logrus.Logger
	*sync.Mutex
	roleBindingExists bool
}

// NewRoleBinding creates new rolebinding middleware
func NewRoleBinding(cb mobile.HTTPRequesterBuilder, namespace string, logger *logrus.Logger, host string) *RoleBinding {
	return &RoleBinding{
		namespace:     namespace,
		clientBuilder: cb,
		logger:        logger,
		khost:         host,
		Mutex:         &sync.Mutex{},
	}
}

// Handle will handle the rolebinding interaction
func (sa *RoleBinding) Handle(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if shouldIgnore(req.URL.Path) {
		next(rw, req)
		return
	}
	// only want one request at a time trying to do this.
	sa.Lock()
	defer sa.Unlock()
	sa.logger.Debug("headers", req.Header)
	token := req.Header.Get(mobile.AuthHeader)
	if req.Header.Get(mobile.SkipSARoleBindingHeader) != "" {
		sa.logger.Debug("skipping sa role binding")
		next(rw, req)
		return
	}
	if err := sa.createRoleBindingIfNotPresent(token); err != nil {
		sa.logger.Error(fmt.Sprintf("error when setting up rolebinding: %+v", err))
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	next(rw, req)
}

func (sa *RoleBinding) createRoleBindingIfNotPresent(token string) error {
	if sa.roleBindingExists {
		sa.logger.Debug("role binding already exists")
		return nil
	}
	has, err := sa.hasRoleBinding(token, "edit")
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	binding := &v1beta1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels: map[string]string{"group": "mobileapp"},
			Name:   "edit",
		},
		RoleRef: v1beta1.RoleRef{
			Name: "edit",
		},
		Subjects: []v1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "mcp-standalone",
				Namespace: sa.namespace,
			},
		},
	}
	data, err := json.Marshal(binding)
	if err != nil {
		return errors.Wrap(err, "failed to marshal rolebinding to create")
	}
	// set insecure via a config flag in future
	client := sa.clientBuilder.Insecure(true).Timeout(3).Build()
	url := sa.khost + fmt.Sprintf("/oapi/v1/namespaces/%s/rolebindings", sa.namespace)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "failed to create request to create a new rolebinding")
	}
	req.Header.Set("Authorization", " Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do request to create rolebinding")
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			sa.logger.Error("attempted to read body after failed rolebinding create ", err.Error())
		}
		if nil != data {
			sa.logger.Error("unexpected response from OpenShift ", string(data))
		}
		return errors.New("unexpected response from OpenShift: " + res.Status)
	}
	//created so we know it now exists
	sa.roleBindingExists = true
	return nil
}

func (sa *RoleBinding) hasRoleBinding(token, name string) (bool, error) {
	//TODO set insecure in config
	client := sa.clientBuilder.Insecure(true).Timeout(3).Build()
	url := sa.khost + fmt.Sprintf("/oapi/v1/namespaces/%s/rolebindings/%s", sa.namespace, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to create request to get rolebinding "+name)
	}
	req.Header.Set("Authorization", " Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "failed to make request to get rolebinding "+name)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		sa.logger.Debug("rolebinding edit does not exist")
		return false, nil
	}
	if res.StatusCode != 200 {
		return false, errors.New(fmt.Sprintf("unexpected response code from openshift %v", res.StatusCode))
	}
	sa.roleBindingExists = true
	return true, nil
}
