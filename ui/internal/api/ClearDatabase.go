package api

import (
	"net/http"
)

func ClearDatabase() error {
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/user", nil)

	client := &http.Client{}
	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil
}
