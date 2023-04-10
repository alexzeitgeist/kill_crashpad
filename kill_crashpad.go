package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	defaultProcessName   = "chrome_crashpad"
	defaultCheckInterval = 60 * time.Second
)

func main() {
	processName := flag.String("process", defaultProcessName, "Name of the process to monitor and kill")
	checkInterval := flag.Duration("interval", defaultCheckInterval, "Interval between process checks")
	flag.Parse()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	processMap := make(map[int]string)

	ticker := time.NewTicker(*checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateProcessMap(processMap)

			for pid, name := range processMap {
				if name == *processName {
					if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
						log.Printf("Error killing process %s with PID %d: %v", *processName, pid, err)
					} else {
						log.Printf("Killed process %s with PID %d\n", *processName, pid)
						delete(processMap, pid) // remove process ID from the map
					}
				}
			}
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Exiting...\n", sig)
			return
		}
	}
}

func updateProcessMap(processMap map[int]string) {
	files, err := os.ReadDir("/proc")
	if err != nil {
		log.Printf("Error reading /proc: %v", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		if name, ok := processMap[pid]; ok {
			if name == "" {
				processMap[pid] = getProcessName(pid)
			}
		} else {
			processMap[pid] = getProcessName(pid)
		}
	}
}

func getProcessName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
