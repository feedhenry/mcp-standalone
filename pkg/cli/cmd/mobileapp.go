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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/spf13/viper"

	"github.com/feedhenry/mcp-standalone/pkg/httpclient"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/spf13/cobra"
)

// mobileappCmd represents the mobileapp command
var getmobileappCmd = &cobra.Command{
	Use:     "mobileapp",
	Short:   "get mobile apps -- mcp get mobileapps -- mcp get mobileapp <appname>",
	Aliases: []string{"mobileapps"},
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		if len(args) > 0 {
			name = args[0]
		}
		var decodeInto interface{}
		httpclient := httpclient.NewHttpClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s ", cmd.Name(), err)
		}
		u.Path = path.Join(u.Path, "/mobileapp")
		if name != "" {
			u.Path = path.Join(u.Path, "/"+name)
		}
		fmt.Println(u.String())
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			log.Fatalf(" %s : failed to create new get request %s ", cmd.Name(), err)
		}
		addAuthorizationHeader(req.Header)
		res, err := httpclient.Do(req)
		if err != nil {
			log.Fatalf(" %s : failed to make the get request %s ", cmd.Name(), err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Fatalf("unexpected response code %v", res.StatusCode)
		}
		decoder := json.NewDecoder(res.Body)
		if name != "" {
			decodeInto = mobile.App{}
		} else {
			decodeInto = []*mobile.App{}
		}
		if err := decoder.Decode(&decodeInto); err != nil {
			log.Fatalf("%s : failed to decode response %s ", cmd.Name(), err)
		}
		data, err := json.MarshalIndent(&decodeInto, "", " ")
		if err != nil {
			log.Fatalf("%s : failed to marshal json %s ", cmd.Name(), err)
		}
		fmt.Println(string(data))
	},
}

func getMobileApp(name string) *mobile.App {
	var decodeInto *mobile.App
	httpclient := httpclient.NewHttpClientBuilder().Insecure(true).Timeout(5).Build()
	u, err := url.Parse(viper.GetString("host"))
	if err != nil {
		log.Fatalf("error parsing mcp host %s ", err)
	}
	u.Path = path.Join(u.Path, "/mobileapp")
	if name != "" {
		u.Path = path.Join(u.Path, "/"+name)
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatalf("failed to create new get request %s ", err)
	}
	addAuthorizationHeader(req.Header)
	res, err := httpclient.Do(req)
	if err != nil {
		log.Fatalf("failed to make the get request %s ", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("unexpected response code %v", res.StatusCode)
	}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&decodeInto); err != nil {
		log.Fatalf("failed to decode response %s ", err)
	}
	return decodeInto
}

var deletemobileCmd = &cobra.Command{
	Use:   "mobileapp",
	Short: "delete a mobile app -- mcp delete mobileapp <name>",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if name == "" {
			log.Fatal("name is required for delete mobileapp")
		}
		httpclient := httpclient.NewHttpClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s ", cmd.Name(), err)
		}
		u.Path = path.Join(u.Path, "/mobileapp/"+name)
		req, err := http.NewRequest("DELETE", u.String(), nil)
		if err != nil {
			log.Fatalf(" %s : failed to create new get request %s ", cmd.Name(), err)
		}
		addAuthorizationHeader(req.Header)
		res, err := httpclient.Do(req)
		if err != nil {
			log.Fatalf(" %s : failed to make the get request %s ", cmd.Name(), err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Fatalf("unexpected response code %v", res.StatusCode)
		}
	},
}

var createmobileCmd = &cobra.Command{
	Use:   "mobileapp",
	Short: "create a mobile app -- mcp create mobileapp <name> <clientType> | mcp create mobileapp -f ./app.json",
	Run: func(cmd *cobra.Command, args []string) {
		file := viper.GetString("file")
		var data []byte
		var err error
		var name, clientType string
		if file != "" {
			data, err = ioutil.ReadFile(file)
			if err != nil {
				log.Fatal("error reading app file ", err)
			}
		} else {
			if len(args) < 2 {
				log.Fatal(" if no file provided name and clientType are required")
			}
			name = args[0]
			clientType = args[1]
			app := &mobile.App{
				Name:       name,
				ClientType: clientType,
			}
			data, err = json.Marshal(app)
			if err != nil {
				log.Fatal("failed to prepare app json ", err)
			}
		}
		httpclient := httpclient.NewHttpClientBuilder().Insecure(true).Timeout(5).Build()
		u, err := url.Parse(viper.GetString("host"))
		if err != nil {
			log.Fatalf(" %s : error parsing mcp host %s ", cmd.Name(), err)
		}
		u.Path = path.Join(u.Path, "/mobileapp")
		req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
		if err != nil {
			log.Fatal("unexpected error creating post request ", err)
		}
		addAuthorizationHeader(req.Header)
		res, err := httpclient.Do(req)
		if err != nil {
			log.Fatal("unexpected error making post request ", err)
		}
		if res.StatusCode != http.StatusCreated {
			log.Fatalf("unexpected response code %v", res.StatusCode)
		}

	},
}

func init() {
	getCmd.AddCommand(getmobileappCmd)
	deleteCmd.AddCommand(deletemobileCmd)
	createCmd.AddCommand(createmobileCmd)
}
