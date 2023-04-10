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
	defProcName   = "chrome_crashpad"
	defCheckIntvl = 60 * time.Second
)

func main() {
	procName := flag.String("process", defProcName, "Name of the process to monitor and kill")
	checkIntvl := flag.Duration("interval", defCheckIntvl, "Interval between process checks")
	flag.Parse()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	procMap := make(map[int]string)
	targetPIDs := make(map[int]struct{})

	ticker := time.NewTicker(*checkIntvl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateProcMap(procMap, targetPIDs, *procName)

			for pid := range targetPIDs {
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					log.Printf("Error killing process %s with PID %d: %v", *procName, pid, err)
				} else {
					log.Printf("Killed process %s with PID %d\n", *procName, pid)
					delete(procMap, pid) // remove process ID from the map
					delete(targetPIDs, pid) // remove process ID from the targetPIDs map
				}
			}
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Exiting...\n", sig)
			return
		}
	}
}

func updateProcMap(procMap map[int]string, targetPIDs map[int]struct{}, targetName string) {
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

		if name, ok := procMap[pid]; ok {
			if name == "" {
				procMap[pid] = getProcName(pid)
				if procMap[pid] == targetName {
					targetPIDs[pid] = struct{}{}
				}
			}
		} else {
			procMap[pid] = getProcName(pid)
			if procMap[pid] == targetName {
				targetPIDs[pid] = struct{}{}
			}
		}
	}
}

func getProcName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
