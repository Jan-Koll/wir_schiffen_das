package main

import (
	"log"

	"github.com/lib/pq"
)

func getAnalysisById(id int) (result job, err error) {
	row := db.QueryRow("SELECT * FROM propulsion_analyzer WHERE job_id = $1", id)
	err = row.Scan(&result.JobId, &result.ConfigurationId, &result.Success, &result.StartingSystem, &result.PowerTransmission, pq.Array(&result.AuxiliaryPto), &result.CreatedAt, &result.LastModified, &result.Status, &result.OrderCreatedAt)
	return
}

func createJob(input job) (result job, err error) {
	result = input
	row := db.QueryRow(
		`INSERT INTO propulsion_analyzer (configuration_id, status, starting_system, power_transmission, auxiliary_pto, order_created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING job_id, created_at, last_modified, order_created_at`,
		input.ConfigurationId,
		input.Status,
		input.StartingSystem,
		input.PowerTransmission,
		pq.Array(input.AuxiliaryPto),
		input.OrderCreatedAt,
	)
	err = row.Scan(&result.JobId, &result.CreatedAt, &result.LastModified, &result.OrderCreatedAt)
	if err != nil {
		result = job{}
	}
	return
}

func setStatus(input job) (result job, err error) {
	result = input
	row := db.QueryRow(
		`UPDATE propulsion_analyzer SET 
				status = $2,
				last_modified = now()
			WHERE job_id = $1
			RETURNING last_modified;`,
		input.JobId,
		input.Status,
	)
	err = row.Scan(&result.LastModified)
	if err != nil {
		result = job{}
	}
	return
}
func setResult(input job) (result job, err error) {
	result = input
	row := db.QueryRow(
		`UPDATE propulsion_analyzer SET
				status = $2,
				success = $3,
				last_modified = now()
			WHERE job_id = $1
			RETURNING last_modified;`,
		input.JobId,
		input.Status,
		input.Success,
	)
	err = row.Scan(&result.LastModified)
	if err != nil {
		result = job{}
	}
	return
}

func deleteByConfigurationId(id int) (err error) {
	_, err = db.Exec("DELETE FROM propulsion_analyzer WHERE configuration_id = $1", id)
	return
}

func getAnalysesByConfigurationId(configuration_id int) (result []job, err error) {
	rows, err := db.Query("SELECT * FROM propulsion_analyzer WHERE configuration_id = $1 ORDER BY job_id ASC", configuration_id)
	if err != nil {
		log.Println("Failed to query database\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var job job
		err = rows.Scan(&job.JobId, &job.ConfigurationId, &job.Success, &job.StartingSystem, &job.PowerTransmission, pq.Array(&job.AuxiliaryPto), &job.CreatedAt, &job.LastModified, &job.Status, &job.OrderCreatedAt)
		if err != nil {
			log.Println("Failed to read row\n", err)
			return
		}
		result = append(result, job)
	}
	return
}
