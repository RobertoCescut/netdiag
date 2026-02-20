package main

import "time"

type Config struct {
	Target      string
	StartPort   int
	EndPort     int
	Threads     int
	Stealth     bool
	ExportCSV   bool
	ExportXLSX  bool
	Timeout     time.Duration
}

