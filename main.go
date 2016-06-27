package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kelseyhightower/mysql-healthz/healthz"
)

var (
	dataSourceName string
	databaseName   string
	healthAddr     string
	tables         string
)

func main() {
	flag.StringVar(&healthAddr, "health-addr", "0.0.0.0:10000", "Healthz HTTP listen address.")
	flag.StringVar(&dataSourceName, "data-source-name", "", "The mysql connect string.")
	flag.StringVar(&databaseName, "database-name", "healthz", "Name of the database to monitor.")
	flag.StringVar(&tables, "tables", "", "Comma seperated list of tables that must exist.")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	hc := &healthz.Config{
		Hostname: hostname,
		Database: healthz.DatabaseConfig{
			DriverName:     "mysql",
			DataSourceName: dataSourceName,
			DatabaseName:   databaseName,
			Tables:         strings.Split(tables, ","),
		},
	}

	healthzHandler, err := healthz.Handler(hc)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/healthz", healthzHandler)
	http.ListenAndServe(healthAddr, nil)
}
