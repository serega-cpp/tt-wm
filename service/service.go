package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"finleap/api"
	"finleap/database"

	"github.com/gorilla/mux"
)

type Service struct {
	db *database.DB
}

func Create(host string, port int, user, pass, dbname string) (*Service, error) {
	db, err := database.Create(host, port, user, pass, dbname)
	if err != nil {
		return nil, err
	}
	return &Service{db}, nil
}

func (s *Service) Destroy() error {
	return s.db.Destroy()
}

func (s *Service) Prepare() error {
	if err := s.db.CreateCitiesTable(); err != nil {
		return fmt.Errorf("CREATE TABLE Cities: %s", err)
	}
	if err := s.db.CreateTemperaturesTable(); err != nil {
		return fmt.Errorf("CREATE TABLE Temperatures: %s", err)
	}
	if err := s.db.CreateWebhooksTable(); err != nil {
		return fmt.Errorf("CREATE TABLE Webhooks: %s", err)
	}
	return nil
}

func (s *Service) Reset() error {
	if err := s.db.ClearCitiesTable(); err != nil {
		return fmt.Errorf("DELETE FROM Cities: %s", err)
	}
	if err := s.db.ClearTemperaturesTable(); err != nil {
		return fmt.Errorf("DELETE FROM Temperatures: %s", err)
	}
	if err := s.db.ClearWebhooksTable(); err != nil {
		return fmt.Errorf("DELETE FROM Webhooks: %s", err)
	}
	return nil
}

///////////////////////////////////////////////////////////
// REST handlers

func (s *Service) CreateCity(w http.ResponseWriter, r *http.Request) {
	city, err := readCity(r)
	if err != nil {
		log.Printf("ERROR: Request parsing failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.db.InsertCity(city); err != nil {
		log.Printf("ERROR: InsertDB failed for %#v: %s", *city, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(city)
	if err != nil {
		panic(err.Error()) // need to know ASAP
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Service) ModifyCity(w http.ResponseWriter, r *http.Request) {
	city, err := readCity(r)
	if err != nil {
		log.Printf("ERROR: Request parsing failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	city.ID, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("ERROR: ID parsing failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.db.UpdateCity(city); err != nil {
		log.Printf("ERROR: UpdateDB failed for %#v: %s", *city, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(city)
	if err != nil {
		panic(err.Error()) // need to know ASAP
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Service) DeleteCity(w http.ResponseWriter, r *http.Request) {
	city_id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("ERROR: ID parsing failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	city := &api.City{
		ID: city_id,
	}
	if err := s.db.DeleteCity(city); err != nil {
		log.Printf("ERROR: DeleteDB failed for %#v: %s", *city, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(city)
	if err != nil {
		panic(err.Error()) // need to know ASAP
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Service) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	temperature, err := readTemperature(r)
	if err != nil {
		log.Printf("ERROR: Request parsing failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.db.InsertTemperature(temperature); err != nil {
		log.Printf("ERROR: InsertDB failed for %#v: %s", *temperature, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.notifyCallbacks(temperature)

	body, err := json.Marshal(temperature)
	if err != nil {
		panic(err.Error()) // need to know ASAP
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Service) GetForecasts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("getForecasts for: %s", vars["city_id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	log.Printf("createWebhook")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("deleteWebhook: %s", vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) notifyCallbacks(temperature *api.Temperature) {
	log.Printf("notifyCallbacks: %d", temperature.CityID)
}

///////////////////////////////////////////////////////////
// Helpers

func readCity(r *http.Request) (*api.City, error) {
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Wrong latitude: %s", err)
	}
	lon, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Wrong longitude: %s", err)
	}

	city := api.City{
		Name:      r.FormValue("name"),
		Latitude:  lat,
		Longitude: lon,
	}
	if err := city.Validate(); err != nil {
		return nil, fmt.Errorf("ERROR: Wrong City value: %s", err)
	}
	return &city, nil
}

func readTemperature(r *http.Request) (*api.Temperature, error) {
	city_id, err := strconv.ParseInt(r.FormValue("city_id"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Wrong city_id: %s", err)
	}
	maxC, err := strconv.ParseFloat(r.FormValue("max"), 32)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Wrong max temperature: %s", err)
	}
	minC, err := strconv.ParseFloat(r.FormValue("min"), 32)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Wrong min temperature: %s", err)
	}

	temperature := api.Temperature{
		CityID:    int(city_id),
		MaxC:      float32(maxC),
		MinC:      float32(minC),
		Timestamp: time.Now().Unix(),
	}
	if err := temperature.Validate(); err != nil {
		return nil, fmt.Errorf("ERROR: Wrong Temperature value: %s", err)
	}
	return &temperature, nil
}
