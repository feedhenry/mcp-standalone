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
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		token := viper.GetString("token")
		if host == "" || token == "" {
			log.Fatal("--host and --token need to be set")
		}
		config := map[string]string{
			"host":  host,
			"token": token,
		}
		data, err := json.Marshal(config)
		if err != nil {
			log.Fatal("failed to marshal the config ", err)
		}
		if err := ioutil.WriteFile(viper.ConfigFileUsed(), data, os.ModePerm); err != nil {
			log.Fatal("failed to write config file at locatation "+viper.ConfigFileUsed(), err)
		}

	},
}

func init() {
	RootCmd.AddCommand(targetCmd)
}
