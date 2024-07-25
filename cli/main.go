package main

import (
    "net/http"
    "semaphore/cli/cmd"
    "semaphore/handlers/nautobot" // Adjust the import path as necessary
    "log"
)

func main() {
    // Run CLI commands in a separate goroutine
    go func() {
        cmd.Execute()
    }()

    // Set up HTTP routes for Nautobot
    http.HandleFunc("/api/nautobot/devices", nautobot_handlers.GetNautobotDevices)

    // Start the HTTP server
    log.Println("Starting server on :3000...")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

