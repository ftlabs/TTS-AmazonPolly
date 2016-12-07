# TTS-AmazonPolly
Backend application which utilises Amazon's Polly service to convert text into speech

## Endpoints

### /convert
#### PUT
Body: Plain text to be converted
VoiceId: one of the Amazon Polly Voices
* choose from `http://docs.aws.amazon.com/polly/latest/dg/API_SynthesizeSpeech.html`
Token: Authorization token


Example:
`curl -X PUT --data "{\"Body\":\"HelloWorld\",\"VoiceId\":\"Geraint\",\"Token\":\"MyToken\"}" localhost:8080/convert`  
    