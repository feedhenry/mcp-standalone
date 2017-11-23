package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createBuildConfigCmd = &cobra.Command{
	Use:     "buildconfig",
	Short:   "",
	Aliases: []string{"bc"},
	Run: func(cmd *cobra.Command, args []string) {
		httpclient := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s", cmd.Name(), err)
		}
		u.Path = path.Join(u.Path, "/build")
		buildconfig := &mobile.BuildConfig{
			AppID: cmd.Flag("host").Value.String(),
			Name:  cmd.Flag("name").Value.String(),
			GitRepo: &mobile.BuildGitRepo{
				URI:             cmd.Flag("giturl").Value.String(),
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

func init() {
	createBuildConfigCmd.PersistentFlags().String("appid", "", "set the id of the app to build")
	createBuildConfigCmd.PersistentFlags().String("name", "", "set a unique name for the build")
	createBuildConfigCmd.PersistentFlags().String("giturl", "", "set the url of the git repository to build from")
	createBuildConfigCmd.PersistentFlags().String("gitref", "master", "set the git ref to build from")
	createBuildConfigCmd.PersistentFlags().String("jenkinsfile", "Jenkinsfile", "set the path of the Jenkinsfile in the repo")

	createCmd.AddCommand(createBuildConfigCmd)
}
