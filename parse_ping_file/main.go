package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

/*
 * Parsing output of
 * PS:\> Test-Connection -Count 86400 -ComputerName 1.1.1.1 | Format-Table @{Name='TimeStamp';Expression={Get-Date}},Address,ProtocolAddress,ResponseTime > 1-1-1-1.txt
 * Convert UTF-16LE file to UTF-8 before reading.
 */

var Pit *time.Time

func checkTimestamp(timeStamp string) {
	layout := "15:04:05"
	ts, err := time.Parse(layout, timeStamp)
	if err != nil {
		fmt.Printf("Error parting timestamp %v\n", err)
	}
	if Pit != nil {
		expect := Pit.Add(time.Second)
		if expect.Day() == 2 { // subtract one day when passing midnight
			expect = expect.Add(-time.Hour * 24)
		}
		if ts != expect && ts != Pit.Add(2*time.Second) { // Permit one packet loss
			fmt.Printf("Expected %v got %v\n", expect, ts)
		}
	}
	Pit = &ts
}

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Println("give input filename as argument.")
		os.Exit(1)
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := strings.Split(scanner.Text(), " ")
		if len(t) > 1 && t[1] != "" {
			timeStamp := t[1]
			checkTimestamp(timeStamp)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
