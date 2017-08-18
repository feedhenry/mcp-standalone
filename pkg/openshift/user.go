package openshift

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	userReadPath = "/oapi/v1/users/~"
)

type UserAccess struct {
	Logger *logrus.Logger
}

func (ua *UserAccess) ReadUserFromToken(host, token string, insecure bool) error {
	u, err := url.Parse(host)
	if err != nil {
		return errors.Wrap(err, "failed to parse openshift host when attempting to read user")
	}
	u.Path = path.Join(userReadPath)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to build request to read user")
	}
	req.Header.Set("authorization", "bearer "+token)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make request to read user from openshift")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return &AuthenticatationError{Message: "access was denied", StatusCode: resp.StatusCode}
		}

		return errors.New(fmt.Sprintf("unexpected response code from openshift %v", resp.StatusCode))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read the response body after reading user")
	}
	ua.Logger.Debug("response from reading user ", string(data))
	return nil
}

type AuthenticatationError struct {
	Message    string
	StatusCode int
}

func (ae *AuthenticatationError) Error() string {
	return ae.Message
}

func (ae *AuthenticatationError) Code() int {
	return ae.StatusCode
}

func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticatationError)
	return ok
}
