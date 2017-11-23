package cmd

import (
	"log"
	"net/http"
	"net/url"
	"path"

	"fmt"
	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createBuildCmd = &cobra.Command{
	Use:     "build",
	Short:   "instantiate a build",
	Aliases: []string{"b"},
	Run: func(cmd *cobra.Command, args []string) {
		httpclient := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s", cmd.Name(), err)
		}
		urlPath := fmt.Sprintf("/build/%s/instantiate", cmd.Flag("name").Value.String())
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

func init() {
	createBuildCmd.PersistentFlags().String("name", "", "name of the build")

	createCmd.AddCommand(createBuildCmd)
}
