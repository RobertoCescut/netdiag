package main

type ServiceInfo struct {
	Name string
	Risk string
}

var knownServices = map[int]ServiceInfo{
	21:   {"FTP", "HIGH"},
	22:   {"SSH", "MEDIUM"},
	23:   {"TELNET", "CRITICAL"},
	53:   {"DNS", "MEDIUM"},
	80:   {"HTTP", "LOW"},
	443:  {"HTTPS", "LOW"},
	445:  {"SMB", "HIGH"},
	554:  {"RTSP", "MEDIUM"},
	3389: {"RDP", "HIGH"},
}

func DetectService(port int) ServiceInfo {
	if s, ok := knownServices[port]; ok {
		return s
	}
	return ServiceInfo{"UNKNOWN", "UNKNOWN"}
}
