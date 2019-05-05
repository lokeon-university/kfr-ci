#!/bin/bash
PROJECT="kfr-ci"
REGISTRY="gcr.io/${PROJECT}"
DOCKERFILE="Dockerfile"
IMAGESDIR="dockerfiles"
IMAGES=(cpp python go javascript)
PUSHEABLE=(cpp python go javascript server bot)
PUSHREGISTRY="$REGISTRY/kfr-"

build-services() {
	docker build -t "$REGISTRY/kfr-bot" -f "tg-bot/$DOCKERFILE" .
	docker build -t "$REGISTRY/kfr-server" -f "server/$DOCKERFILE" .
}

build-ci() {
	cd ci || exit
	CGO_ENABLED=0 GOOS=linux go build -v -o kfr-ci
}

build-images() {
	for image in "${IMAGES[@]}"; do
		cd "$IMAGESDIR/$image" || exit
		docker build -t "$REGISTRY/kfr-$image" -f "$DOCKERFILE" .
		cd - || exit
	done
}

deploy-services() {
	gcloud beta run deploy --image "$REGISTRY/kfr-bot" --update-env-vars "$(cat tg-bot/.env)"
	gcloud beta run deploy --image "$REGISTRY/kfr-server" --update-env-vars "$(cat server/.env)"
}

push-images() {
	for image in "${PUSHEABLE[@]}"; do
		docker push "$PUSHREGISTRY$image"
	done
}

all() {
	#build-services
	#build-images
	push-images
	deploy-services
	build-ci
}

all
