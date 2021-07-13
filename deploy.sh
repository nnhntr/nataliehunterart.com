#!/usr/bin/env bash

set -ex

go test ./...

export VERSION
VERSION="0.1.0"

export TAG
TAG="registry.digitalocean.com/crhntr/nnhntr:${VERSION}"

mkdir -p bin

GOOS=linux GOARCH=amd64 go build \
 -ldflags="-w -s" \
 -gcflags=-trimpath="${PWD}" \
 -asmflags=-trimpath="${PWD}" \
 -o bin/nnhntr ./

docker build -t "${TAG}" .

rm -r bin

docker push "${TAG}"

yq w --doc=1 -i deployment.yml 'spec.template.spec.containers[0].image' "${TAG}"

kubectl apply -f deployment.yml

yq d --doc=1 -i deployment.yml 'spec.template.spec.containers[0].image'

echo kubectl get pods -l app=nnhntr --watch