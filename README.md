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

## Credentials

Can be found in ftlabs lastpass

## building ...locally (a reminder)

* set up golang (with correct folder hierarchy and $GOPATH, etc)
* put the repo code in the correct spot: go/src/github.com/ftlabs/TTS-AmazonPolly
* pull in all the appropriate dependencies: go get -u -v github.com/ftlabs/TTS-AmazonPolly
* build the executable: go install
* ensure the .gitignore file includes .env
* ensure the assorted env vars are in .env
   * AWS_ACCESS_ID
   * AWS_ACCESS_KEY (the AWS secret)
   * TOKEN (for authenticating access to this service)
* actually, since we are not auto-parsing the .env file, manually export all the env vars
* run the executable: $GOPATH/bin/TTS-AmazonPolly

## testing

* start your local server
* cd test/command-line
* ./mac.sh (if you are on a mac, or tweak it if on a different OS)
* listen to the TTS

## deploying to heroku (a reminder)

* If you have made changes to any dependencies (added, or updated)
   * install and run godep (see https://github.com/tools/godep)
   * commit the changes
* Push to heroku.
