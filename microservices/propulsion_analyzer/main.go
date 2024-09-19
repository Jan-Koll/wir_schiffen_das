package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var jobs []*job
var jobChan chan *job

func worker(ctx context.Context, jobChan <-chan *job) {
	for {
		select {
		case <-ctx.Done():
			return //exit worker when ctx is canceled
		case job, ok := <-jobChan:
			if !ok {
				return // Channel is closed, no more jobs. Worker exits.
			}
			if job.Status == Canceled {
				sendNotification(*job)
				removeJobFromQueue(job)
				return
			}
			processJob(job)
		}
	}
}
func removeJobFromQueue(job_to_remove *job) {
	for i, jobs_item := range jobs {
		if jobs_item == job_to_remove {
			jobs = append(jobs[:i], jobs[i+1:]...)
			return
		}
	}
}

type notification struct {
	Message job    `json:"message"`
	Sender  string `json:"sender"`
}

func sendNotification(updated_job job) {
	frontend_host := "frontend"
	_, exists := os.LookupEnv("FRONTEND_PORT")
	if !exists {
		frontend_host = "127.0.0.1"
	}
	frontend_url := fmt.Sprintf("http://%s:3000/analyze/sse/all/write", frontend_host) //todo: port for api somehow is always 3000?
	newNotification := notification{Message: updated_job, Sender: "propulsion_analyzer"}
	body, err := json.Marshal(newNotification)
	if err != nil {
		log.Println("Could not create JSON request")
	}

	req, err := http.NewRequest("POST", frontend_url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Could not create request", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		log.Println("Sending notification failed", err)
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("Sending notification was not successfull", response.StatusCode)
		return
	}

}

var client *http.Client
var stopChan chan bool = make(chan bool)

