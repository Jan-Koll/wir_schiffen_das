package main

import "time"

type job struct {
	JobId                   int       `json:"id"`
	ConfigurationId         int       `json:"configuration_id"`
	Status                  State     `json:"status"`
	Success                 bool      `json:"success"`
	EngineManagementSystem  bool      `json:"engine_management_system"`
	MonitoringControlSystem string    `json:"monitoring_control_system"`
	CreatedAt               string    `json:"created_at"`
	LastModified            string    `json:"last_modified"`
	OrderCreatedAt          time.Time `json:"order_created_at"`
}
