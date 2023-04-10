# Chrome Crashpad Handler Killer

This program helps mitigate high CPU usage caused by `chrome_crashpad_handler` processes. The `chrome_crashpad_handler` is part of Google Chrome and handles crash reports, but sometimes it consumes unnecessary CPU resources.

This is a lightweight application that runs in the background. It periodically scans running processes on a Linux system, identifies the `chrome_crashpad_handler` processes, and terminates them.

