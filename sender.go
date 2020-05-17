package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	// webURL = os.Getenv("WEB_URL") // url to web endpoint
	webURL = "http://tractum.serveo.net/soliveboa/services/live"
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

	c := make(chan sendStatus)
	go message.Post(c)

	return nil
}

// Post to the API
func (m SanitizedMessage) Post(c chan sendStatus) {

	j, err := json.Marshal(m.Content)

	// fmt.Println(string(j))
	if err != nil {
		log.Printf("[X] error to marshal packageMessage: %v", err)
		c <- sendStatus{liveID: m.Content.IDLive, status: false, err: err}
	}

	req, err := http.NewRequest("POST", webURL, bytes.NewBuffer(j))

	if err != nil {
		log.Printf("[X] error to create new http request: %v", err)
		c <- sendStatus{liveID: m.Content.IDLive, status: false, err: err}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", m.Headers.CorrelationID)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[X] error to send the request: %v", err)
		c <- sendStatus{liveID: m.Content.IDLive, status: false, err: err}
	}

	defer resp.Body.Close()

	log.Print("[*] LiveID: " + m.Content.IDLive)

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		fmt.Print("Status: FALSE - StatusCode:", resp.Status)
		c <- sendStatus{liveID: m.Content.IDLive, code: resp.StatusCode, status: false, err: errors.New("Status Code: " + strconv.Itoa(resp.StatusCode))}

		// return errors.New("Status Code: " + strconv.Itoa(resp.StatusCode))
	}

	// message for log
	fmt.Print("Status: TRUE - StatusCode:", resp.Status)
	fmt.Println("")

	c <- sendStatus{liveID: m.Content.IDLive, code: resp.StatusCode, status: true}

}
