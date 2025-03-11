package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)


func AddSubscription(data map[string]string) (bool, error) {
  // data := map[string]string{"name": m.inputs[username].Value(), "password": m.inputs[password].Value()}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error parsing string to JSON")
		return false, err
	}

	resp, err := http.Post("http://localhost:8080/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	statusCode := resp.StatusCode
	defer resp.Body.Close()

	return (statusCode != http.StatusOK), nil
}
