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
	"github.com/rs/zerolog/log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// RecordWriter is a simple wrapper around InfluxDB client allowing
// adding records in a convenient way.
type RecordWriter[T Influxable] struct {
	db *InfluxDBAdapter
}

// AddRecord adds a record to a respective InfluxDB database
// and measurement.
func (c *RecordWriter[T]) AddRecord(measurement string, rec T) {
	tags, values := rec.ToInfluxDB()
	p := influxdb2.NewPointWithMeasurement(measurement)
	p.SetTime(rec.GetTime())
	for tn, tv := range tags {
		p.AddTag(tn, tv)
	}
	for field, value := range values {
		p.AddField(field, value)
	}
	c.db.WritePoint(p)
}

// NewRecordWriter is a factory function for RecordWriter
func NewRecordWriter[T Influxable](db *InfluxDBAdapter) *RecordWriter[T] {
	return &RecordWriter[T]{db}
}

// RunWriteConsumerSync reads from incomingData channel and stores the data
// to via a provided InfluxDBAdapter ('db' arg.). In case 'db' is nil, the
// function just listens to 'incomingData' and does nothing.
// Typically, this function should run in its own goroutine.
func RunWriteConsumerSync[T Influxable](
	db *InfluxDBAdapter,
	measurement string,
	incomingData <-chan T,
) {
	if db != nil {
		var err error
		client := NewRecordWriter[T](db)
		for rec := range incomingData {
			client.AddRecord(measurement, rec)
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to write influxDB record")
		}

	} else {
		for range incomingData {
		}
	}
}
