package nautobot

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

type NautobotClient struct {
    BaseURL  string
    APIToken string
}

func NewNautobotClient(baseURL, apiToken string) *NautobotClient {
    return &NautobotClient{
        BaseURL:  baseURL,
        APIToken: apiToken,
    }
}

func (c *NautobotClient) DoRequest(method, endpoint string, body interface{}) (*http.Response, error) {
    url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request body: %w", err)
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }

    if resp.StatusCode >= 400 {
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        return nil, fmt.Errorf("API request error: %s", string(bodyBytes))
    }

    return resp, nil
}

func (c *NautobotClient) GetDevices() (map[string]interface{}, error) {
    resp, err := c.DoRequest("GET", "dcim/devices/", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return result, nil
}

func main() {
    baseURL := os.Getenv("NAUTOBOT_BASE_URL")
    apiToken := os.Getenv("NAUTOBOT_API_TOKEN")

    client := NewNautobotClient(baseURL, apiToken)
    devices, err := client.GetDevices()
    if err != nil {
        fmt.Printf("Error getting devices: %v\n", err)
        return
    }
    fmt.Printf("Devices: %v\n", devices)
}
