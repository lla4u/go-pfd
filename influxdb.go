package main

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

// Influxdb
const (
	database = "BBOX"
	username = "admin"
	password = "generic"
)

// influxDBClient function
func influxDBClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return c
}

// InsertInflux function
// Todo souhaitable de sortir log Fatal pour log Error ?
func InsertInflux(l *SafeCounter, c client.Client) {

	// records types are made of ENGINE, AIRU, GPS types.
	groups := [4]string{"ENGINE", "AIRU", "GPS", "Unknown"}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})
	if err != nil {
		log.Fatalln("Influxdb Insert Error: ", err)
	}

	// tm := (l.agg["GpsTime"]).(time.Time)
	tm := time.Now()

	for _, t := range groups {

		tags := map[string]string{}

		fields := map[string]interface{}{}

		for k, v := range l.agg {

			if v.Type == t {
				fields[k] = v.Value
			}
		}
		if len(fields) > 0 {
			// Make Point
			point, err := client.NewPoint(
				t,
				tags,
				fields,
				tm,
			)
			if err != nil {
				log.Fatalln("Adding Point Error: ", err)
			}
			// Add Point in Batch
			bp.AddPoint(point)
		}
	}
	// Save points in DB
	err = c.Write(bp)
	if err != nil {
		log.Fatalln("Writing Points Error: ", err)
	}
}
