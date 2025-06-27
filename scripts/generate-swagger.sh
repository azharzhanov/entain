#!/usr/bin/env bash
# https://goswagger.io/install.html

# Declare variables.
option=$1
alias swagger="docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger:v0.30.0"

# Generate swagger documentation.
docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger:v0.30.0 generate spec -o ./api/generated/swagger/app.swagger.json

# Check if the "validate" option specified.
if [ "$option" = 'validate' ]; then
    echo "Validate the provided swagger document against a swagger spec"
    docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger:v0.30.0 validate ./api/generated/swagger/app.swagger.json --skip-warnings
fi
