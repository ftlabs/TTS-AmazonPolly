package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type textToSpeechHandler struct {
	service textToSpeechService
}

func newTextToSpeechHandler(service textToSpeechService) textToSpeechHandler {
	return textToSpeechHandler{service: service}
}

func (h *textToSpeechHandler) convertToSpeech(w http.ResponseWriter, req *http.Request) {
	var body io.Reader = req.Body

	dec := json.NewDecoder(body)
	params, _ := decodeJSON(dec)

	resp, err, statusCode := h.service.convertToSpeech(params)

	if err != nil {
		if strings.Contains(err.Error(), "TextLengthExceededException") {
			log.Errorf("Character limit exceeded", err.Error())
			w.WriteHeader(400)
			w.Write([]byte("Character limit exceeded. Service restricted to 1500 characters per request"))
			return
		}
		http.Error(w, err.Error(), statusCode)
	} else {
		w.WriteHeader(statusCode)
		w.Header().Add("Content-Type", "application/mpeg")
		streamBytes, _ := ioutil.ReadAll(resp.AudioStream)
		w.Write(streamBytes)
	}
}

func decodeJSON(dec *json.Decoder) (interface{}, error) {
	c := request{}
	err := dec.Decode(&c)
	return c, err
}
