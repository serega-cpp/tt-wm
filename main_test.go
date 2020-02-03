package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"finleap/api"
	"finleap/config"
	"finleap/service"

	"github.com/gorilla/mux"
)

var rest *service.Service

func TestMain(m *testing.M) {
	cfgFileName := flag.String("cfg", "config_test.yaml", "config file")
	flag.Parse()

	cfg, err := config.ReadConfig(*cfgFileName)
	if err != nil {
		log.Fatal(err)
	}
	svc, err := service.Create(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Pass, cfg.DB.DBname)
	if err != nil {
		log.Fatal(err)
	}
	if err = svc.Prepare(); err != nil {
		log.Fatal(err)
	}
	if err = svc.Reset(); err != nil {
		log.Fatal(err)
	}

	rest = svc
	os.Exit(m.Run())
}

func TestCreateCity(t *testing.T) {
	city := api.City{
		Name:      "Moscow",
		Longitude: 37.618,
		Latitude:  55.751,
	}
	body := url.Values{
		"name":      {city.Name},
		"longitude": {strconv.FormatFloat(city.Longitude, 'f', -1, 64)},
		"latitude":  {strconv.FormatFloat(city.Latitude, 'f', -1, 64)},
	}

	req, err := http.NewRequest("POST", "/cities", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/cities", rest.CreateCity)
	router.ServeHTTP(rec, req)

	validateStatus(t, rec.Code, http.StatusOK)

	var result api.City
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	city.ID = result.ID
	if city != result || result.ID == 0 {
		t.Errorf("Wrong result: got %v expected %v (expected ID > 0)", result, city)
	}
}

func TestModifyCity(t *testing.T) {
	city := api.City{
		ID:        1,
		Name:      "Berlin",
		Longitude: 13.41,
		Latitude:  52.52,
	}
	body := url.Values{
		"name":      {city.Name},
		"longitude": {strconv.FormatFloat(city.Longitude, 'f', -1, 64)},
		"latitude":  {strconv.FormatFloat(city.Latitude, 'f', -1, 64)},
	}

	req, err := http.NewRequest("PATCH", "/cities/1", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/cities/{id}", rest.ModifyCity)
	router.ServeHTTP(rec, req)

	validateStatus(t, rec.Code, http.StatusOK)

	var result api.City
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if city != result {
		t.Errorf("Wrong result: got %v expected %v", result, city)
	}
}

func TestDeleteCity(t *testing.T) {
	city := api.City{
		ID:        1,
		Name:      "Berlin",
		Longitude: 13.41,
		Latitude:  52.52,
	}

	req, err := http.NewRequest("DELETE", "/cities/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/cities/{id}", rest.DeleteCity)
	router.ServeHTTP(rec, req)

	validateStatus(t, rec.Code, http.StatusOK)

	var result api.City
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if city != result {
		t.Errorf("Wrong result: got %v expected %v", result, city)
	}
}

func validateStatus(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("Status code: got %v expected %v", got, expected)
	}
}
