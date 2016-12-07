package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
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

	resp, err := h.service.convertToSpeech(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Add("Content-Type", "application/mpeg")
	streamBytes, _ := ioutil.ReadAll(resp.AudioStream)
	w.Write(streamBytes)
}

func decodeJSON(dec *json.Decoder) (interface{}, error) {
	c := request{}
	err := dec.Decode(&c)
	return c, err
}
