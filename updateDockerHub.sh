#!/usr/bin/env bash

version="v1.6"

# build images
docker build . --tag camerondugan/neorgify:"$version"
docker build . --tag camerondugan/neorgify:latest

# push images
docker push camerondugan/neorgify:"$version"
docker push camerondugan/neorgify:latest
