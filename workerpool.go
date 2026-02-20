package main

import "sync"

type Job struct {
	IP   string
	Port int
}

func WorkerPool(jobs <-chan Job, results chan<- ScanResult, wg *sync.WaitGroup, cfg Config) {
	defer wg.Done()

	for job := range jobs {
		state := ScanPort(job.IP, job.Port, cfg.Timeout)
		if state == "OPEN" {
			service := DetectService(job.Port)
			results <- ScanResult{
				IP:      job.IP,
				Port:    job.Port,
				State:   state,
				Service: service.Name,
				Risk:    service.Risk,
			}
		}
	}
}
