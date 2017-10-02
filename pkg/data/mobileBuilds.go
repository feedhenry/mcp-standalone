package data

import (
	"fmt"

	"strconv"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/build"
	"github.com/feedhenry/mcp-standalone/pkg/openshift/client"
	"github.com/pkg/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/v1"
	kapi "k8s.io/client-go/pkg/api/v1"
)

type BuildRepo struct {
	buildConfClient client.BuildConfigInterface
	buildClient     client.BuildInterface
	secretClient    corev1.SecretInterface
	validator       MobileBuildValidator
}

func NewBuildRepo(bc client.BuildConfigInterface, buildc client.BuildInterface, sc corev1.SecretInterface) *BuildRepo {
	br := &BuildRepo{
		buildConfClient: bc,
		secretClient:    sc,
		buildClient:     buildc,
	}
	br.validator = DefaultMobileBuildValidator{}
	return br
}

func (br *BuildRepo) Create(b *mobile.BuildConfig) error {
	if err := br.validator.PreCreate(b); err != nil {
		return errors.Wrap(err, "failed to create build as validation failed")
	}
	bc := convertMobileBuildToBuildConfig(b)
	if _, err := br.buildConfClient.Create(bc); err != nil {
		return errors.Wrap(err, "build repo failed to create underlying buildconfig")
	}
	return nil
}

func (br BuildRepo) AddDownload(buildName string, dl *mobile.BuildDownload) error {
	build, err := br.buildClient.Get(buildName, meta_v1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to add download settings to build")
	}
	fmt.Println("got build", build)
	build.Annotations["downloadURL"] = dl.URL
	build.Annotations["downloadExp"] = fmt.Sprintf("%d", dl.Expires)
	build.Annotations["downloadToken"] = dl.Token
	if _, err := br.buildClient.Update(build); err != nil {
		return errors.Wrap(err, "failed to update the build with the download settings ")
	}
	return nil
}

func (br *BuildRepo) GetDownload(buildName string) (*mobile.BuildDownload, error) {
	build, err := br.buildClient.Get(buildName, meta_v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add download settings to build")
	}
	downLoadURL := build.Annotations["downloadURL"]
	expire := build.Annotations["downloadExp"]
	token := build.Annotations["downloadToken"]
	if downLoadURL == "" || expire == "" {
		return nil, errors.New("no download settings available")
	}
	expireVal, err := strconv.ParseInt(expire, 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse expire value for download")
	}
	return &mobile.BuildDownload{
		Expires: expireVal,
		URL:     downLoadURL,
		Token:   token,
	}, nil
}

// AddBuildAsset will create a secret and return its name
func (br *BuildRepo) AddBuildAsset(asset mobile.BuildAsset) (string, error) {
	labels := map[string]string{"type": string(asset.Type), "platform": asset.Platform}
	if asset.BuildName != "" {
		labels["buildID"] = asset.BuildName
	}
	if asset.AppName != "" {
		labels["mobileAppID"] = asset.AppName
	}
	secret := &v1.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels: labels,
			Name:   asset.Name,
		},
		Data: asset.AssetData,
	}
	s, err := br.secretClient.Create(secret)
	if err != nil {
		return "", errors.Wrap(err, "failed to add build asset")
	}
	return s.Name, nil

}

func (br *BuildRepo) Update(config *build.BuildConfig) (*build.BuildConfig, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func convertMobileBuildToBuildConfig(b *mobile.BuildConfig) *build.BuildConfig {
	bc := &build.BuildConfig{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels: map[string]string{"mobile-appid": b.AppID},
			Name:   b.Name,
		},
		Spec: build.BuildConfigSpec{
			CommonSpec: build.CommonSpec{
				Strategy: build.BuildStrategy{
					Type: build.JenkinsPipelineBuildStrategyType,
					JenkinsPipelineStrategy: &build.JenkinsPipelineBuildStrategy{

						JenkinsfilePath: b.GitRepo.JenkinsFilePath,
					},
				},
				Source: build.BuildSource{
					Git: &build.GitBuildSource{

						URI: b.GitRepo.URI,
						Ref: b.GitRepo.Ref,
					},
				},
			},
		},
	}

	if b.GitRepo.Private {
		bc.Spec.Source.SourceSecret = &kapi.LocalObjectReference{
			Name: b.GitRepo.PublicKeyID,
		}
	}
	return bc
}

func convertBuildConfigToMobileBuild(b *build.BuildConfig) (*mobile.BuildConfig, error) {
	return nil, nil
}

// MobileBuildValidator defines what a validator should do
type MobileBuildValidator interface {
	PreCreate(a *mobile.BuildConfig) error
	PreUpdate(old *mobile.BuildConfig, new *mobile.BuildConfig) error
}

// NewBuildsRepoBuilder provides an implementation of mobile.ServiceRepoBuilder
func NewBuildsRepoBuilder(clientBuilder mobile.K8ClientBuilder, ocClientBuilder mobile.OSClientBuilder, namespace, saToken string) mobile.BuildRepoBuilder {
	return &BuildRepoBuilder{
		clientBuilder:   clientBuilder,
		ocClientBuilder: ocClientBuilder,
		saToken:         saToken,
		namespace:       namespace,
	}
}

type BuildRepoBuilder struct {
	clientBuilder   mobile.K8ClientBuilder
	ocClientBuilder mobile.OSClientBuilder
	token           string
	namespace       string
	saToken         string
}

func (marb *BuildRepoBuilder) WithToken(token string) mobile.BuildRepoBuilder {
	return &BuildRepoBuilder{
		clientBuilder:   marb.clientBuilder,
		ocClientBuilder: marb.ocClientBuilder,
		token:           token,
		saToken:         marb.saToken,
		namespace:       marb.namespace,
	}
}

//UseDefaultSAToken delegates off to the service account token setup with the MCP. This should only be used for APIs where no real token is provided and should always be protected
func (marb *BuildRepoBuilder) UseDefaultSAToken() mobile.BuildRepoBuilder {
	return &BuildRepoBuilder{
		clientBuilder:   marb.clientBuilder,
		ocClientBuilder: marb.ocClientBuilder,
		token:           marb.saToken,
		saToken:         marb.saToken,
		namespace:       marb.namespace,
	}

}

// Build builds the final repo
func (marb *BuildRepoBuilder) Build() (mobile.BuildCruder, error) {
	k8client, err := marb.clientBuilder.WithToken(marb.token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "MobileAppRepoBuilder failed to build a configmap client")
	}
	ocClient, err := marb.ocClientBuilder.WithToken(marb.token).BuildClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build ocClient ")
	}
	return NewBuildRepo(ocClient.BuildConfigs(marb.namespace), ocClient.Builds(marb.namespace), k8client.CoreV1().Secrets(marb.namespace)), nil
}
