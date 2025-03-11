package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackdzi/feederizer/ui/internal/config"
)


func CheckUser(user string, pass string) (bool, error) {
	passDecrypt, _ := config.Decrypt(pass)
	data := map[string]string{"name": user, "password": passDecrypt}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error parsing string to JSON")
		return false, err
	}

	resp, err := http.Post("http://localhost:8080/user/check", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	statusCode := resp.StatusCode
	defer resp.Body.Close()

	return (statusCode == http.StatusOK), nil
}
