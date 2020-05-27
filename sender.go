package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	// webURL = os.Getenv("WEB_URL") // url to web endpoint
	webURL = "http://soliveboa.com.br/services/live"
)

type sendStatus struct {
	liveID string
	code   int
	status bool
	err    error
}

// PostToProcess - receive the message that will be processed
func (m MessageResponse) PostToProcess() error {

	// call sanitizer method
	message, err := m.Sanitizer()

	if err != nil {
		log.Println(err)
		return err
	}

	//c := make(chan sendStatus)
	//go message.Post(c)
	err = message.Post()

	if err != nil {
		return err
	}

	return nil
}

// Post to the API
func (m SanitizedMessage) Post() error {

	j, err := json.Marshal(m.Content)

	// fmt.Println(string(j))
	if err != nil {
		log.Printf("[X] error to marshal packageMessage: %v", err)
		//c <- sendStatus{liveID: m.Content.IDLive, status: false, err: err}
		return err
	}

	req, err := http.NewRequest("POST", webURL, bytes.NewBuffer(j))

	if err != nil {
		log.Printf("[X] error to create new http request: %v", err)
		//c <- sendStatus{liveID: m.Content.IDLive, status: false, err: err}
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", m.Headers.CorrelationID)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[X] error to send the request: %v", err)
		//c <- sendStatus{liveID: m.Content.IDLive, status: false, err: err}
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("LiveID: %v Status: FALSE - StatusCode: %v", m.Content.IDLive, resp.Status)
		// c <- sendStatus{liveID: m.Content.IDLive, code: resp.StatusCode, status: false, err: errors.New("Status Code: " + strconv.Itoa(resp.StatusCode))}
		return err
	}

	// message for log
	fmt.Printf("LiveID: %v Status: TRUE - StatusCode: %v", m.Content.IDLive, resp.Status)
	fmt.Println("")

	return nil

}
