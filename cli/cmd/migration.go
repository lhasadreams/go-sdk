//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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
	"regexp"
	"strings"

	"github.com/lacework/go-sdk/lwconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// MigrateProfile does an automatic migration of a single profile,
// which means that the user needs to run a single command on every
// profile to trigger the migration
func (c *cliState) MigrateProfile() (err error) {
	c.Event.Feature = featMigrateConfigV2
	defer func() {
		if err == nil {
			c.SendHoneyvent()
		}
	}()

	//c.Event.AddFeatureField("os", osInfo.Name)
	err = c.VerifySettings()
	if err != nil {
		return err
	}

	account, errOrgInfo := c.LwApi.Account.GetOrganizationInfo()
	if errOrgInfo != nil {
		c.Event.Error = err.Error()
		return err
	}

	// TODO @afiune what if there is no config file?
	// what if the user send the required settings via Env Variables or Flags?
	migratedProfile := lwconfig.Profile{
		Account:   c.Account,
		ApiKey:    c.KeyID,
		ApiSecret: c.Secret,
		Version:   2,
	}

	if account.OrgAccount {

		// substract the account name from the full domain ACCOUNT.lacework.net
		rx, err := regexp.Compile(`\.lacework\.net.*`)
		if err != nil {
			return errors.Wrap(err, "unable to substract account name from full domain")
		}
		accountSplit := rx.Split(strings.ToLower(account.OrgAccountURL), -1)
		if len(accountSplit) != 0 {
			// set the right organization account name
			migratedProfile.Account = accountSplit[0]
		}

		// if the user is accessing a sub-account, that is, if the previous
		// account is different from the organizational account name, set it
		if migratedProfile.Account != strings.ToLower(c.Account) {
			migratedProfile.Subaccount = c.Account
		}
	}

	if err := lwconfig.StoreProfileAt(viper.ConfigFileUsed(), cli.Profile, migratedProfile); err != nil {
		return errors.Wrap(err, "unable to configure the command-line")
	}

	// set config version = 2

	// make a backup

	// store it

	return nil
}
