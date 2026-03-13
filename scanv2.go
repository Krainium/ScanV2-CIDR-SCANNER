package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	G  = "\033[32m"
	O  = "\033[33m"
	GR = "\033[37m"
	R  = "\033[31m"
	RS = "\033[0m"
)

var (
	mu      sync.Mutex
	hitFile *os.File
	httpRe  = regexp.MustCompile(`^HTTP/\d(\.\d)?`)
	scanned int64
	total   int64
	debug   bool
	startTime time.Time
)

func scanner(host, port string) {
	timeout := 3 * time.Second
	sock, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		return
	}
	defer sock.Close()
	sock.SetDeadline(time.Now().Add(timeout))

	payload := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", host)
	sock.Write([]byte(payload))

	buf := make([]byte, 1024)
	n, _ := sock.Read(buf)
	if n == 0 {
		return
	}

	response := string(buf[:n])
	hit := false
	for _, line := range strings.Split(response, "\r\n") {
		if httpRe.MatchString(line) {
			fmt.Printf("\n%s[HIT] %s → %s%s\n", G, host, line, RS)
		}
		if strings.Contains(strings.ToLower(line), "content-type:") {
			hit = true
		}
	}

	if hit {
		mu.Lock()
		if hitFile != nil {
			hitFile.WriteString(host + "\n")
		}
		mu.Unlock()
	}
}

func expandCIDR(cidr string) ([]string, error) {
	cidr = strings.TrimSpace(cidr)
	if !strings.Contains(cidr, "/") {
		if net.ParseIP(cidr) != nil {
			return []string{cidr}, nil
		}
		return nil, fmt.Errorf("invalid IP/CIDR: %s", cidr)
	}
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	startIP := network.IP.To4()
	if startIP == nil {
		return nil, fmt.Errorf("only IPv4 supported: %s", cidr)
	}
	start := binary.BigEndian.Uint32(startIP)
	mask := binary.BigEndian.Uint32([]byte(network.Mask))
	end := (start & mask) | (^mask)

	ips := make([]string, 0, end-start+1)
	for i := start; i <= end; i++ {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, i)
		ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3]))
	}
	return ips, nil
}

func main() {
	fmt.Printf("\n%sModified Scanner (High Performance)%s\n", O, RS)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the IP ranges file: ")
	filepath, _ := reader.ReadString('\n')
	filepath = strings.TrimSpace(filepath)

	fmt.Print("Enter port number: ")
	portStr, _ := reader.ReadString('\n')
	portStr = strings.TrimSpace(portStr)

	fmt.Print("Enter concurrent threads (default 200): ")
	threadStr, _ := reader.ReadString('\n')
	threads, _ := strconv.Atoi(strings.TrimSpace(threadStr))
	if threads <= 0 { threads = 200 }

	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("%sCould not open file: %v%s\n", R, err, RS)
		os.Exit(1)
	}
	defer f.Close()

	var ipList []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" { continue }
		ips, _ := expandCIDR(line)
		ipList = append(ipList, ips...)
	}

	total = int64(len(ipList))
	if total == 0 {
		fmt.Println("No IPs found.")
		return
	}

	hitFile, _ = os.OpenFile("live.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Printf("\n%s[INFO] Starting scan on %d IPs with %d threads...%s\n", G, total, threads, RS)

	jobs := make(chan string, threads)
	var wg sync.WaitGroup
	startTime = time.Now()

	// Start Worker Pool
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range jobs {
				scanner(ip, portStr)
				curr := atomic.AddInt64(&scanned, 1)

				// Calculate ETA every 10 IPs to save CPU
				if curr % 10 == 0 || curr == total {
					elapsed := time.Since(startTime).Seconds()
					ipsPerSec := float64(curr) / elapsed
					remainingIps := float64(total - curr)
					etaSeconds := int(remainingIps / ipsPerSec)
					
					eta := time.Duration(etaSeconds) * time.Second
					
					// Clear line and print status
					fmt.Printf("\r\033[K%s[PROBING] %d/%d (%d%%) | Rate: %.1f ip/s | ETA: %s%s", 
						GR, curr, total, (curr*100)/total, ipsPerSec, eta.String(), RS)
				}
			}
		}()
	}

	// Fill queue
	for _, ip := range ipList {
		jobs <- ip
	}
	close(jobs)
	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("\n\n%s[DONE] Finished in %s. Hits saved to live.txt%s\n", G, duration.String(), RS)
}