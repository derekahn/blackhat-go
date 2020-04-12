package main

import (
	"fmt"
	"net"
	"sort"
)

var (
	limit  int    = 65535 // ALL ports
	target string = "scanme.nmap.org:%d"
)

func main() {
	ports := make(chan int, 250)
	results := make(chan int)

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	var openPorts []int
	go func() {
		for i := 0; i < limit; i++ {
			ports <- i
		}
	}()

	for i := 0; i < limit; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openPorts)
	fmt.Printf("\nOpen ports: %+v\n", openPorts)
}

func worker(ports <-chan int, results chan<- int) {
	for p := range ports {
		address := fmt.Sprintf(target, p)

		conn, err := net.Dial("tcp", address)
		fmt.Printf("checking: %+v\n", address)
		if err != nil {
			results <- 0
			continue
		}

		conn.Close()
		results <- p
	}
}
