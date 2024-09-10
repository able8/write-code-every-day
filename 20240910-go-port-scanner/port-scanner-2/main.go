package main

import (
	"fmt"
	"net"
	"sort"
	"time"
	"unsafe"
)

func main() {
	// tcpScanByGoroutineWithChannelAndSort("127.0.0.1", 1, 65535)
	tcpScanByGoroutineWithChannelAndSort("114.114.114.114", 1, 65535)
}

// The function handles checking if ports are open or closed for a given IP address.
func handleWorker(ip string, ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", ip, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("[debug] ip %s Close \n", address)
			results <- (-p)
			continue
		}
		fmt.Printf("[debug] ip %s Open \n", address)
		conn.Close()
		results <- p
	}
}

func tcpScanByGoroutineWithChannelAndSort(ip string, portStart int, portEnd int) {
	start := time.Now()
	// Parameter verification
	isok := verifyParam(ip, portStart, portEnd)
	if !isok {
		fmt.Printf("[Exit]\n")
	}

	ports := make(chan int, 50)
	results := make(chan int)
	var openSlice []int
	var closeSlice []int

	// Task producer-distribute tasks (create a new goroutinue to distribute data)
	go func(a int, b int) {
		for i := a; i <= b; i++ {
			ports <- i
		}
	}(portStart, portEnd)

	// Task consumer-processing tasks (each port number is assigned a goroutinue for scanning)
	// Result producer-every time the result is obtained, it is written into the result chan
	for i := 0; i < cap(ports); i++ {
		go handleWorker(ip, ports, results)
	}

	// Result consumer - waiting to collect results (goroutinue in main continuously reads data from chan in a blocking manner)
	for i := portStart; i <= portEnd; i++ {
		resPort := <-results
		if resPort > 0 {
			openSlice = append(openSlice, resPort)
		} else {
			closeSlice = append(closeSlice, -resPort)
		}
	}

	// close chan
	close(ports)
	close(results)

	// sort
	sort.Ints(openSlice)
	sort.Ints(closeSlice)

	// output
	for _, p := range openSlice {
		fmt.Printf("[info] %s:%-8d Open\n", ip, p)
	}
	// for _, p := range closeSlice {
	// fmt.Printf("[info] %s:%-8d Close\n", ip, p)
	// }

	cost := time.Since(start)
	fmt.Printf("[tcpScanByGoroutineWithChannelAndSort] cost %s second \n", cost)
}

func verifyParam(ip string, portStart int, portEnd int) bool {
	netip := net.ParseIP(ip)
	if netip == nil {
		fmt.Println("[Error] ip type is must net.ip")
		return false
	}
	fmt.Printf("[Info] ip=%s | ip type is: %T | ip size is: %d \n", netip, netip, unsafe.Sizeof(netip))

	if portStart < 1 || portEnd > 65535 {
		fmt.Println("[Error] port is must in the range of 1~65535")
		return false
	}
	fmt.Printf("[Info] port start:%d end:%d \n", portStart, portEnd)

	return true
}
