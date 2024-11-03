package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	pinger []string
	mu     sync.Mutex
)

func ping(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "1 received")
}

func scanNetwork(wg *sync.WaitGroup, ip string) {
	defer wg.Done()
	if ping(ip) {
		mu.Lock()
		pinger = append(pinger, ip)
		mu.Unlock()
	}
}

func main() {
	startTime := time.Now()
	netIp := "192.168.1.0"
	// fake udp-dial for detect main interface
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localAddrIp := localAddr.IP.String()
	if localAddrIp != "" {
		netIp = localAddrIp
	}

	if len(os.Args) > 1 {
		netIp = os.Args[1]
	}

	re := regexp.MustCompile(`^(?P<firstThreeOctets>\d{1,3}\.\d{1,3}\.\d{1,3})\.\d{1,3}$`)
	if match := re.FindStringSubmatch(netIp); match != nil {
		netIp = match[1]
	} else {
		fmt.Println("Incorrect IP-address. Need: 192.168.1.0")
		return
	}

	var wg sync.WaitGroup
	for i := 1; i < 255; i++ {
		ip := fmt.Sprintf("%s.%d", netIp, i)
		wg.Add(1)
		go scanNetwork(&wg, ip)
	}
	wg.Wait()

	sort.Slice(pinger, func(i, j int) bool {
		lastOctetI, _ := strconv.Atoi(strings.Split(pinger[i], ".")[3])
		lastOctetJ, _ := strconv.Atoi(strings.Split(pinger[j], ".")[3])
		return lastOctetI < lastOctetJ
	})

	if len(pinger) == 0 {
		fmt.Println("No results for " + netIp + ".0")
	} else {
		for _, v := range pinger {
			if v == localAddrIp  || (len(os.Args) > 1 && os.Args[1] == v){
				color.Green(v)
			} else {
				fmt.Println(v)
			}
		}
		fmt.Printf("Took: %s\n", time.Since(startTime).Truncate(time.Millisecond))
	}
}
