package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const authServiceURL = "http://authentication-service/authenticate"

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	err := app.pushToQueue("broker_hit", r.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Message = "Received request"

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(out)
}

func (app *Config) BrokerAuth(w http.ResponseWriter, r *http.Request) {
	// create a variable matching the structure of the JSON we expect from the front end
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// read posted json into our variable
	_ = readJSON(w, r, &requestPayload)

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	// call the authentication-service
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		_ = errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	// create variable we'll read the response.Body from the authentication-service into
	var jsonFromService jsonResponse

	// decode the json we get from the authentication-service into our variable
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send json back to our end user
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(out)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		log.Println(err)
		return err
	}

	payload := Payload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "    ")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}