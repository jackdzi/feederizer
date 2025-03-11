package api

import "net/http"

func InitDb() (bool, error) {
	resp, err := http.Post("http://localhost:8080/init", "application/json", nil)
	if err != nil {
		return false, err
	}
	statusCode := resp.StatusCode
	defer resp.Body.Close()

	return (statusCode != http.StatusOK), nil
}
