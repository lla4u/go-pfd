package main

import (
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
)

// temporary directory location
var backupDir = filepath.FromSlash("/tmp/")

func portal() {

	// return a `.tgz` file for `/backup` route
	http.HandleFunc("/backup", func(res http.ResponseWriter, req *http.Request) {

		// Execute influx backup
		cmd := exec.Command("influxd", "backup", "-portable", filepath.Join(backupDir, "backup"))
		err := cmd.Run()

		if err != nil {
			log.Fatal("Backup Error: ", err)
		}

		// compress influx backup
		cmd = exec.Command("tar", "cvfz", filepath.Join(backupDir, "backup.tgz"), filepath.Join(backupDir, "backup"))
		err = cmd.Run()

		if err != nil {
			log.Fatal("Backup Error: ", err)
		}

		// Send backup
		http.ServeFile(res, req, filepath.Join(backupDir, "backup.tgz"))

		// Cleanup backup dir
		cmd = exec.Command("rm", "-rf", filepath.Join(backupDir, "backup"))
		err = cmd.Run()

		if err != nil {
			log.Fatal("Backup Error: ", err)
		}

	})

	// start HTTP server with `http.DefaultServeMux` handler
	log.Fatal(http.ListenAndServe(":9000", nil))

}
