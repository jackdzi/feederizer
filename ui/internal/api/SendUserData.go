package api

import (
	"net/http"
  "fmt"
	"bytes"
	"encoding/json"
)


func SendUserData(data map[string]string) (bool, bool, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error parsing string to JSON")
		return false, false, err
	}

	resp, err := http.Post("http://localhost:8080/credentials", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, false, err
	}
  statusCode := resp.StatusCode
	defer resp.Body.Close()

	return (statusCode == http.StatusInternalServerError), (statusCode == http.StatusUnauthorized), nil
}
