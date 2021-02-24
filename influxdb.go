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
// TODO: Group fields IE Engine, GPS, ...
func InsertInflux(l *SafeCounter, c client.Client) {

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})
	if err != nil {
		log.Fatalln("Influxdb Insert Error: ", err)
	}

	// tm := (l.agg["GpsTime"]).(time.Time)
	tm := time.Now()

	for k, v := range l.agg {

		// fmt.Printf("k: %v - v: %v\n", k, v)

		// No need to store GpsTime in Measurements
		if k == "GpsTime" {
			continue
		}

		// Empty tag
		tags := map[string]string{}

		fields := map[string]interface{}{
			k: v,
		}
		// Make Point
		point, err := client.NewPoint(
			k,
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
	// Save points in DB
	err = c.Write(bp)
	if err != nil {
		log.Fatalln("Writing Points Error: ", err)
	}
}
