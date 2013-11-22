package main

import (
	"io/ioutil"
	"net/http"
)

// This function fetch the content of a URL will return it as an
// array of bytes if retrieved successfully.
func getContent(url string) ([]byte, []byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	// Defer the closing of the body
	defer resp.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	// Read the header into a byte array
	/*defer resp.Header.Close()
	header, err := ioutil.ReadAll(resp.Header)
	if err != nil {
		return nil, nil, err
	}*/
	// At this point we're done - simply return the bytes
	return nil, body, nil
}
