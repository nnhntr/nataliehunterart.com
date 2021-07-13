#!/usr/bin/env bash

set -ex

go test ./...

export VERSION
VERSION="$(doctl registry repository list-tags nnhntr --output json | jq -r .[0].tag | awk '{split($0,a,"."); print a[3] "."  a[2] + 1 "." a[1]}')"

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