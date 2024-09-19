package main

import (
	"testing"
	"time"
)

func TestRunAnalysis(t *testing.T) {
	// Test case 1: Analysis finishes
	var success, canceled bool
	var err error
	maxRetries := 2
	for i := 0; i < maxRetries; i++ {
		success, canceled, err = runAnalysis()
		if err == nil {
			break
		}
		if i == maxRetries-1 {
			t.Errorf("Expected analysis to succeed, but it failed after %d retries: %v", maxRetries, err)
			return
		}
	}
	if canceled {
		t.Errorf("Expected analysis to succeed, but it was canceled")
	}

	// Test case 2: Analysis is canceled
	go func() {
		time.Sleep(1 * time.Second)
		cancelRunningJobs()
	}()
	for i := 0; i < maxRetries; i++ {
		success, canceled, err = runAnalysis()
		if err == nil {
			break
		}
		if i == maxRetries-1 {
			t.Errorf("Expected analysis to succeed, but it failed after %d retries: %v", maxRetries, err)
			return
		}
	}
	if success {
		t.Errorf("Expected analysis to be canceled, but it succeeded")
	}
	if !canceled {
		t.Errorf("Expected analysis to be canceled, but it was not")
	}
}

func TestRemoveJobFromQueue(t *testing.T) {
	// Test case 1: Remove a job from a non-empty queue & that is not in the queue
	jobToRemove := job{JobId: 1}
	removeJobFromQueue(&jobToRemove)
	if len(jobs) != 0 {
		t.Errorf("Expected queue length to be 0, got %d", len(jobs))
	}

	// Test case 2: Remove the only job from the queue
	jobToRemove = job{JobId: 1}
	jobs = append(jobs, &jobToRemove)
	if len(jobs) != 1 {
		t.Errorf("Expected queue length to be 1, got %d", len(jobs))
	}
	removeJobFromQueue(&jobToRemove)
	if len(jobs) != 0 {
		t.Errorf("Expected queue length to be 0, got %d", len(jobs))
	}

	// Test case 3: Remove a job that is not in the queue from a non-empty queue
	jobToRemove = job{JobId: 1}
	job2 := job{JobId: 2}
	job3 := job{JobId: 3}
	jobs = append(jobs, &jobToRemove)
	jobs = append(jobs, &job2)
	jobs = append(jobs, &job3)
	if len(jobs) != 3 {
		t.Errorf("Expected queue length to be 3, got %d", len(jobs))
	}
	removeJobFromQueue(&jobToRemove)
	if len(jobs) != 2 {
		t.Errorf("Expected queue length to be 2, got %d", len(jobs))
	}
	if jobs[0].JobId != job2.JobId || jobs[1].JobId != job3.JobId {
		t.Errorf("Deleted the wrong item. Expected first job to be job2 and second job to be job3, got %v", jobs)
	}
}
