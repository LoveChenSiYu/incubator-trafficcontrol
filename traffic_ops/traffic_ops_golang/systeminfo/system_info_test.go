package systeminfo

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"encoding/json"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"
)

func TestGetSystemInfo(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := test.ColsFromStructByTag("db", tc.ParameterNullable{})
	rows := sqlmock.NewRows(cols)

	lastUpdated := tc.TimeNoMod{Time: time.Now()}
	configFile := "global"
	secure := false

	firstID := 1
	firstName := "paramname1"
	firstVal := "val1"

	secondID := 2
	secondName := "paramname2"
	secondVal := "val2"

	var sysInfoParameters = []tc.ParameterNullable{
		tc.ParameterNullable{
			ConfigFile:  &configFile,
			ID:          &firstID,
			LastUpdated: &lastUpdated,
			Name:        &firstName,
			Profiles:    json.RawMessage(`["foo","bar"]`),
			Secure:      &secure,
			Value:       &firstVal,
		},

		tc.ParameterNullable{
			ConfigFile:  &configFile,
			ID:          &secondID,
			LastUpdated: &lastUpdated,
			Name:        &secondName,
			Profiles:    json.RawMessage(`["foo","bar"]`),
			Secure:      &secure,
			Value:       &secondVal,
		},
	}

	for _, ts := range sysInfoParameters {
		rows = rows.AddRow(
			ts.ConfigFile,
			ts.ID,
			ts.LastUpdated,
			ts.Name,
			ts.Profiles,
			ts.Secure,
			ts.Value,
		)
	}

	mock.ExpectQuery("SELECT.*WHERE p.config_file='global'").WillReturnRows(rows)
	sysinfo, err := getSystemInfo(db, auth.PrivLevelReadOnly)
	if err != nil {
		t.Errorf("getSystemInfo expected: nil error, actual: %v", err)
	}

	if len(sysinfo) != 2 {
		t.Errorf("getSystemInfo expected: len(sysinfo) == 2, actual: %v", len(sysinfo))
	}
}
