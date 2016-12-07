package main

type content struct {
	BodyXml         string  `json:"bodyXML,omitempty"`
}

type request struct {
	Body 	  string  `json:"body,omitempty"`
	VoiceId	  string  `json:"voiceId,omitempty"`
	Token     string  `json:"token,omitempty`
}