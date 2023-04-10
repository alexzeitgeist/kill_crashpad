# Chrome Crashpad Handler Killer

This simple Go program helps mitigate high CPU usage caused by `chrome_crashpad_handler` processes. The `chrome_crashpad_handler` is part of Google Chrome and handles crash reports, but sometimes it consumes unnecessary CPU resources.

It runs in the background and periodically scans running processes on a Linux system, identifies the `chrome_crashpad_handler` processes, and terminates them.

## Usage

```bash
./kill_crashpad -process=<process_name> -interval=<check_interval>
```

- `process`: Name of the process to monitor and kill (default: "chrome_crashpad")
- `interval`: Interval between process checks (default: 60s)
