// Copyright 2022 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2022 Martin Zimandl <martin.zimandl@gmail.com>
// Copyright 2022 Charles University - Faculty of Arts,
//                Institute of the Czech National Corpus
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package influx

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// ConnectionConf specifies a configuration required to store data
// to an InfluxDB database
type ConnectionConf struct {
	Server       string `json:"server"`
	Token        string `json:"token"`
	Organization string `json:"organization"`
	Bucket       string `json:"bucket"`
}

// IsConfigured tests whether the configuration is considered
// to be enabled (i.e. no error checking just enabled/disabled)
func (conf *ConnectionConf) IsConfigured() bool {
	return conf.Server != ""
}

// ------

// Influxable represents any type which is able
// to export its data in a format required by InfluxDB.
type Influxable interface {

	// ToInfluxDB defines a method providing data
	// to be written to a database. The first returned
	// value is for tags, the second one for fields.
	ToInfluxDB() (map[string]string, map[string]any)

	// GetTime provides a date and time when the record
	// was created.
	GetTime() time.Time
}

type influxDBAdapter struct {
	api             influxdb2api.WriteAPI
	address         string
	onErrorHandlers []func(error)
}

func (db *influxDBAdapter) WritePoint(p *write.Point) {
	db.api.WritePoint(p)
}

func (db *influxDBAdapter) Address() string {
	return db.address
}

func ConnectAPI(conf *ConnectionConf, errListen <-chan error) *influxDBAdapter {
	ans := new(influxDBAdapter)
	ans.onErrorHandlers = make([]func(error), 0, 10)
	var influxClient influxdb2.Client
	if conf.IsConfigured() {
		ans.address = conf.Server
		influxClient = influxdb2.NewClient(conf.Server, conf.Token)
		ans.api = influxClient.WriteAPI(
			conf.Organization,
			conf.Bucket,
		)
		go func() {
			rt := make(chan error)
			for err := range ans.api.Errors() {
				rt <- err
			}
			close(rt)
		}()
		return ans
	}
	return nil
}
