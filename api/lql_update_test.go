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

package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/stretchr/testify/assert"
)

var (
	updateQueryText = `my_lql { source { CloudTrailRawEvents } return { INSERT_ID, INSERT_TIME } }`
	updateQuery     = api.UpdateQuery{
		QueryText: newQueryText,
	}
	updateQueryJSON = fmt.Sprintf(`{
	"evaluatorId": "%s",
	"queryId": "%s",
	"queryText": "%s"
}`, queryEvaluator, queryID, updateQueryText)
)

func TestQueryUpdateMethod(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method, "Update should be a PATCH method")
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Update(queryID, updateQuery)
	assert.Nil(t, err)
}

func TestQueryUpdateBadInput(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{}")
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Update("", api.UpdateQuery{})
	assert.Equal(t, "query ID must be provided", err.Error())
}

func TestQueryUpdateOK(t *testing.T) {
	mockResponse := mockQueryDataResponse(updateQueryJSON)

	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResponse)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	updateExpected := api.QueryResponse{}
	_ = json.Unmarshal([]byte(mockResponse), &updateExpected)

	var updateActual api.QueryResponse
	updateActual, err = c.V2.Query.Update(queryID, updateQuery)
	assert.Nil(t, err)

	assert.Equal(t, updateExpected, updateActual)
}

func TestQueryUpdateNotFound(t *testing.T) {
	fakeServer := lacework.MockServer()
	fakeServer.UseApiV2()
	fakeServer.MockAPI(
		"Queries/my_lql",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, lqlErrorReponse, http.StatusBadRequest)
		},
	)
	defer fakeServer.Close()

	c, err := api.NewClient("test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)
	assert.Nil(t, err)

	_, err = c.V2.Query.Update(queryID, updateQuery)
	assert.NotNil(t, err)
}
