//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
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
//

package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Profiles is the representation of the ~/.lacework.toml
//
// Example:
//
// [default]
// account = "example"
// api_key = "EXAMPLE_0123456789"
// api_secret = "_0123456789"
//
// [dev]
// account = "dev"
// api_key = "DEV_0123456789"
// api_secret = "_0123456789"
type Profiles map[string]credsDetails

type credsDetails struct {
	Account   string `toml:"account" json:"account"`
	ApiKey    string `toml:"api_key" json:"api_key" survey:"api_key"`
	ApiSecret string `toml:"api_secret" json:"api_secret" survey:"api_secret"`
}

// apiKeyDetails represents the details of an API key, we use this struct
// internally to unmarshal the JSON file provided by the Lacework WebUI
type apiKeyDetails struct {
	KeyID  string `json:"keyId"`
	Secret string `json:"secret"`
}

var (
	// configureJsonFile is the API key file downloaded form the Lacework WebUI
	configureJsonFile string

	// configureCmd represents the configure command
	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "configure the Lacework CLI",
		Args:  cobra.NoArgs,
		Long: `
Configure settings that the Lacework CLI uses to interact with the Lacework
platform. These include your Lacework account, API access key and secret.

To create a set of API keys, log in to your Lacework account via WebUI and
navigate to Settings > API Keys and click + Create New. Enter a name for
the key and an optional description, then click Save. To get the secret key,
download the generated API key file.

Use the argument --json_file to preload the downloaded API key file.

If this command is run with no arguments, the Lacework CLI will store all
settings under the default profile. The information in the default profile
is used any time you run a Lacework CLI command that doesn't explicitly
specify a profile to use.

You can configure multiple profiles by using the --profile argument. If a
config file does not exist (the default location is ~/.lacework.toml), the
Lacework CLI will create it for you.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return promptConfigureSetup()
		},
	}
)

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&configureJsonFile,
		"json_file", "j", "", "loads the generated API key JSON file from the WebUI",
	)
}

func promptConfigureSetup() error {
	cli.Log.Debugw("configuring cli", "profile", cli.Profile)

	if len(configureJsonFile) != 0 {
		auth, err := loadKeysFromJsonFile(configureJsonFile)
		if err != nil {
			return errors.Wrap(err, "unable to load keys from the provided json file")
		}
		cli.KeyID = auth.KeyID
		cli.Secret = auth.Secret
	}

	questions := []*survey.Question{
		{
			Name: "account",
			Prompt: &survey.Input{
				Message: "Account: ",
				Default: cli.Account,
			},
			Validate: promptRequiredStringLen(0,
				"The account subdomain of URL is required. (i.e. <ACCOUNT>.lacework.net)",
			),
		},
		{
			Name: "api_key",
			Prompt: &survey.Input{
				Message: "Access Key ID: ",
				Default: cli.KeyID,
			},
			Validate: promptRequiredStringLen(55,
				"The API access key id must have more than 55 characters.",
			),
		},
		{
			Name: "api_secret",
			Prompt: &survey.Input{
				Message: "Secret Access Key: ",
				Default: cli.Secret,
			},
			Validate: promptRequiredStringLen(30,
				"The API secret access key must have more than 30 characters.",
			),
		},
	}

	newCreds := credsDetails{}
	err := survey.Ask(questions, &newCreds,
		survey.WithIcons(promptIconsFunc),
	)
	if err != nil {
		return err
	}

	var (
		profiles = Profiles{}
		buf      = new(bytes.Buffer)
		confPath = viper.ConfigFileUsed()
	)

	if confPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		confPath = path.Join(home, ".lacework.toml")
		cli.Log.Debugw("generating new config file",
			"path", confPath,
		)
	} else {
		if _, err := toml.DecodeFile(confPath, &profiles); err != nil {
			cli.Log.Debugw("unable to decode profiles from config, trying previous config",
				"path", confPath, "error", err,
			)

			var oldcreds credsDetails
			if _, err2 := toml.DecodeFile(confPath, &oldcreds); err2 != nil {
				cli.Log.Debugw("unable to decode old config, no more options, exit",
					"error", err2,
				)
				return err
			}
			profiles["default"] = oldcreds
		}
		cli.Log.Debugw("profiles loaded from config, updating", "profiles", profiles)
	}

	profiles[cli.Profile] = newCreds
	cli.Log.Debugw("storing updated profiles", "profiles", profiles)
	if err := toml.NewEncoder(buf).Encode(profiles); err != nil {
		return err
	}

	err = ioutil.WriteFile(confPath, buf.Bytes(), 0600)
	if err != nil {
		return err
	}

	cli.OutputHuman("\nYou are all set!\n")
	return nil
}

func loadKeysFromJsonFile(file string) (*apiKeyDetails, error) {
	cli.Log.Debugw("loading API key JSON file", "path", file)
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cli.Log.Debugw("keys from file", "raw", string(jsonData))
	var auth apiKeyDetails
	err = json.Unmarshal(jsonData, &auth)
	return &auth, err
}
