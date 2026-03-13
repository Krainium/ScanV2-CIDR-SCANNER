# ⚡ ScanV2 – High Performance HTTP Service Scanner

![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/status-active-success)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20windows%20%7C%20macOS-blue)

**ScanV2** is a fast, multithreaded **Go-based network scanner** designed to probe large ranges of IP addresses and detect active HTTP services on a specific port.

The tool expands **CIDR ranges into individual IPs**, scans them concurrently using a worker pool, and detects live hosts by analyzing HTTP responses. All discovered live hosts are automatically saved to a file.

This scanner is optimized for **speed, concurrency, and large-scale scanning tasks**.

---

# ✨ Features

* ⚡ **High-performance Go scanner**
* 🌐 Supports **CIDR ranges and single IPs**
* 🧵 **Concurrent worker pool scanning**
* 📊 Real-time **scan progress, rate, and ETA**
* 🔎 Detects **HTTP services via response headers**
* 💾 Saves discovered hosts to `live.txt`
* 📈 Displays **scan rate (IP/s) and progress percentage**
* ⏱ Shows **estimated time remaining (ETA)**

---

# 🎥 Demo

Example scanner run:

```
Modified Scanner (High Performance)

Enter the IP ranges file: ranges.txt
Enter port number: 80
Enter concurrent threads (default 200): 500

[INFO] Starting scan on 65536 IPs with 500 threads...

[PROBING] 32000/65536 (48%) | Rate: 2400 ip/s | ETA: 12s

[HIT] 192.168.1.20 → HTTP/1.1 200 OK
[HIT] 192.168.1.45 → HTTP/1.1 403 Forbidden

[DONE] Finished in 28s. Hits saved to live.txt
```

You can optionally add a **terminal demo GIF** here later.

---

# 🧰 Requirements

* **Go 1.18+**

ScanV2 uses only **Go standard library packages**, so **no external dependencies are required**.

---

# 📦 Installation

Clone the repository:

```
git clone https://github.com/krainium/scanv2-http-scanner.git
cd scanv2-http-scanner
```

---

# 🚀 Usage

Run the scanner:

```
go run scanv2.go
```

Or build a binary:

```
go build scanv2.go
```

Run the compiled binary:

```
./scanv2
```

---

# 📝 Input Parameters

The program will ask for:

### 1️⃣ IP ranges file

A text file containing **IP addresses or CIDR ranges**.

Example `ranges.txt`:

```
192.168.1.1
192.168.1.0/24
10.0.0.0/24
8.8.8.8
```

---

### 2️⃣ Target port

The port to scan for HTTP services.

Example:

```
80
8080
8000
443
```

---

### 3️⃣ Thread count

Number of concurrent scanning threads.

Example:

```
200
500
1000
```

If left empty, the default is **200 threads**.

---

# 📂 Output

All discovered hosts that respond with HTTP headers are saved to:

```
live.txt
```

Example output file:

```
192.168.1.10
192.168.1.22
192.168.1.45
```

---

# 📊 Performance

Approximate benchmark results (varies by network and hardware):

| Threads | Speed           |
| ------- | --------------- |
| 200     | ~800–1200 IP/s  |
| 500     | ~1800–2500 IP/s |
| 1000    | ~3000+ IP/s     |

Go's concurrency model allows ScanV2 to process **tens of thousands of IPs very quickly**.

---

# 📂 Project Structure

```
scanv2-http-scanner/
│
├── scanv2.go
├── ranges.txt
├── live.txt
└── README.md
```

---

# ⚠️ Disclaimer

This tool is intended for:

* Network diagnostics
* Security testing
* Educational purposes

Only scan systems **you own or have permission to test**. Unauthorized scanning may violate laws or network policies.

---

# 🤝 Contributions

Contributions are welcome.

If you'd like to improve ScanV2:

1. Fork the repository
2. Create a new branch
3. Commit your changes
4. Submit a pull request

Possible contribution areas:

* Additional protocol detection (HTTPS, FTP, SSH)
* IPv6 support
* Faster CIDR expansion
* Improved response analysis
* Performance optimizations
* Logging improvements

---

# 🚀 Future Updates

Planned improvements:

* 🔐 **HTTPS detection**
* 🌍 **IPv6 scanning support**
* 📡 **Multi-port scanning**
* 📊 **Live statistics dashboard**
* 🧠 **Service fingerprinting**
* 📁 **JSON / CSV export**
* ⚡ **Automatic port detection**
* 🔎 **Banner grabbing**
* 🧵 **Adaptive thread scaling**
* 📦 **Precompiled binaries**

---

# ⭐ Support

If you find this project useful, consider **starring the repository**.

---

# 👨‍💻 Author

Krainium
GitHub: [https://github.com/krainium](https://github.com/krainium)

---

# 📜 License

MIT License