func processJob(job *job) {
	job.Status = Running
	sendNotification(*job)
	_, err := setStatus(*job)
	if err != nil {
		log.Printf("Could not update the status of job %+v in the database", *job)
	}

	log.Printf("Starting job: %+v", *job)
	job.Status = Running
	defer removeJobFromQueue(job)
	result, canceled, err := runAnalysis()
	if err != nil {
		job.Status = Failed
		log.Printf("Failed job: %+v", *job)
		sendNotification(*job)
		_, err = setResult(*job)
		if err != nil {
			log.Printf("Could not update result of job %+v in the database", *job)
		}
		failedJobs.Inc()
		return
	}
	if canceled {
		job.Status = Canceled
		log.Printf("Canceled job: %+v", *job)
		sendNotification(*job)
		_, err = setResult(*job)
		if err != nil {
			log.Printf("Could not update result of job %+v in the database", *job)
		}
		return
	}
	job.Status = Ready
	job.Success = result
	log.Printf("Completed job: %+v", *job)
	sendNotification(*job)
	_, err = setResult(*job)
	if err != nil {
		log.Printf("Could not update result of job %+v in the database", *job)
	}
	if !job.Success {
		unsuccessfulJobs.Inc()
	}
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
		panic("")
	}

	defer db.Close()
	log.Println("Successfully connected to database")

	const maxQueueSize = 5 //todo: falsch
	const amountWorkers = 2
	jobChan = make(chan *job, maxQueueSize)

	//Start worker goroutines
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < amountWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, jobChan)
		}()
	}
	// Setup HTTP client to make requests
	client = &http.Client{}
	// Setup HTTP server with versioned routes
	router := http.NewServeMux()
	router.HandleFunc("GET /health", healthCheckHandler)
	router.HandleFunc("GET /status", statusHandler)
	router.HandleFunc("POST /job", postAnalysisHandler)
	router.HandleFunc("GET /job/{job_id}", getAnalysisHandler)
	router.HandleFunc("GET /configuration/{configuration_id}/jobs", getAnalysesByConfigurationHandler)
	router.HandleFunc("PUT /job/{id}", cancelAnalysisHandler)
	router.Handle("/metrics", promhttp.Handler())

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	port, exists := os.LookupEnv("PROPULSION_ANALYZER_PORT")
	if !exists {
		port = "8085"
		log.Println("No port provided. Using default port " + port)
	}
	server := http.Server{Addr: ":" + port, Handler: Logging(Cors(v1))}
	log.Println("Starting server on port " + port)
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	// Signal handling for shutdown handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-signalChan
	log.Printf("Received signal: %v\n", sig)
	wg.Add(1)
	cancelRunningJobs()
	go func() {
		for _, job := range jobs {
			job.Status = Canceled
			setStatus(*job)
			sendNotification(*job)
			ctx.Done()
		}
		defer wg.Done()
		defer wg.Done()
		log.Println("Shutting down...")
	}()
	// Wait for workers to finish and cancel context
	wg.Wait()
	cancel()
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(map[string]interface{}{"healthy": checkHealth(), "queue_length": len(jobs)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to create json response")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
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

func postAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input job
	err := decoder.Decode(&input)
	if err != nil {
		log.Println("Could not decode body.", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode body. Please only provide a valid configuration_id and optionally a order_created_at timestamp"))
		return
	}
	// Required fields
	if input.ConfigurationId == 0 {
		log.Println("Bad request, please provide a valid configuration_id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, please provide a valid configuration_id"))
		return
	}
	// Plausibility checks
	if input.OrderCreatedAt.IsZero() {
		input.OrderCreatedAt = time.Now()
	}
	if input.OrderCreatedAt.After(time.Now()) {
		log.Println("Bad request, order_created_at is in the future. Please provide a valid timestamp.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request, order_created_at is in the future. Please provide a valid timestamp."))
		return
	}

	// Fetch fields from configuration_manager
	//todo: url builder
	requestUrl := fmt.Sprintf("http://configuration_manager:8081/v1/configuration/%d", input.ConfigurationId)
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println("Could not get data from configuration_manager", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not get data from configuration_manager\n"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Println("This configuration does not exist, thus the analysis could not be started", resp.StatusCode)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("This configuration does not exist, thus the analysis could not be started"))
			return
		}
		log.Println("Unexpected status code from configuration_manager", resp.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected status code from configuration_manager\n"))
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read the response body from configuration_manager", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read the response body from configuration_manager\n"))
		return
	}
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("Error unmarshaling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read the response body from configuration_manager\n"))
		return
	}
	var ok bool
	input.StartingSystem, ok = responseData["starting_system"].(bool)
	if !ok {
		fmt.Println("Error: starting_system is not a bool")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read the response for starting_system from configuration_manager\n"))
		return
	}
	input.PowerTransmission, ok = responseData["power_transmission"].(bool)
	if !ok {
		fmt.Println("Error: power_transmission is not a bool")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read the response for power_transmission from configuration_manager\n"))
		return
	}

	// Type conversion for []interface -> []string
	auxiliaryPtoInterface, ok := responseData["auxiliary_pto"].([]interface{})
	if !ok {
		log.Println("Error: auxiliary_pto is not a slice", responseData["auxiliary_pto"])
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not read the response for auxiliary_pto from configuration_manager\n"))
		return
	}
	if len(input.AuxiliaryPto) == 0 {
		// initialize to not violate non-null constraint in postgres
		input.AuxiliaryPto = []string{}
	}
	for _, value := range auxiliaryPtoInterface {
		auxiliaryPtoString, ok := value.(string)
		if !ok {
			log.Println("Warning: Element in auxiliary_pto is not a string")
			continue
		}
		input.AuxiliaryPto = append(input.AuxiliaryPto, auxiliaryPtoString)
	}
	input.Status = Queued
	newJob, err := createJob(input)
	if err != nil {
		log.Println("Failed to save the job to the database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save the job to the database\n"))
		return
	}

	sendNotification(newJob)
	jobs = append(jobs, &newJob)
	jobChan <- &newJob
	log.Printf("Queued new job: %+v", newJob)
	response, err := json.Marshal(newJob)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln("Failed to create json response")
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getAnalysesByConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("configuration_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := getAnalysesByConfigurationId(id)
	if err == sql.ErrNoRows || len(result) == 0 {
		w.Write([]byte("[]"))
		return
	}
	if err != nil {
		log.Println("Failed to query database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to query database."))
		return
	}

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

func getAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("job_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provided id is not valid"))
		return
	}
	job, err := getAnalysisById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No analysis with that id not found\n", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No analysis with that id found"))
			return
		}
		log.Println("Failed to query database\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to query database."))
		return
	}
	response, err := json.Marshal(job)
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

func cancelAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idx := slices.IndexFunc(jobs, func(job *job) bool { return job.JobId == id })
	if idx == -1 {
		w.WriteHeader(404)
		w.Write([]byte("No running analysis with this id found"))
		return
	}

	if jobs[idx].Status == Failed || jobs[idx].Status == Ready {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Job already finished. Skipping not possible."))
		return
	}

	result, err := setResult(*jobs[idx])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not update analysis in database. "))
		return
	}

	jobs[idx].Status = Canceled // Perform actual canceling
	//todo: make runAnalysis stopable with context cancel

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln("Failed to create json response")
	}
	log.Printf("Job #%d canceled", idx)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func runAnalysis() (sucess bool, canceled bool, err error) {
	error_probability := 0.15
	seconds_minimum, seconds_maximum := 8, 15
	timer := time.NewTimer(time.Duration(rand.IntN(seconds_maximum-seconds_minimum)+seconds_minimum) * time.Second)
	defer timer.Stop()
	select {
	case <-timer.C:
		if rand.IntN(99)+1 <= int(error_probability*100) {
			err = errors.New("analysis crashed")
		}
		sucess = rand.IntN(99)+1 >= 20
		return
	case <-stopChan:
		log.Println("Job canceled")
		canceled = true
		return
	}
}

func cancelRunningJobs() {
	stopChan <- true
}
