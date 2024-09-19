package main

import "time"

type job struct {
	JobId           int       `json:"id"`
	ConfigurationId int       `json:"configuration_id"`
	Status          State     `json:"status"`
	Success         bool      `json:"success"`
	CoolingSystem   string    `json:"cooling_system"`
	ExhaustSystem   bool      `json:"exhaust_system"`
	CreatedAt       string    `json:"created_at"`
	LastModified    string    `json:"last_modified"`
	OrderCreatedAt  time.Time `json:"order_created_at"`
}
