package cmd

import (
	json "encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"bytes"
	"fmt"
	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createBuildCmd = &cobra.Command{
	Use:     "clientbuild",
	Short:   "instantiate a build of a mobile client",
	Aliases: []string{"cb"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("expected app id and git url", "mcp create clientbuild myapp https://example-git-url.com")
		}
		clientBuildApp := args[0]
		clientBuildName := cmd.Flag("name").Value.String()
		clientBuildGitURL := args[1]
		if clientBuildName == "" {
			clientBuildName = clientBuildApp
		}

		httpclient := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s", cmd.Name(), err)
		}
		u.Path = path.Join(u.Path, "/build")
		buildconfig := &mobile.BuildConfig{
			AppID: clientBuildApp,
			Name:  clientBuildName,
			GitRepo: &mobile.BuildGitRepo{
				URI:             clientBuildGitURL,
				Ref:             cmd.Flag("gitref").Value.String(),
				JenkinsFilePath: cmd.Flag("jenkinsfile").Value.String(),
			},
		}
		data, err := json.Marshal(buildconfig)
		if err != nil {
			log.Fatal("failed to prepare buildconfig JSON", err)
		}
		req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
		if err != nil {
			log.Fatal("failed to create POST request", err)
		}
		addAuthorizationHeader(req.Header)
		res, err := httpclient.Do(req)
		if err != nil {
			log.Fatal("failed to perform POST request", err)
		}
		if res.StatusCode != http.StatusCreated {
			log.Fatalf("unexpected response code %v", res.StatusCode)
		}
	},
}

var startBuildCmd = &cobra.Command{
	Use:     "clientbuild",
	Short:   "instantiate a build of a mobile client",
	Aliases: []string{"cb"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("expected client build name", "mcp start clientbuild myapp")
		}
		clientBuildName := args[0]
		httpclient := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s", cmd.Name(), err)
		}
		urlPath := fmt.Sprintf("/build/%s/instantiate", clientBuildName)
		u.Path = path.Join(u.Path, urlPath)
		if err != nil {
			log.Fatal("failed to prepare buildconfig JSON", err)
		}
		req, err := http.NewRequest("POST", u.String(), nil)
		if err != nil {
			log.Fatal("failed to create POST request", err)
		}
		addAuthorizationHeader(req.Header)
		res, err := httpclient.Do(req)
		if err != nil {
			log.Fatal("failed to perform POST request", err)
		}
		if res.StatusCode != http.StatusCreated {
			log.Fatalf("unexpected response code %v", res.StatusCode)
		}
	},
}

func initCreateBuildCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("appid", "", "set the id of the app to build")
	cmd.PersistentFlags().String("name", "", "set a unique name for the build")
	cmd.PersistentFlags().String("giturl", "", "set the url of the git repository to build from")
	cmd.PersistentFlags().String("gitref", "master", "set the git ref to build from")
	cmd.PersistentFlags().String("jenkinsfile", "Jenkinsfile", "set the path of the Jenkinsfile in the repo")
}

func initStartBuildCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("name", "", "name of the build")
}

func init() {
	initCreateBuildCmd(createBuildCmd)
	initStartBuildCmd(startBuildCmd)

	startCmd.AddCommand(startBuildCmd)
	createCmd.AddCommand(createBuildCmd)
}
