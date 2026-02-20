package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type ScanResult struct {
	IP      string
	Port    int
	State   string
	Service string
	Risk    string
	Elapsed time.Duration
	
}

func ScanPort(ip string, port int, timeout time.Duration) string {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	
	if err != nil {
		
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "FILTERED"
		}
		
		if strings.Contains(err.Error(), "connection refused") {
			return "CLOSED"
		}
		
		if strings.Contains(err.Error(), "no route to host") {
			return "UNREACHABLE"
		}
		
		return "FILTERED"
	}
	
	conn.Close()
	return "OPEN"
}

func generateHosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		ips = append(ips, ipCopy.String())
	}

	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}

	return ips, nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type scanTask struct {
	ip   string
	port int
}

func StartScan(cfg Config) {
	
	startTime := time.Now()
	var totalResponseTime time.Duration
	
	
	hosts, err := generateHosts(cfg.Target)
	if err != nil {
		println("Errore CIDR:", err.Error())
		return
	}
	
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 1 * time.Second
	}
	cfg.Timeout = timeout
	
	var results []ScanResult
	totalTasks := len(hosts) * (cfg.EndPort - cfg.StartPort + 1)
	completed := 0
	
	tasks := make(chan scanTask)
	resultsChan := make(chan ScanResult)
	
	// Contatori professionali
	var openCount, closedCount, filteredCount, unreachableCount int
	
	// Worker
	for i := 0; i < cfg.Threads; i++ {
		go func() {
			for task := range tasks {
				
				scanStart := time.Now()
				state := ScanPort(task.ip, task.port, cfg.Timeout)
				elapsed := time.Since(scanStart)
				
				
				resultsChan <- ScanResult{
					IP:    task.ip,
					Port:  task.port,
					State: state,
					Elapsed: elapsed,
				}
			}
		}()
	}
	
	// Generazione task
	go func() {
		for _, ip := range hosts {
			for port := cfg.StartPort; port <= cfg.EndPort; port++ {
				tasks <- scanTask{ip: ip, port: port}
			}
		}
		close(tasks)
	}()
	
	// Raccolta risultati
	for i := 0; i < totalTasks; i++ {
		
		result := <-resultsChan
		completed++
		
		totalResponseTime += result.Elapsed
		
		// Conteggio stati
		switch result.State {
		case "OPEN":
			openCount++
		case "CLOSED":
			closedCount++
		case "FILTERED":
			filteredCount++
		case "UNREACHABLE":
			unreachableCount++
		}
		
		percent := float64(completed) / float64(totalTasks) * 100
		renderProgressBar(percent)
		
		if result.State == "OPEN" {
			
			serviceInfo := DetectService(result.Port)
			
			var banner string
			
			switch result.Port {
			case 80, 443:
				banner = grabHTTPBanner(result.IP, result.Port)
			case 22:
				banner = grabSSHBanner(result.IP, result.Port)
			}
			// Analisi banner
			if banner != "" {
				detected := AnalyzeBanner(banner)
				if detected != "" {
					serviceInfo.Name = detected
				}
			}
			
			
			fmt.Printf("\n%s : %d OPEN [%s | %s]\n",
				result.IP,
				result.Port,
				serviceInfo.Name,
				serviceInfo.Risk,
			)
			
			if banner != "" {
				
				// stampa banner
				fmt.Println("   Banner:", banner)
				
				// analizza banner
				detected := AnalyzeBanner(banner)
				if detected != "" {
					serviceInfo.Name = detected
				}
			}
			
			results = append(results, ScanResult{
				IP:      result.IP,
				Port:    result.Port,
				State:   result.State,
				Service: serviceInfo.Name,
				Risk:    serviceInfo.Risk,
			})
		}
	}
	
	fmt.Println("\nScansione completata.")
	
	// Summary professionale
	duration := time.Since(startTime)
	speed := float64(totalTasks) / duration.Seconds()
	avgLatency := totalResponseTime / time.Duration(totalTasks)
	
	fmt.Println("\n==============================================")
	fmt.Println("                SCAN SUMMARY")
	fmt.Println("==============================================")
	fmt.Printf("Open ports: %d\n", openCount)
	fmt.Printf("Closed ports: %d\n", closedCount)
	fmt.Printf("Filtered ports: %d\n", filteredCount)
	fmt.Printf("Unreachable: %d\n", unreachableCount)
	fmt.Printf("Tempo totale: %s\n", duration)
	fmt.Printf("Velocità media: %.2f porte/sec\n", speed)
	fmt.Printf("Latenza media stimata: %s\n", avgLatency)
	fmt.Println("==============================================")
	
	// Raggruppa porte per IP
	devicePorts := make(map[string][]int)
	
	for _, r := range results {
		devicePorts[r.IP] = append(devicePorts[r.IP], r.Port)
	}
	
	// Classificazione dispositivi
	for ip, ports := range devicePorts {
		device := DetectDevice(ports)
		fmt.Println("Device rilevato su", ip, ":", device)
		fmt.Println("----------------------------")
	}
	
	
	
	// EXPORT CSV
	if cfg.ExportCSV {
		file, err := os.Create("netdiag_results.csv")
		if err != nil {
			println("Errore creazione file")
		} else {
			defer file.Close()
			
			file.WriteString("IP,Port,State,Service,Risk\n")
			
			for _, r := range results {
				
				cleanService := strings.ReplaceAll(r.Service, "\n", " ")
				cleanService = strings.ReplaceAll(cleanService, "\r", " ")
				
				cleanRisk := strings.ReplaceAll(r.Risk, "\n", " ")
				cleanRisk = strings.ReplaceAll(cleanRisk, "\r", " ")
				
				line := fmt.Sprintf("%s,%d,%s,%s,%s\n",
					r.IP,
					r.Port,
					r.State,
					cleanService,
					cleanRisk,
				)
				
				file.WriteString(line)
			}
			
			println("Export completato: netdiag_results.csv")
		}
	}
	// EXPORT EXCEL
	if cfg.ExportXLSX {
		err := ExportExcel(results)
		if err != nil {
			fmt.Println("Errore export Excel:", err)
		} else {
			fmt.Println("Export completato: scan_results.xlsx")
		}
	}
}
	
func renderProgressBar(percent float64) {
	barWidth := 30
	filled := int(percent / 100 * float64(barWidth))

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"

	fmt.Printf("\r%s %.0f%%", bar, percent)
}
func grabHTTPBanner(ip string, port int) string {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return ""
	}
	defer conn.Close()

	request := "HEAD / HTTP/1.0\r\n\r\n"
	conn.Write([]byte(request))

	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}

	raw := string(buffer[:n])
	lines := strings.Split(raw, "\n")

	var result string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "HTTP/") {
			result += line + " | "
		}

		if strings.HasPrefix(strings.ToLower(line), "server:") {
			result += line
		}
	}

	return result
}

func grabSSHBanner(ip string, port int) string {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return ""
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}

	return string(buffer[:n])
}
	