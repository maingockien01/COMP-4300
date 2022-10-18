package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Get
func GetJson(url string, object any) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("Status code is not success")
	}

	json.NewDecoder(res.Body).Decode(object)

	return nil
}

// Post
func PostJson(url string, object []byte) ([]byte, error) {
	body, err := json.Marshal(object)

	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK || res.StatusCode != http.StatusCreated || res.StatusCode != http.StatusNoContent {
		return nil, errors.New("Status code is not success")
	}

	resBody, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	return resBody, nil

}

// Delete
func DeleteJson(url string, jsonBody []byte) ([]byte, error) {
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodDelete, url, bodyReader)

	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

//Websocket
//Open
//Read
//Send message
