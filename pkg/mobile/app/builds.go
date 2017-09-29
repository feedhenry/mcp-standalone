package app

import (
	"crypto/rand"
	"crypto/rsa"

	"bytes"
	"crypto/x509"
	"encoding/pem"

	"encoding/asn1"

	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

type Build struct {
}

func NewBuild() *Build {
	return &Build{}
}

type AppBuildCreatedResponse struct {
	PublicKey string `json:"publicKey"`
	BuildID   string `json:"buildID"`
}

func (b *Build) CreateAppBuild(buildRepo mobile.BuildCruder, build *mobile.Build) (*AppBuildCreatedResponse, error) {
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
