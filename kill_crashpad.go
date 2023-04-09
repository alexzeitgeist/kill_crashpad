package main

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"syscall"
	"time"
)

func findProcessByName(processName string) (int32, error) {
	processList, err := process.Processes()
	if err != nil {
		return -1, err
	}

	for _, proc := range processList {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		if name == processName {
			return proc.Pid, nil
		}
	}

	return -1, nil
}

func killProcess(pid int32) error {
	targetProcess, err := os.FindProcess(int(pid))
	if err != nil {
		return fmt.Errorf("failed to find process with PID %d: %v", pid, err)
	}

	err = targetProcess.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to kill process with PID %d: %v", pid, err)
	}

	return nil
}

func main() {
	processName := flag.String("process", "chrome_crashpad_handler", "Name of the process to monitor and kill")
	checkInterval := flag.Duration("interval", 600*time.Second, "Interval between process checks")
	flag.Parse()

	for {
		pid, err := findProcessByName(*processName)
		if err != nil {
			fmt.Printf("Error finding %s process: %v\n", *processName, err)
		} else if pid != -1 {
			fmt.Printf("Found %s process with PID %d. Killing it...\n", *processName, pid)
			err := killProcess(pid)
			if err != nil {
				fmt.Printf("Error killing %s process: %v\n", *processName, err)
			} else {
				fmt.Printf("Killed %s process with PID %d.\n", *processName, pid)
			}
		}

		time.Sleep(*checkInterval)
	}
}
