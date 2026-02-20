package main

import (
	"flag"
	"fmt"
	"time"
)
const version = "0.1.0"

func printBanner() {
	fmt.Println(`
███╗   ██╗███████╗████████╗██████╗ ██╗ █████╗  ██████╗ 
████╗  ██║██╔════╝╚══██╔══╝██╔══██╗██║██╔══██╗██╔════╝ 
██╔██╗ ██║█████╗     ██║   ██║  ██║██║███████║██║  ███╗
██║╚██╗██║██╔══╝     ██║   ██║  ██║██║██╔══██║██║   ██║
██║ ╚████║███████╗   ██║   ██████╔╝██║██║  ██║╚██████╔╝
╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═════╝ ╚═╝╚═╝  ╚═╝ ╚═════╝ 

		NetDiag - Network Diagnostic & Exposure Analyzer
		Professional Edition
`)
	
	fmt.Println("        Version:", version)
	fmt.Println()
}


const green = "\033[32m"
const reset = "\033[0m"

func clearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func main() {

	printBanner()

	target := flag.String("target", "", "Target CIDR (es 192.168.1.0/24)")
	startPort := flag.Int("start", 1, "Porta iniziale")
	endPort := flag.Int("end", 1024, "Porta finale")
	threads := flag.Int("threads", 50, "Numero worker")
	stealth := flag.Bool("stealth", false, "Modalità stealth")
	export := flag.Bool("export", false, "Export CSV")

	flag.Parse()

	var cfg Config

	// 🔥 MODALITÀ CLI
	if *target != "" {

		fmt.Println("Modalità CLI professionale")

		cfg = Config{
			Target:     *target,
			StartPort:  *startPort,
			EndPort:    *endPort,
			Threads:    *threads,
			Stealth:    *stealth,
			ExportXLSX: *export,
			Timeout:    300 * time.Millisecond,
		}

	} else {

		// 🔥 MODALITÀ INTERATTIVA
		fmt.Println("Modalità interattiva")

		var t string
		fmt.Print("Target (IP o CIDR): ")
		fmt.Scanln(&t)

		var sp, ep int
		fmt.Print("Porta iniziale: ")
		fmt.Scanln(&sp)

		fmt.Print("Porta finale: ")
		fmt.Scanln(&ep)

		var th int
		fmt.Print("Numero thread (default 50): ")
		fmt.Scanln(&th)
		if th == 0 {
			th = 50
		}

		var s string
		fmt.Print("Modalità stealth? (y/n): ")
		fmt.Scanln(&s)

		var eCSV string
		var eXLSX string
		
		fmt.Print("Export CSV? (y/n): ")
		fmt.Scanln(&eCSV)
		
		fmt.Print("Export Excel? (y/n): ")
		fmt.Scanln(&eXLSX)
		
		cfg = Config{
			Target:     t,
			StartPort:  sp,
			EndPort:    ep,
			Threads:    th,
			Stealth:    s == "y",
			ExportCSV:  eCSV == "y",
			ExportXLSX: eXLSX == "y",
			Timeout:    300 * time.Millisecond,
		}
		
	}

	StartScan(cfg)
}
