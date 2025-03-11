package api

import (
	"net/http"
  "fmt"
	"bytes"
	"encoding/json"
)

func AddUser(data map[string]string) (bool, error) {
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
