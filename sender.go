package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	guuid "github.com/google/uuid"
)

// MessageContent content message
type MessageContent struct {
	Package []ContentMessage
}

var (
	packageSize    = 10                   // maximum size per package
	packageCount   = 0                    // count in package
	packageMessage MessageContent         // package message
	webURL         = os.Getenv("WEB_URL") // url to web endpoint
)

// PostToProcess - receive the message that will be processed
func (m MessageResponse) PostToProcess() error {

	// call sanitizer method
	message, err := m.Sanitizer()

	if err != nil {
		log.Println(err)
		return err
	}

	message.addMessageToPackage()

	fmt.Println(message.Content.Likes)

	PostPackage()

	return nil

}

// PostPackage to the API
func PostPackage() error {

	j, err := json.Marshal(&packageMessage.Package)
	// j := []byte(`{"title":"Buy cheese and bread for breakfast."}`)

	if err != nil {
		log.Printf("[X] error to marshal packageMessage: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", webURL, bytes.NewBuffer(j))

	if err != nil {
		log.Printf("[X] error to create new http request: %v", err)
		return err
	}

	x := guuid.New().String()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-ID", x)

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
		return errors.New("Status Code: ")
	}

	fmt.Println("response Headers:", resp.Header)

	return nil

}

func (m SanitizedMessage) addMessageToPackage() {

	// add message to the package

	packageMessage.Package = append(packageMessage.Package, m.Content)

	// add to the package count
	packageCount++
}
