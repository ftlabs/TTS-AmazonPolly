package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"net/http"
	"strings"
)

type textToSpeechService interface {
	convertToSpeech(thing interface{}) (*polly.SynthesizeSpeechOutput, error, int)
}

type textToSpeechServiceImpl struct {
	awsCreds credentials.Credentials
	token    string
}

func newTextToSpeechService(awsCreds credentials.Credentials, token string) textToSpeechService {
	return &textToSpeechServiceImpl{awsCreds: awsCreds, token: token}
}

func (tts *textToSpeechServiceImpl) convertToSpeech(thing interface{}) (*polly.SynthesizeSpeechOutput, error, int) {
	input := thing.(request)

	if tts.token != input.Token {
		return &polly.SynthesizeSpeechOutput{}, errors.New("Token " + input.Token + " is invalid!"), http.StatusUnauthorized
	}

	if len(input.Body) == 0 || input.Body == "" {
		return &polly.SynthesizeSpeechOutput{}, errors.New("Unable to process input text:" + input.Body), http.StatusBadRequest
	}

	sess, err := session.NewSession()
	if err != nil {
		return &polly.SynthesizeSpeechOutput{}, fmt.Errorf("Failed to create session", err), http.StatusInternalServerError
	}

	pollyService := polly.New(sess, &aws.Config{Region: aws.String("eu-west-1"), Credentials: &tts.awsCreds})

	// Simple check for the text starting with <speak>,
	// from which we infer whether it should be treated as SSML (or plain text, by default).
	// N.B., once in ssml mode, broken syntax will cause the TTS to fail, e.g. with no closing </speak>.
	textType := "text"
	trimmedText := strings.TrimSpace(input.Body)
	if strings.HasPrefix(trimmedText, "<speak>") {
		textType = "ssml"
	}

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(input.Body),
		VoiceId:      aws.String(input.VoiceId),
		TextType:     aws.String(textType),
	}

	ssiResp, err := pollyService.SynthesizeSpeech(params)

	if err != nil {
		return &polly.SynthesizeSpeechOutput{}, fmt.Errorf("Unable to access Amazon Polly %v", err), http.StatusInternalServerError
	}
	return ssiResp, nil, http.StatusOK
}
