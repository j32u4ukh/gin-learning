#!/bin/sh

cd /usr/src/app || exit
# cp .env.local .env

if [ ! -f "/app/initialized" ]; then

    # exec command
    go mod tidy    
    go mod download

    # save flag
    touch /app/initialized
fi

# start app in background

# keep docker running
tail -f  /dev/null