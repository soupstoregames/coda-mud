#!/bin/bash

ret=$?
if [ $ret -eq 0 ]; then
  STATE="success"
else
  STATE="failure"
fi

echo -e "[Webhook]: Sending webhook to Discord...\\n";

IMAGE=$1

case $STATE in
  "success" )
    EMBED_COLOR=3066993
    STATUS_MESSAGE="Pushed image"
    ;;

  "failure" )
    EMBED_COLOR=15158332
    STATUS_MESSAGE="Failed to push"
    ;;

  * )
    EMBED_COLOR=0
    STATUS_MESSAGE="Status Unknown"
    ;;
esac


if [ "$AUTHOR_NAME" == "$COMMITTER_NAME" ]; then
  CREDITS="$AUTHOR_NAME authored & committed"
else
  CREDITS="$AUTHOR_NAME authored & $COMMITTER_NAME committed"
fi

TIMESTAMP=$(date --utc +%FT%TZ)
WEBHOOK_DATA='{
  "username": "",
  "avatar_url": "https://cdn.iconscout.com/public/images/icon/free/png-512/docker-logo-35698963e9c1a96c-512x512.png",
  "embeds": [ {
    "color": '$EMBED_COLOR',
    "author": {
      "name": "'"$IMAGE"'",
      "url": "https://travis-ci.org/'"$TRAVIS_REPO_SLUG"'/builds/'"$TRAVIS_BUILD_ID"'"
    },
    "title": "'"$COMMIT_SUBJECT"'",
    "url": "'"$URL"'",
    "description": "'"$STATUS_MESSAGE"'",
    "fields": [],
    "timestamp": "'"$TIMESTAMP"'"
  } ]
}'

(curl --fail --progress-bar -A "Dockerhub-Webhook" -H Content-Type:application/json -H X-Author:k3rn31p4nic#8383 -d "$WEBHOOK_DATA" "$2" \
  && echo -e "\\n[Webhook]: Successfully sent the webhook.") || echo -e "\\n[Webhook]: Unable to send webhook."
