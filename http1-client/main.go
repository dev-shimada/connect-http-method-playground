package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type apiServicePostRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}
type apiServiceGetResponse struct {
	ID string `json:"id"`
}

func main() {
	var apiResp apiServiceGetResponse
	if p, err := sendPostRequest(); err == nil {
		println("POST Response:", string(p))
		if err := json.Unmarshal(p, &apiResp); err == nil {
			println("Parsed ID:", apiResp.ID)
		} else {
			println("Error parsing ID:", err.Error())
		}
	} else {
		println("POST Error:", err.Error())
	}

	if p, err := sendGetRequest(apiResp.ID); err == nil {
		println("GET Response:", string(p))
	} else {
		println("GET Error:", err.Error())
	}
}

func sendPostRequest() ([]byte, error) {
	reqBody := apiServicePostRequest{
		UserID: "12345",
		Text:   "Hello, World!",
	}
	j, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "http://api:8081/api/v1/ApiService", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func sendGetRequest(id string) ([]byte, error) {
	u, err := url.Parse("http://api:8081/api/v1/ApiService")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("id", id)
	u.RawQuery = q.Encode()
	url := u.String()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
