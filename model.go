package main

type request struct {
	Body    string `json:"body,omitempty"`
	VoiceId string `json:"voiceId,omitempty"`
	Token   string `json:"token,omitempty`
}
