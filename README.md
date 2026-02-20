
# NetDiag

**NetDiag – Network Diagnostic & Exposure Analyzer**
Professional Edition

NetDiag is a high-performance network scanning and exposure analysis tool written in Go.
It is designed for controlled environments, internal auditing, and security assessment workflows.

---

## Overview

NetDiag performs concurrent port scanning across single IPs or CIDR ranges.
It includes configurable threading, stealth mode, and structured export capabilities.

The tool is optimized for speed, modularity, and clean CLI usability.

---

## Features

* CIDR target scanning (e.g., 192.168.1.0/24)
* Custom port range selection
* Multi-threaded worker pool engine
* Optional stealth mode
* CSV export support
* Modular architecture (scanner, services, fingerprinting, exporter)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/RobertoCescut/netdiag.git
cd netdiag
```

Build:

```bash
go build -o netdiag
```

---

## Usage

Basic example:

```bash
./netdiag -target 192.168.1.0/24
```

Advanced example:

```bash
./netdiag -target 192.168.1.0/24 -start 1 -end 1024 -threads 100 -export
```

---

## Flags

| Flag     | Description              | Default  |
| -------- | ------------------------ | -------- |
| -target  | Target IP or CIDR range  | required |
| -start   | Starting port            | 1        |
| -end     | Ending port              | 1024     |
| -threads | Number of worker threads | 50       |
| -stealth | Enable stealth mode      | false    |
| -export  | Export results to CSV    | false    |

---

## Architecture

NetDiag is structured into modular components:

* `scanner.go` → scanning engine
* `workerpool.go` → concurrency management
* `services.go` → service mapping
* `fingerprint.go` → response fingerprinting
* `exporter.go` → output handling
* `config.go` → runtime configuration

The tool leverages Go's concurrency model for high-performance parallel scanning.

---

## Legal Notice

This tool is intended for:

* Authorized internal network audits
* Lab environments
* Security research within legal boundaries

Unauthorized scanning of networks without explicit permission may violate laws and regulations.

The author assumes no liability for misuse.

---

## Version

Current release: **v0.1.0**

---

## Author

Roberto Cescut

---

