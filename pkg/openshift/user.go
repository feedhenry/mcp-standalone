package openshift

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"encoding/json"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

const (
	userReadPath = "/oapi/v1/users/~"
)

type UserChecker func(host, token string, skipTLS bool) (*mobile.User, error)

type UserAccess struct{}

type userResponse struct {
	Identities []string `json:"identities"`
	Groups     []string `json:"groups"`
	Metadata   struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

func (ua *UserAccess) ReadUserFromToken(host, token string, insecure bool) (*mobile.User, error) {
	user := &mobile.User{}
	u, err := url.Parse(host)
	if err != nil {
		return user, errors.Wrap(err, "failed to parse openshift host when attempting to read user")
	}
	u.Path = path.Join(userReadPath)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return user, errors.Wrap(err, "failed to build request to read user")
	}
	req.Header.Set("authorization", "bearer "+token)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return user, errors.Wrap(err, "failed to make request to read user from openshift")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Error("failed to close response body. can cause file handle leaks ", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return user, &AuthenticationError{Message: "access was denied", StatusCode: resp.StatusCode}
		}

		return user, errors.New("unexpected response code from openshift " + strconv.Itoa(resp.StatusCode))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, errors.Wrap(err, "failed to read the response body after reading user")
	}

	userData := &userResponse{}
	err = json.Unmarshal(data, userData)
	if err != nil {
		return user, err
	}
	user.User = userData.Metadata.Name
	user.Groups = userData.Groups

	return user, nil
}

type AuthenticationError struct {
	Message    string
	StatusCode int
}

func (ae *AuthenticationError) Error() string {
	return ae.Message
}

func (ae *AuthenticationError) Code() int {
	return ae.StatusCode
}

func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticationError)
	return ok
}

type UserRepoBuilder struct {
	token       string
	client      mobile.UserAccessChecker
	host        string
	ignoreCerts bool
}

func NewUserRepoBuilder(host string, ignoreCerts bool) mobile.UserRepoBuilder {
	return &UserRepoBuilder{host: host, ignoreCerts: ignoreCerts}
}

func (urb *UserRepoBuilder) WithToken(token string) mobile.UserRepoBuilder {
	return &UserRepoBuilder{token: token, client: urb.client, host: urb.host, ignoreCerts: urb.ignoreCerts}
}

func (urb *UserRepoBuilder) WithClient(client mobile.UserAccessChecker) mobile.UserRepoBuilder {
	return &UserRepoBuilder{token: urb.token, client: client, host: urb.host, ignoreCerts: urb.ignoreCerts}
}

func (urb *UserRepoBuilder) Build() mobile.UserRepo {
	return &UserRepo{token: urb.token, client: urb.client, host: urb.host, ignoreCerts: urb.ignoreCerts}
}

type UserRepo struct {
	token       string
	client      mobile.UserAccessChecker
	host        string
	ignoreCerts bool
}

func (ur *UserRepo) GetUser() (*mobile.User, error) {
	return ur.client.ReadUserFromToken(ur.host, ur.token, ur.ignoreCerts)
}
