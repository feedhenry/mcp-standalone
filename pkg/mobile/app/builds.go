package app

import (
	"crypto/rand"
	"crypto/rsa"

	"bytes"
	"crypto/x509"
	"encoding/pem"

	"encoding/asn1"

	"io"

	"io/ioutil"

	"time"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Build struct {
	artifactRetriever mobile.ArtifactRetriever
	saToken           string
}

func NewBuild(artifactRetriever mobile.ArtifactRetriever, saToken string) *Build {
	return &Build{
		artifactRetriever: artifactRetriever,
		saToken:           saToken,
	}
}

type AppBuildCreatedResponse struct {
	PublicKey string `json:"publicKey"`
	BuildID   string `json:"buildID"`
}

func (b *Build) CreateAppBuild(buildRepo mobile.BuildCruder, build *mobile.BuildConfig) (*AppBuildCreatedResponse, error) {
	var res = &AppBuildCreatedResponse{BuildID: build.Name}
	if build.GitRepo.JenkinsFilePath == "" {
		build.GitRepo.JenkinsFilePath = "Jenkinsfile"
	}
	if !build.GitRepo.Private {
		if err := buildRepo.Create(build); err != nil {
			return nil, errors.Wrap(err, "CreateAppBuild: failed to create build")
		}
		return res, nil
	}
	assetName, pubkey, err := b.CreateBuildSrcKeySecret(buildRepo, build.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup src keys when creating app build ")
	}
	res.PublicKey = string(pubkey)
	build.GitRepo.PublicKeyID = assetName
	if err := buildRepo.Create(build); err != nil {
		return nil, errors.Wrap(err, "CreateAppBuild: failed to create build")
	}
	return res, nil

}

// CreateBuildSrcKeySecret creates a public private key pair and returns the secret name it is stored in and the public part of the key as bytes
func (b *Build) CreateBuildSrcKeySecret(br mobile.BuildCruder, buildName string) (string, []byte, error) {
	var (
		buildAsset    = mobile.BuildAsset{}
		privateKeyVal *bytes.Buffer
		publicKeyVal  *bytes.Buffer
		reader        = rand.Reader
		bitSize       = 2048
	)
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return "", nil, errors.Wrap(err, "CreateAppBuild: failed to generate a new rsa key pair when creating new app build")
	}
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	privateKeyVal = &bytes.Buffer{}
	err = pem.Encode(privateKeyVal, privateKey)
	if err != nil {
		return "", nil, errors.Wrap(err, "CreateAppBuild: failed to encode private key for new app build ")
	}
	pKey, err := asn1.Marshal(key.PublicKey)
	if err != nil {
		return "", nil, errors.Wrap(err, "CreateAppBuild: failed to marshal public key when creating app build")
	}
	var pubkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pKey,
	}
	publicKeyVal = &bytes.Buffer{}
	if err = pem.Encode(publicKeyVal, pubkey); err != nil {
		return "", nil, errors.Wrap(err, "CreateAppBuild: failed to encode public key for new app build")
	}
	buildAsset.BuildName = buildName
	buildAsset.Type = mobile.BuildAssetTypeSourceCredential
	buildAsset.Name = buildName + "ssh-key"
	buildAsset.AssetData = map[string][]byte{"ssh-privatekey": privateKeyVal.Bytes(), "ssh-publickey": publicKeyVal.Bytes()}
	assetName, err := br.AddBuildAsset(buildAsset)
	if err != nil {
		return "", nil, errors.Wrap(err, "CreateAppBuild: failed to add build asset ssh-key ")
	}
	return assetName, publicKeyVal.Bytes(), nil
}

// EnableDownload will enabling downloading of a build artefact for a set amount of time
func (b *Build) EnableDownload(br mobile.BuildCruder, buildName string) (*mobile.BuildDownload, error) {
	download := &mobile.BuildDownload{}
	token := uuid.NewV4().String()
	download.URL = "/build/" + buildName + "/download?token=" + token
	download.Expires = time.Now().Add(time.Minute * 30).Unix() //TODO this should be in config
	download.Token = token
	if err := br.AddDownload(buildName, download); err != nil {
		return nil, errors.Wrap(err, "enabling download failed when trying to add the download to the build")
	}
	return download, nil
}

func (b *Build) Download(br mobile.BuildCruder, buildName string) (io.ReadCloser, error) {
	buildStatus, err := br.Status(buildName)
	if err != nil {
		return nil, errors.Wrap(err, "download failed to get build status")
	}
	if buildStatus.Phase != "Complete" {
		return nil, &mobile.StatusError{Code: 404, Message: "no artifact found, build not completed yet. Build status : " + buildStatus.Phase}
	}
	u, err := buildStatus.ArtifactURL()
	if err != nil {
		return nil, errors.Wrap(err, "during download failed to get artifact url from build status")
	}
	return b.artifactRetriever.Retrieve(u, b.saToken)
}

func (b *Build) AddBuildAsset(br mobile.BuildCruder, resource io.Reader, asset *mobile.BuildAsset) (string, error) {
	data, err := ioutil.ReadAll(resource)
	if err != nil {
		return "", errors.Wrap(err, "failed to read the data for asset ")
	}
	if asset.Platform != "android" {
		return "", errors.New("android is the only supported build platform currently")
	}
	asset.AssetData = map[string][]byte{"p12": data, "password": []byte(asset.Password)}
	return br.AddBuildAsset(*asset)
}
