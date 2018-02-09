#!/usr/bin/env bash

# This shell script is written to run on a Mac,
# but should be tweakable to work on other OSs.
# You'll need to establish which exe can generate audio on the command line.
# For the mac, that is afplay.

# Make sure your instance of the TTS service is running,
# and ensure LOCAL_SERVER references it.
# Make sure the TOKEN value in the JSON files is correct for whichever server you use.

LOCAL_SERVER=localhost:8080

# If you run this script in it's folder, you should have some example JSON texts.

for JSONFILE in *.json; do
  echo "processing JSONFILE=${JSONFILE}"
  MP3FILE=${JSONFILE}.mp3
  rm -f $MP3FILE

  curl \
    -X PUT \
    --data @${JSONFILE} \
    --header "Content-Type: application/json" \
    ${LOCAL_SERVER}/convert \
    > ${MP3FILE}

  if [ -s $MP3FILE ]; then
    afplay $MP3FILE
  else
    echo "ERROR: failed to generate MP3 file"
  fi

done
