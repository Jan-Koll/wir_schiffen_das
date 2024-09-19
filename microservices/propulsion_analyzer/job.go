package main

import "time"

type job struct {
	JobId             int       `json:"id"`
	ConfigurationId   int       `json:"configuration_id"`
	Status            State     `json:"status"`
	Success           bool      `json:"success"`
	StartingSystem    bool      `json:"starting_system"`
	PowerTransmission bool      `json:"power_transmission"`
	AuxiliaryPto      []string  `json:"auxiliary_pto"`
	CreatedAt         string    `json:"created_at"`
	LastModified      string    `json:"last_modified"`
	OrderCreatedAt    time.Time `json:"order_created_at"`
}
