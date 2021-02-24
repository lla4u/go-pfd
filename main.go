package main

import (
	"flag"
	"net"
	"os"
	"sync"
	"time"

	"github.com/brutella/can"
	log "github.com/sirupsen/logrus"
)

// SafeCounter struc to store measurements such as Engine, AHRS and GPS
// sync Mutex is used to prevent map simultaneous read and write.
type SafeCounter struct {
	mu  sync.Mutex
	agg map[string]interface{}
}

// Global variables
var verbose *bool // Flag verbose mode
var frequency int // Flag sampling frequency
var gps *bool     // Flag to enable external Gps

var sc SafeCounter = SafeCounter{agg: make(map[string]interface{})} // Cumulative measurements

func main() {

	// Verbose mode flag
	verbose = flag.Bool("v", false, "Run in verbose mode")

	// Sampling frequency flag
	sampling := flag.Int("s", 1, "Sampling rate is second")

	// Verbose mode flag
	gps = flag.Bool("gps", false, "Run with external gps (gpsd)")

	flag.Parse()

	// set frequency from flag parameter
	frequency = *sampling

	// Log file config

	if *verbose {
		// Log format config
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true})

		// Log output on stdout
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		log.WithFields(log.Fields{"verbose": *verbose}).Info("Flag verbose mode enabled")

	} else {
		logFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		// SetLevel can be invoqued only once
		log.SetOutput(logFile)
		log.SetLevel(log.WarnLevel)
	}

	log.WithFields(log.Fields{"frequency": frequency}).Info("Flag sampling frequency")
	log.WithFields(log.Fields{"gps": *gps}).Info("Flag external GPS mode enabled")

	// Infludb client setup
	log.Info("Starting Database ...")
	client := influxDBClient()
	defer client.Close()

	// Goroutine to flush into db at frequency
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	go func() {
		for range ticker.C {
			sc.mu.Lock()

			if *verbose {
				log.WithFields(log.Fields{"measurements": sc.agg}).Info("Measurements")
			}

			InsertInflux(&sc, client)

			// Flush map for next round
			sc.agg = make(map[string]interface{})
			sc.mu.Unlock()
		}
	}()

	// CAN setup starting by interface
	log.Info("Starting Can ...")
	iface, err := net.InterfaceByName("can0")

	if err != nil {
		log.Fatalf("Could not find network interface %s (%v)", "can0", err)
	}

	// CAN open RW
	conn, err := can.NewReadWriteCloserForInterface(iface)

	if err != nil {
		log.Fatal(err)
	}

	// CAN bus init
	canBus := can.NewBus(conn)

	canBus.SubscribeFunc(logDakuFrame)

	canBus.ConnectAndPublish()
}
