package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type configuration struct {
	ConfigurationId         int      `json:"id"`
	Description             string   `json:"description"`
	Engine                  string   `json:"engine"`
	GearboxType             string   `json:"gearbox_type"`
	EngineManagementSystem  bool     `json:"engine_management_system"`
	MonitoringControlSystem string   `json:"monitoring_control_system"`
	StartingSystem          bool     `json:"starting_system"`
	PowerTransmission       bool     `json:"power_transmission"`
	AuxiliaryPto            []string `json:"auxiliary_pto"`
	OilSystem               []string `json:"oil_system"`
	FuelSystem              []string `json:"fuel_system"`
	CoolingSystem           string   `json:"cooling_system"`
	ExhaustSystem           bool     `json:"exhaust_system"`
	MountingSystem          bool     `json:"mounting_system"`
	GearboxOption           []string `json:"gearbox_option"`
	CreatedAt               string   `json:"created_at"`
	LastModified            string   `json:"last_modified"`
}

var db *sql.DB

func main() {
	// Connect to database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	var connectionError error
	db, connectionError = sql.Open("postgres", connStr)
	if connectionError != nil {
		log.Fatal("Failed to connect to database.\n", connectionError)
	}

	defer db.Close()
	log.Println("Successfully connected to database")

	// Setup HTTP server with versioned routes
	router := http.NewServeMux()
	router.HandleFunc("GET /health", healthCheckHandler)

	router.HandleFunc("POST /configuration", postConfigurationHandler)
	router.HandleFunc("GET /configuration", getAllConfigurationHandler)
	router.HandleFunc("GET /configuration/{id}", getConfigurationHandler)
	router.HandleFunc("POST /configuration/{id}", postConfigurationHandler)
	router.HandleFunc("DELETE /configuration/{id}", deleteConfigurationHandler)
	router.Handle("/metrics", promhttp.Handler())

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	port, exists := os.LookupEnv("CONFIGURATION_MANAGER_PORT")
	if !exists {
		port = "8081"
		log.Println("No port provided. Using default port " + port)
	}
	server := http.Server{Addr: ":" + port, Handler: Logging(Cors(v1))}
	log.Println("Starting server on port " + port)
	log.Fatal(server.ListenAndServe())
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	healthy := checkHealth()
	w.Header().Set("Content-Type", "application/json")
	if healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	response, err := json.Marshal(map[string]bool{"healthy": healthy})
	if err != nil {
		log.Fatalln("Failed to create json response")
	}
	log.Println(string(response))
	w.Write(response)
}

func checkHealth() (healthy bool) {
	err := db.Ping()
	if db != nil && err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func postConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var configuration configuration
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("Could not decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode body"))
		return
	}
	// Required fields
	if configuration.Description == "" || configuration.Engine == "" || configuration.GearboxType == "" {
		log.Println("Bad request, please provide at least a description, an engine and a gearbox_type")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, please provide at least a description, an engine and a gearbox_type"))
		return
	}
	// Get the id either from path-param or from body --> handler can be used for multiple method-types
	// Prioritizes path-param
	if r.PathValue("id") == "" {
		configuration.ConfigurationId, _ = strconv.Atoi(r.Form.Get("id"))
	} else {
		configuration.ConfigurationId, _ = strconv.Atoi(r.PathValue("id"))
	}
	result, err := createOrUpdate(configuration)
	if err != nil {
		if err == sql.ErrNoRows && configuration.ConfigurationId != 0 {
			log.Println("No configuration with that id not found\n", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No configuration with that id found"))
			return
		}
		log.Println("Failed to update database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update database"))
		return
	}
	//	response, err := json.Marshal(map[string]int{"id": created_configuration_id})
	response, err := json.Marshal(result)
	if err != nil {
		log.Println("Failed to create JSON response\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create JSON response."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getAllConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	configurations, err := getAllConfigurations()
	if err != nil {
		log.Println("Failed to query the database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to query the database."))
		return
	}
	response, err := json.Marshal(configurations)
	if err != nil {
		log.Println("Failed to create JSON response\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create JSON response."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	configuration, err := getConfigurationById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No configuration with that id not found\n", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No configuration with that id found"))
			return
		}
		log.Println("Failed to query database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to query database."))
		return
	}
	response, err := json.Marshal(configuration)
	if err != nil {
		log.Println("Failed to create JSON response\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create JSON response."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func deleteConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = deleteById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No configuration with that id not found\n", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No configuration with that id found"))
			return
		}
		log.Println("Failed to delete the entry in the database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete the entry in the database."))
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}
