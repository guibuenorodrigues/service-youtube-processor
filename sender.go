package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	webURL = os.Getenv("WEB_URL") // url to web endpoint
)

// PostToProcess - receive the message that will be processed
func (m MessageResponse) PostToProcess() error {

	// call sanitizer method
	message, err := m.Sanitizer()

	if err != nil {
		log.Println(err)
		return err
	}

	message.Post()

	fmt.Println(message.Content.Likes)

	return nil

}

// Post to the API
func (m SanitizedMessage) Post() error {

	j, err := json.Marshal(m.Content)

	if err != nil {
		log.Printf("[X] error to marshal packageMessage: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", webURL, bytes.NewBuffer(j))

	if err != nil {
		log.Printf("[X] error to create new http request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", m.Headers.CorrelationID)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[X] error to send the request: %v", err)
		return err
	}

	defer resp.Body.Close()

	//
	if resp.StatusCode != http.StatusOK {
		fmt.Println("[X] error - response Status:", resp.Status)
		return errors.New("Status Code: " + strconv.Itoa(resp.StatusCode))
	}

	fmt.Println("response Headers:", resp.Header)

	return nil

}
