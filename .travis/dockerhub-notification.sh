#!/bin/bash

ret=$?
if [ $ret -eq 0 ]; then
  STATE="success"
else
  STATE="failure"
fi

echo -e "[Webhook]: Sending webhook to Discord...\\n";


case $STATE in
  "success" )
    EMBED_COLOR=3066993
    STATUS_MESSAGE="Pushed"
    ;;

  "failure" )
    EMBED_COLOR=15158332
    STATUS_MESSAGE="Failed"
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
  "avatar_url": "https://iconscout.com/icon/docker-7",
  "embeds": [ {
    "color": '$EMBED_COLOR',
    "author": {
      "name": "$TRAVIS_REPO_SLUG",
    },
    "title": "PUSHING CONTAINER AND VERSION TO DOCKERHUB",
    "url": "URL TO DOCKERHUB HERE",
    "description": "PUT THE VERSION PUSHED HERE",
    "timestamp": "'"$TIMESTAMP"'"
  } ]
}'

(curl --fail --progress-bar -A "Dockerhub-Webhook" -H Content-Type:application/json -H X-Author:k3rn31p4nic#8383 -d "$WEBHOOK_DATA" "$2" \
  && echo -e "\\n[Webhook]: Successfully sent the webhook.") || echo -e "\\n[Webhook]: Unable to send webhook."
