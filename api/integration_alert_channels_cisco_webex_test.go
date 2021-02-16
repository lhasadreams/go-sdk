//
// Author:: Darren Murray (<darren.murray@lacework.net>)
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

package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/intgguid"
	"github.com/lacework/go-sdk/internal/lacework"
)

func TestIntegrationsNewCiscoWebexAlertChannel(t *testing.T) {
	subject := api.NewCiscoWebexAlertChannel("integration_name",
		api.CiscoWebexChannelData{
			WebhookURL: "https://webexapis.com/v1/webhooks/incoming/api-token",
		},
	)
	assert.Equal(t, api.CiscoWebexChannelIntegration.String(), subject.Type)
}

func TestIntegrationsCreateCiscoWebexAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "CreateCiscoWebexAlertChannel should be a POST method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "CISCO_SPARK_WEBHOOK", "wrong integration type")
			assert.Contains(t, body, "https://webexapis.com/v1/webhooks/incoming/api-token", "wrong webhook url")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, ciscoWebexChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewCiscoWebexAlertChannel("integration_name",
		api.CiscoWebexChannelData{
			WebhookURL: "https://webexapis.com/v1/webhooks/incoming/api-token",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "CiscoWebex integration name mismatch")
	assert.Equal(t, "CISCO_SPARK_WEBHOOK", data.Type, "a new CiscoWebex integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new CiscoWebex integration should be enabled")

	response, err := c.Integrations.CreateCiscoWebexAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "https://webexapis.com/v1/webhooks/incoming/api-token", resData.Data.WebhookURL)
	}
}

func TestIntegrationsGetCiscoWebexAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "GetCiscoWebexAlertChannel should be a GET method")
		fmt.Fprintf(w, ciscoWebexChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.GetCiscoWebexAlertChannel(intgGUID)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	if assert.Equal(t, 1, len(response.Data)) {
		resData := response.Data[0]
		assert.Equal(t, intgGUID, resData.IntgGuid)
		assert.Equal(t, "integration_name", resData.Name)
		assert.True(t, resData.State.Ok)
		assert.Equal(t, "https://webexapis.com/v1/webhooks/incoming/api-token", resData.Data.WebhookURL)
	}
}

func TestIntegrationsUpdateCiscoWebexAlertChannel(t *testing.T) {
	var (
		intgGUID   = intgguid.New()
		apiPath    = fmt.Sprintf("external/integrations/%s", intgGUID)
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI(apiPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method, "UpdateCiscoWebexAlertChannel should be a PATCH method")

		if assert.NotNil(t, r.Body) {
			body := httpBodySniffer(r)
			assert.Contains(t, body, intgGUID, "INTG_GUID missing")
			assert.Contains(t, body, "integration_name", "integration name is missing")
			assert.Contains(t, body, "CISCO_SPARK_WEBHOOK", "wrong integration type")
			assert.Contains(t, body, "https://webexapis.com/v1/webhooks/incoming/api-token", "wrong webhook url")
			assert.Contains(t, body, "ENABLED\":1", "integration is not enabled")
		}

		fmt.Fprintf(w, ciscoWebexChannelIntegrationJsonResponse(intgGUID))
	})
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	data := api.NewCiscoWebexAlertChannel("integration_name",
		api.CiscoWebexChannelData{
			WebhookURL: "https://webexapis.com/v1/webhooks/incoming/api-token",
		},
	)
	assert.Equal(t, "integration_name", data.Name, "CiscoWebex integration name mismatch")
	assert.Equal(t, "CISCO_SPARK_WEBHOOK", data.Type, "a new CiscoWebex integration should match its type")
	assert.Equal(t, 1, data.Enabled, "a new CiscoWebex integration should be enabled")
	data.IntgGuid = intgGUID

	response, err := c.Integrations.UpdateCiscoWebexAlertChannel(data)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SUCCESS", response.Message)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, intgGUID, response.Data[0].IntgGuid)
}

func TestIntegrationsListCiscoWebexAlertChannel(t *testing.T) {
	var (
		intgGUIDs  = []string{intgguid.New(), intgguid.New(), intgguid.New()}
		fakeServer = lacework.MockServer()
	)
	fakeServer.MockAPI("external/integrations/type/CISCO_SPARK_WEBHOOK",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "ListCiscoWebexAlertChannel should be a GET method")
			fmt.Fprintf(w, ciscoWebexChanMultiIntegrationJsonResponse(intgGUIDs))
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	response, err := c.Integrations.ListCiscoWebexAlertChannel()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Ok)
	assert.Equal(t, len(intgGUIDs), len(response.Data))
	for _, d := range response.Data {
		assert.Contains(t, intgGUIDs, d.IntgGuid)
	}
}

func ciscoWebexChannelIntegrationJsonResponse(intgGUID string) string {
	return `
{
  "data": [` + singleCiscoWebexChanIntegration(intgGUID) + `],
  "ok": true,
  "message": "SUCCESS"
}
`
}

func ciscoWebexChanMultiIntegrationJsonResponse(guids []string) string {
	integrations := []string{}
	for _, guid := range guids {
		integrations = append(integrations, singleCiscoWebexChanIntegration(guid))
	}
	return `
{
"data": [` + strings.Join(integrations, ", ") + `],
"ok": true,
"message": "SUCCESS"
}
`
}

func singleCiscoWebexChanIntegration(id string) string {
	return `
{
  "INTG_GUID": "` + id + `",
  "CREATED_OR_UPDATED_BY": "user@email.com",
  "CREATED_OR_UPDATED_TIME": "2020-Jul-16 19:59:22 UTC",
  "DATA": {
    "WEBHOOK": "https://webexapis.com/v1/webhooks/incoming/api-token"
  },
  "ENABLED": 1,
  "IS_ORG": 0,
  "NAME": "integration_name",
  "STATE": {
    "lastSuccessfulTime": "2020-Jul-16 18:26:54 UTC",
    "lastUpdatedTime": "2020-Jul-16 18:26:54 UTC",
    "ok": true
  },
  "TYPE": "CISCO_SPARK_WEBHOOK",
  "TYPE_NAME": "CISCO_SPARK_WEBHOOK"
}
`
}
