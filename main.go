package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"finleap/config"
	"finleap/service"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Weather Service started")
	cfgFileName := flag.String("cfg", "config.yaml", "config file")
	flag.Parse()

	log.Printf("Configuration file: %s", *cfgFileName)
	cfg, err := config.ReadConfig(*cfgFileName)
	if err != nil {
		log.Fatalf("FATAL: Config load error (%s): %s", *cfgFileName, err)
	}

	log.Printf("Database instance: %s:xxx@%s:%d/%s", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBname)
	svc, err := service.Create(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Pass, cfg.DB.DBname)
	if err != nil {
		log.Fatalf("FATAL: Connection to PostgreSQL failed: %s", err)
	}
	defer svc.Destroy()

	if err = svc.Prepare(); err != nil {
		log.Fatalf("FATAL: Prepare: %s", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/cities", svc.CreateCity).Methods("POST")
	r.HandleFunc("/cities/{id:[0-9]+}", svc.ModifyCity).Methods("PATCH")
	r.HandleFunc("/cities/{id:[0-9]+}", svc.DeleteCity).Methods("DELETE")
	r.HandleFunc("/temperatures", svc.CreateMeasurement).Methods("POST")
	r.HandleFunc("/forecasts/{id:[0-9]+}", svc.GetForecasts).Methods("GET")
	r.HandleFunc("/webhooks", svc.CreateWebhook).Methods("POST")
	r.HandleFunc("/webhooks/{id:[0-9]+}", svc.DeleteWebhook).Methods("DELETE")

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("REST server is listening: %s", addr)
	if err = srv.ListenAndServe(); err != nil {
		log.Fatalf("FATAL: REST server on %s failed: %s", addr, err)
	}
}
