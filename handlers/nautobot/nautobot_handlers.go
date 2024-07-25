package nautobot_handlers

import (
    "encoding/json"
    "net/http"
    "semaphore/services/nautobot"
)

func GetNautobotDevices(w http.ResponseWriter, r *http.Request) {
    client := nautobot.NewNautobotClient()
    devices, err := client.GetDevices()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(devices)
}
