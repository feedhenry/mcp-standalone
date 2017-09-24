// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sdkconfigCmd represents the sdkconfig command
var getsdkconfigCmd = &cobra.Command{
	Use:   "sdkconfig",
	Short: "will get the current sdk config for a mobile app",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("expected a mobile app name", "mcp get sdkconfig myapp")
		}
		appName := args[0]
		app := getMobileApp(appName)
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf("error parsing mcp host %s ", err)
		}
		u.Path = path.Join(u.Path, fmt.Sprintf("/sdk/mobileapp/%s/config", app.ID))
		fmt.Println("url is ", u.String())
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			log.Fatalf("error creating request %s ", err)
		}
		req.Header.Set(mobile.AppAPIKeyHeader, app.APIKey)
		httpclient := httpclient.NewClientBuilder().Insecure(true).Timeout(5).Build()
		res, err := httpclient.Do(req)
		if err != nil {
			log.Fatalf("error doint request %s ", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatal("unexpected status code response " + res.Status)
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("error reading request body %s ", err)
		}
		fmt.Println(string(data))

	},
}

func init() {
	getCmd.AddCommand(getsdkconfigCmd)
}
