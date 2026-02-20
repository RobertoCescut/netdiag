package main

import "strings"


func DetectDevice(openPorts []int) string {

	has := func(p int) bool {
		for _, port := range openPorts {
			if port == p {
				return true
			}
		}
		return false
	}

	switch {

	case has(445) && has(3389):
		return "Windows Machine"

	case has(22) && has(80):
		return "Linux Server"

	case has(53) && has(80):
		return "Router / Gateway"

	case has(21) && has(80):
		return "Embedded Appliance"

	case has(5000) || has(5001):
		return "NAS Device"

	default:
		return "Generic Device"
	}
}

func AnalyzeBanner(banner string) string {
	
	b := strings.ToLower(banner)
	
	switch {
		
	case strings.Contains(b, "nginx"):
		return "nginx"
		
	case strings.Contains(b, "apache"):
		return "Apache"
		
	case strings.Contains(b, "openssh"):
		return "OpenSSH"
		
	case strings.Contains(b, "microsoft-iis"):
		return "Microsoft IIS"
		
	default:
		return ""
	}
}
